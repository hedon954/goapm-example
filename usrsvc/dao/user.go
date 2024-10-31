package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hedon954/goapm/apm"
	"github.com/redis/go-redis/v9"
)

type UserDao struct {
	DB  *sql.DB
	RDB *redis.Client
}

func (dao *UserDao) Get(ctx context.Context, uid int64) map[string]any {
	userCache := dao.RDB.Get(ctx, userKey(uid)).Val()
	if userCache != "" {
		userInfo := make(map[string]any)
		err := json.Unmarshal([]byte(userCache), &userInfo)
		if err == nil {
			return userInfo
		}
	}

	apm.Logger.Debug(ctx, "userDao.Get", map[string]any{
		"uid": uid,
	})

	raw, err := dao.DB.QueryContext(ctx,
		"select * from `t_user` where id = ?", uid)
	if err != nil {
		apm.Logger.Error(ctx, "userDao.Get", err, map[string]any{
			"uid": uid,
		})
		return nil
	}
	info := apm.DBUtils.QueryFirst(raw, raw.Err())

	apm.Logger.Debug(ctx, "userDao.Get", map[string]any{
		"info": info,
	})

	if info != nil {
		userInfo, err := json.Marshal(info)
		if err == nil {
			dao.RDB.Set(ctx, userKey(uid), string(userInfo), time.Minute*10)
		}
	}

	fmt.Println("info", info)

	return info
}

func userKey(uid int64) string {
	return fmt.Sprintf("%s:%s:%d", "usersc", "uinfo", uid)
}
