package heimdall

import (
	"github.com/jackc/pgx"
	"gopkg.in/inconshreveable/log15.v2"
	"context"
)

func Listen(pool *pgx.ConnPool, config *Config) error {
	logger := log15.New("MODULE", "Listen")
	workerPool := NewWorkerPool(config.MaxWorkers)

	conn, err := pool.Acquire()
	if err != nil {
		logger.Crit("unable to acquire connection", "error", err)
		return err
	}

	defer pool.Release(conn)

	err = conn.Listen(config.DB.ProcedureName)
	if err != nil {
		logger.Crit("unable to start listening", "error", err)
		return err
	}

	for {
		notification, err := conn.WaitForNotification(context.Background())
		if err != nil {
			logger.Crit("unable to wait for notification", "error", err)
			return err
		}

		workerPool.Schedule(func() {
			work := NewDispatchNotification(notification.Channel, notification.Payload)
			logger.Debug("payload prepared", "payload", work.Payload, "channel", work.Channel)
			if err := work.Handle(workerPool, config); err != nil {
				logger.Error("error happened while handling notification", "error", err)
			}
		})
	}
}