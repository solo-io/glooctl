package k8s

import (
	"fmt"
	"log"
	"time"

	gluev1 "github.com/solo-io/glue/pkg/api/types/v1"
	crdclient "github.com/solo-io/glue/pkg/platform/kube/crd/client/clientset/versioned"
	"github.com/solo-io/glue/pkg/platform/kube/crd/solo.io/v1"
	platform "github.com/solo-io/gluectl/platform"
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

func (e *Executor) RunCreateUpstream(gparams *platform.GlobalParams, uparams *platform.UpstreamParams) {

	if uparams.Name == "" || uparams.UType == "" {
		log.Fatal("Both Name and Type of the Upstream must be provided")
	}

	x := upstreamFromArgs(uparams.Name, uparams.UType, uparams.Spec)
	e.client.GlueV1().Upstreams(gparams.Namespace).Create(x)
	err := e.wait(gparams.WaitSec, func(e *Executor) bool {
		s := e.getUpstreamCrdStatus(uparams.Name, gparams.Namespace)
		if s != "" {
			log.Printf("Create Upstream Status: %s\n", s)
			return true
		}
		return false
	})
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Upstream created")
	}
}

func getClientConfig(kubeConfig string) (*rest.Config, error) {
	if kubeConfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeConfig)
	}
	return rest.InClusterConfig()
}

func upstreamFromArgs(name, utype string, spec map[string]interface{}) *v1.Upstream {

	return &v1.Upstream{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.DeepCopyUpstream{
			Name: name,
			Type: gluev1.UpstreamType(utype),
			Spec: spec,
		},
		Status: "",
	}
}

func (e *Executor) getUpstreamCrdStatus(name, namespace string) string {
	o, err := e.client.GlueV1().Upstreams(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(o.Status)
}

func (e *Executor) wait(w int, cb func(e *Executor) bool) error {
	if w <= 0 {
		return nil
	}
	for i := 0; i < w; i++ {
		if cb(e) {
			return nil
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("Wait timeout")
}
