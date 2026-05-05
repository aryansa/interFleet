package device

type IDeviceReadRepository interface {
	GetById(id string) (*Device, error)
}
