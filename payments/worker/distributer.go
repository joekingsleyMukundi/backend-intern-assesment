package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributer interface {
	DistributetaskInitiatePayment(
		crx context.Context,
		payload *PayloadInitiatePayment,
		opts ...asynq.Option,
	) error
}
type RedisTaskDistributer struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOPt asynq.RedisClientOpt) TaskDistributer {
	client := asynq.NewClient(redisOPt)
	return &RedisTaskDistributer{
		client: client,
	}
}
