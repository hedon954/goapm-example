package dao

import (
	"context"
	"database/sql"

	"github.com/hedon954/goapm/apm"
	"github.com/redis/go-redis/v9"
)

type SkuDao struct {
	DB  *sql.DB
	RDB *redis.Client
}

func (s *SkuDao) Get(ctx context.Context, id int64) map[string]any {
	raw, err := s.DB.QueryContext(ctx,
		"select * from `t_sku` where id = ?", id)
	if err != nil {
		return nil
	}
	info := apm.DBUtils.QueryFirst(raw, raw.Err())
	return info
}

func (s *SkuDao) Decr(ctx context.Context, id int64, num int32) (sql.Result, error) {
	return s.DB.ExecContext(ctx,
		"update `t_sku` set num = num - ? where id = ? and (num - ? >= 0)", num, id, num)
}
