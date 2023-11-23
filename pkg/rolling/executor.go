package rolling

var (
	ExecutorOptionsMap = map[string]Options{}
	ExecutorFactoryMap = map[string]ExecutorFactory{}
)

type ExecutorFactory func(o Options) Executor
type Executor interface {
	Prepare() error
	FilterNodes()
	RestartNode() error
}
