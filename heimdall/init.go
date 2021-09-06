package heimdall

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/log/log15adapter"
	"gopkg.in/inconshreveable/log15.v2"
)


func upsertDBNotificationProcedure(pool *pgx.ConnPool, name string, logger log15.Logger) error {
	query := fmt.Sprintf(queryCreateOrReplacePublishFunction, name, name)

	if operation, err := pool.Exec(query); err != nil {
		logger.Error("unable to create procedure", "error", err)
		return err
	} else {
		logger.Info("created procedure", "output", operation)
	}

	return nil
}

func installTableTriggers(pool *pgx.ConnPool, notifications []NotificationConfig, dbConfig DBConfig, logger log15.Logger) error {
	for _, notification :=  range notifications {
		if notification.TableName == "" {
			logger.Crit("ignoring invalid table name", "configuration", notification)
			return errors.New("empty table name")
		}

		triggerName := notification.TriggerName
		if triggerName == "" {
			triggerName = fmt.Sprintf("trigger_%s_on_%s", dbConfig.ProcedureName, notification.TableName)
		}

		dropQuery := fmt.Sprintf(queryDropTriggerStatement, triggerName, notification.TableName)
		if _, err := pool.Exec(dropQuery); err != nil {
			logger.Crit("Unable to drop previous trigger", "error", err)
			return err
		}

		query := fmt.Sprintf(queryCreateTriggerStatement, triggerName, notification.TableName, dbConfig.ProcedureName)
		if _, err := pool.Exec(query); err != nil {
			logger.Crit("unable to install table trigger", "error", err)
			return err
		}
	}

	logger.Debug("tables triggers installed successfully")
	return nil
}

// Bootup starts the system up with given config
// This initializes any triggers/procedures and starts listening to connection for notifications
func Bootup(config *Config) (*pgx.ConnPool, error) {
	logger := log15.New()
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     config.DB.Host,
			Port: 	  config.DB.Port,
			User:     config.DB.User,
			Password: config.DB.Password,
			Database: config.DB.Database,
			Logger:   log15adapter.NewLogger(logger),
		},

		MaxConnections: config.DB.MaxConnections,
	}

	pool, err := pgx.NewConnPool(connPoolConfig)
	if err != nil {
		log15.Crit("unable to create connection pool", "error", err)
		return nil, err
	}

	err = upsertDBNotificationProcedure(pool, config.DB.ProcedureName, logger)
	if err != nil {
		return pool, err
	}

	err = installTableTriggers(pool, config.Notifications, config.DB, logger)
	if err != nil {
		return pool, err
	}

	return pool, nil
}

