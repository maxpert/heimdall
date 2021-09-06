package heimdall

type HttpHookConfig struct {
	Url           string            `mapstructure:"url"`
	Method 		  string 			`mapstructure:"method"`
	Headers       map[string]string `mapstructure:"headers"`
	Retry         uint              `mapstructure:"retry"`
	RetryStrategy string            `mapstructure:"retry_strategy"`
	Timeout       int               `mapstructure:"timeout"`
}

type DBConfig struct {
	Host           string `mapstructure:"host"`
	Port           uint16 `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Database       string `mapstructure:"db_name"`
	MaxConnections int    `mapstructure:"max_connections"`
	ProcedureName  string `mapstructure:"procedure_name"`
}

type LoggerConfig struct {
	Context []string `mapstructure:"context"`
}

type NotificationConfig struct {
	TriggerName string          `mapstructure:"trigger_name"`
	TableName   string          `mapstructure:"table_name"`
	OnOperation map[string]bool `mapstructure:"on"`
}

type Config struct {
	DB            DBConfig                  `mapstructure:"database"`
	Logger        LoggerConfig              `mapstructure:"logger"`
	HttpHooks     map[string]HttpHookConfig `mapstructure:"http_hooks"`
	Notifications []NotificationConfig      `mapstructure:"notifications"`
	MaxWorkers    int                       `mapstructure:"max_workers"`
}
