package api

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/hedon954/goapm/apm"
	"github.com/spf13/cast"
	"go.opentelemetry.io/otel/trace"

	"github.com/hedon954/goapm-example/ordersvc/grpcclient"
	"github.com/hedon954/goapm-example/ordersvc/metric"
	"github.com/hedon954/goapm-example/protos"
)

type Order struct {
	Tracer trace.Tracer
	DB     *sql.DB
}

func (o *Order) Add(w http.ResponseWriter, request *http.Request) {
	ctx, span := o.Tracer.Start(request.Context(), "order.Add-Start")
	defer span.End()

	// get request body
	values := request.URL.Query()
	var (
		uid, _   = strconv.Atoi(values.Get("uid"))
		skuID, _ = strconv.Atoi(values.Get("sku_id"))
		num      = cast.ToInt32(values.Get("num"))
	)

	// check user info
	userInfo, err := grpcclient.UserClient.GetUser(ctx, &protos.User{
		Id: int64(uid),
	})
	if err != nil {
		apm.Logger.Error(ctx, "get user info from user service", err, map[string]any{
			"uid":    uid,
			"sku_id": skuID,
			"num":    num,
		})
		HttpStatus.Error(w, err.Error(), nil)
		return
	}
	if userInfo.Id == 0 {
		HttpStatus.Error(w, "user not found from user service", nil)
		return
	}

	// deduct stock
	res, err := grpcclient.SkuClient.DecreaseStock(ctx, &protos.Sku{
		Id:  int64(skuID),
		Num: num,
	})
	if err != nil {
		apm.Logger.Error(ctx, "createOrder", err, map[string]any{
			"sku_id": skuID,
			"num":    num,
		})
		HttpStatus.Error(w, err.Error(), nil)
		return
	}

	// create order
	_, err = o.DB.ExecContext(ctx,
		"INSERT INTO `t_order` (`order_id`, `sku_id`, `num`, `price`, `uid`) VALUES (?, ?, ?, ?, ?)",
		uuid.NewString(), skuID, num, res.Price, uid,
	)
	if err != nil {
		apm.Logger.Error(ctx, "createOrder", err, map[string]any{
			"uid":    uid,
			"sku_id": skuID,
			"num":    num,
		})
		HttpStatus.Error(w, err.Error(), nil)
		return
	}

	// add order-success metric
	metric.OrderSuccessCounter.WithLabelValues(strconv.Itoa(skuID)).Inc()
	log.Println("order success", skuID)

	// return
	HttpStatus.Ok(w)
}
