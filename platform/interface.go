package platform

type Executor interface {
	RunCreateUpstreamFromFile(file, namespace string, wait int)
	RunCreateUpstream(name, namespace, utype, spec string, wait int)
}
