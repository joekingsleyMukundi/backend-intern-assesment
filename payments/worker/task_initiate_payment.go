package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskInitiatePayment = "task:initiate_payment"

type PayloadInitiatePayment struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Amount   string `json:"amount"`
}

func (distributor *RedisTaskDistributer) DistributetaskInitiatePayment(
	ctx context.Context,
	payload *PayloadInitiatePayment,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal: %s", err)
	}
	task := asynq.NewTask(TaskInitiatePayment, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enque a task: %s", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}
