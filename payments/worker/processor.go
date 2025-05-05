package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/joekingsleyMukundi/backend-intern-assesment/common/db/sqlc"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritiacal = "critiacal"
	QueueDefaut    = "default"
)

type TaskProccessor interface {
	Start() error
	ProssesstaskInitiatePayment(cts context.Context, task *asynq.Task) error
}

type RedisTaskProccessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisProccessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProccessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritiacal: 10,
				QueueDefaut:    5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("Proccesstype failed")
			}),
			Logger: NewLogger(),
		},
	)
	return &RedisTaskProccessor{
		server: server,
		store:  store,
	}
}

func (proccessor *RedisTaskProccessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskInitiatePayment, proccessor.ProssesstaskInitiatePayment)
	proccessor.server.Start(mux)
	return nil
}
