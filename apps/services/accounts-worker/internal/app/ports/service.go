package ports

// AccountService will be a placeholder
type AccountService interface {
	PublishAccountCreated(queueName string) error
}
