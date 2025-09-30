package hello

// Service defines the behavior of our Hello service
type Service interface {
	SayHello(name string) string
	SayGoodbye(name string) string
}

// helloService implements the Service interface
type helloService struct{}

func (helloService) SayHello(name string) string {
	return "Hello, " + name
}

func (helloService) SayGoodbye(name string) string {
	return "Goodbye, " + name
}

// NewService is a constructor for helloService
func NewService() Service {
	return helloService{}
}
