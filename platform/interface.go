package platform

type GlobalParams struct {
	FileName  string
	Namespace string
	WaitSec   int
}

type UpstreamParams struct {
	Name  string
	UType string
	Spec  map[string]interface{}
}

type Executor interface {
	RunCreateUpstream(gparams *GlobalParams, uparams *UpstreamParams)
	RunUpdateUpstream(gparams *GlobalParams, uparams *UpstreamParams)
	RunDeleteUpstream(gparams *GlobalParams, uparams *UpstreamParams)
	RunGetUpstream(gparams *GlobalParams, uparams *UpstreamParams)
	RunDescribeUpstream(gparams *GlobalParams, uparams *UpstreamParams)
}
