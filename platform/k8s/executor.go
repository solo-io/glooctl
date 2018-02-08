package k8s

import (
	"log"
	"strings"
	"time"

	gluev1 "github.com/solo-io/glue/pkg/api/types/v1"
	crdclient "github.com/solo-io/glue/pkg/platform/kube/crd/client/clientset/versioned"
	"github.com/solo-io/glue/pkg/platform/kube/crd/solo.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Executor struct {
	cfg    *rest.Config
	client *crdclient.Clientset
}

func NewExecutor(config interface{}) *Executor {
	s, ok := config.(string)
	if !ok {
		s = ""
	}
	cfg, err := getClientConfig(s)
	if err != nil {
		log.Fatal("Cannot create k8s client", err)
	}
	client, err := crdclient.NewForConfig(cfg)
	if err != nil {
		log.Fatal("Cannot create glue CRDs clientset", err)
	}
	return &Executor{
		cfg:    cfg,
		client: client,
	}
}

func (e *Executor) RunCreateUpstreamFromFile(file, namespace string, wait int) {

}

func (e *Executor) RunCreateUpstream(name, namespace, utype, spec string, wait int) {
	x := upstreamFromArgs(name, utype, spec)
	e.client.GlueV1().Upstreams(namespace).Create(x)
	e.wait(wait, func(e *Executor) bool {
		return true
	})
}

func getClientConfig(kubeConfig string) (*rest.Config, error) {
	if kubeConfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeConfig)
	}
	return rest.InClusterConfig()
}

func upstreamFromArgs(name, utype, spec string) *v1.Upstream {
	ss := strings.Split(spec, ";")
	s := make(map[string]interface{})
	for _, v := range ss {
		kw := strings.Split(v, "=")
		if len(kw) == 2 {
			s[kw[0]] = kw[1]
		}
	}

	return &v1.Upstream{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.DeepCopyUpstream{
			Name: name,
			Type: gluev1.UpstreamType(utype),
			Spec: s,
		},
	}
}

func (e *Executor) wait(w int, cb func(e *Executor) bool) {
	for i := 0; i < w; i++ {
		if cb(e) {
			break
		}
		time.Sleep(time.Second)
	}
}
