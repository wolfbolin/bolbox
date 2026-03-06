package configs

type Options struct {
	ExitOnHelp bool
	ParseFlows []Flow
}

type Flow string

const (
	FlowEnv  Flow = "env"
	FlowFlag Flow = "flag"
)

func DefaultOptions() *Options {
	return &Options{
		ExitOnHelp: true,
		ParseFlows: []Flow{FlowEnv, FlowFlag},
	}
}
