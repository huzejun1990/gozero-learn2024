// @Author huzejun 2024/11/10 21:38:00
package queue

import (
	"context"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"user-api/internal/svc"
)

func Consumers(kqConf kq.KqConf, ctx context.Context, serviceContext *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(kqConf, NewLogConsumer(ctx, serviceContext)),
	}
}
