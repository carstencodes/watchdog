package notifications

type Args interface {
	GetBool(key string) (bool, error)
	GetFloat(key string) (float64, error)
	GetInt(key string) (int64, error)
	GetString(key string) (string, error)
}

type Notifier interface {
	Connect() error
	Disconnect() error
	Send(message Message, messageArgs Args) error
}

type NotifierCreatorFunc func() (Notifier, error)
