package heimdall

import (
	"encoding/json"
	"github.com/pkg/errors"
	"fmt"
	"strings"
)

type notificationInfo struct {
	Table  string `json:"table"`
	Action string `json:"action"`
}

type DispatchNotification struct {
	Channel string
	Payload string
	Info    *notificationInfo
}

var (
	PayloadEmptyError = errors.New("Empty payload")
)

func NewDispatchNotification(channel, payload string) *DispatchNotification {
	return &DispatchNotification{
		channel,
		payload,
		nil,
	}
}

func (n *DispatchNotification) readPayload() error {
	info := &notificationInfo{}
	if err := json.Unmarshal([]byte(n.Payload), info); err != nil {
		return err
	}

	n.Info = info
	return nil
}

func (n *DispatchNotification) GetMessage() string {
	return n.Payload
}

func (n *DispatchNotification) GetTable() string {
	if n.Info == nil {
		panic(PayloadEmptyError)
	}

	return n.Info.Table
}

func (n *DispatchNotification) GetAction() string {
	if n.Info == nil {
		panic(PayloadEmptyError)
	}

	return n.Info.Action
}

func (n *DispatchNotification) filterMatchingHooks(hooks []NotificationConfig) ([]NotificationConfig, error) {
	if n.Info == nil {
		return nil, PayloadEmptyError
	}

	filtered := make([]NotificationConfig, 0)
	for _, hook := range hooks {
		if hook.TableName == n.Info.Table {
			if hook.TriggerName == "" {
				return nil, errors.New(fmt.Sprintf("Hooks not configured for %s", n.Info.Table))
			}

			if hook.OnOperation != nil {
				filtered = append(filtered, hook)
				continue
			}

			if allow, matches := hook.OnOperation[strings.ToLower(n.Info.Action)]; allow && matches {
				filtered = append(filtered, hook)
				continue
			}
		}
	}

	return filtered, nil
}

func (n *DispatchNotification) triggerNotifications(
	workerPool WorkerPool,
	notifications []NotificationConfig,
	hooks map[string]HttpHookConfig,
) error {
	for _, notification := range notifications {
		if _, ok := hooks[notification.TriggerName]; !ok {
			return errors.New(fmt.Sprintf("Missing hook configration %s", notification.TriggerName))

		}
	}

	for _, notification := range notifications {
		hook, _ := hooks[notification.TriggerName]
		workerPool.Schedule(func() {
			NewHttpHookPublisher(hook).Send(n)
		})
	}

	return nil
}

func (n *DispatchNotification) Handle(workerPool WorkerPool, config *Config) error {
	if err := n.readPayload(); err != nil {
		return err
	}

	notifications, err := n.filterMatchingHooks(config.Notifications)
	if err != nil {
		return err
	}

	n.triggerNotifications(workerPool, notifications, config.HttpHooks)
	return nil
}
