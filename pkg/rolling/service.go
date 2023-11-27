package rolling

var (
	ServiceOptionsMap = map[string]Options{}
	ServiceFactoryMap = map[string]ServiceFactory{}
)

type ServiceFactory func(o Options) (Service, error)
type Service interface {
	Prepare() error
	FilterNodes()
	RestartNode() error
}
