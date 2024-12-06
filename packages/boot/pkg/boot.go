package boot

type Service struct {
	Name string
}

func NewService(serviceName string) *Service {
	return &Service{
		Name: serviceName,
	}
}

func (s Service) GetServiceName() string {
	return s.Name
}
