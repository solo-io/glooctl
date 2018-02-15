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

type VhostParams struct {
	Name string
}

type Executor interface {
	RunCreate(gparams *GlobalParams, sparams interface{})
	RunUpdate(gparams *GlobalParams, sparams interface{})
	RunDelete(gparams *GlobalParams, sparams interface{})
	RunGet(gparams *GlobalParams, sparams interface{})
	RunDescribe(gparams *GlobalParams, sparams interface{})
}
