// @Author huzejun 2024/11/10 21:33:00
package queue

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"user-api/internal/svc"
)

type LogConsumer struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewLogConsumer(ctx context.Context, svc *svc.ServiceContext) *LogConsumer {
	return &LogConsumer{
		ctx: ctx,
		svc: svc,
	}

}

func (*LogConsumer) Consume(ctx context.Context, key, val string) error {
	logx.Info("consume log: key=%s, val=%s", val)
	return nil
}
