package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/tonisco/simple-bank-go/db/sqlc"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(opts asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(opts, asynq.Config{})
	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}
