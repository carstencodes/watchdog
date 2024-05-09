package notifications

type nullNotifier struct{}

func newNullNotifier() (Notifier, error) {
	return nullNotifier{}, nil
}

func (nn nullNotifier) Connect() error {
	return nil
}

func (nn nullNotifier) Disconnect() error {
	return nil
}

func (nn nullNotifier) Send(msg Message, messageArgs Args) error {
	return nil
}

func init() {
	notificationClients[""] = newNullNotifier
}
