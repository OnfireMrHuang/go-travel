package service

import (
	"context"
	"go-travel/blog-service/global"
	"go-travel/blog-service/internal/dao"

	otgorm "github.com/smacker/opentracing-gorm"
)

type Service struct {
	ctx context.Context
	dao *dao.Dao
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = dao.New(otgorm.SetSpanToGorm(svc.ctx, global.DBEngine))
	return svc
}
