package api

import (
	"context"

	"github.com/spf13/cast"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hedon954/goapm-example/protos"
	"github.com/hedon954/goapm-example/usrsvc/dao"
)

type User struct {
	protos.UnimplementedUserServiceServer
	Dao *dao.UserDao
}

func (u *User) GetUser(ctx context.Context, req *protos.User) (*protos.User, error) {
	info := u.Dao.Get(ctx, req.Id)
	if info == nil {
		return nil, status.Errorf(codes.NotFound, "user not found in db")
	}
	return &protos.User{
		Name: info["name"].(string),
		Id:   cast.ToInt64(info["id"]),
	}, nil
}
