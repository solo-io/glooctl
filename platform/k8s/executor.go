package k8s

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
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
		Fatal("Cannot create k8s client", err)
	}
	client, err := crdclient.NewForConfig(cfg)
	if err != nil {
		Fatal("Cannot create glue CRDs clientset", err)
	}
	return &Executor{
		cfg:    cfg,
		client: client,
	}
}

func (e *Executor) RunCreateUpstream(gparams *platform.GlobalParams, uparams *platform.UpstreamParams) {
	e.updateUpstream(gparams, uparams, true)
}

func (e *Executor) RunUpdateUpstream(gparams *platform.GlobalParams, uparams *platform.UpstreamParams) {
	e.updateUpstream(gparams, uparams, false)
}

func (e *Executor) RunDeleteUpstream(gparams *platform.GlobalParams, uparams *platform.UpstreamParams) {
	if uparams.Name == "" {
		Fatal("Name of the Upstream must be provided")
	}
	err := e.client.GlueV1().Upstreams(gparams.Namespace).Delete(uparams.Name, &metav1.DeleteOptions{})
	if err != nil {
		Fatal(err)
	}
	err = e.wait(gparams.WaitSec, func(e *Executor) bool {
		s := e.getUpstreamCrdStatus(uparams.Name, gparams.Namespace, false)
		if s != "" {
			fmt.Printf("Upstream Status: %s\n", s)
			return true
		}
		return false
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Upstream deleted")
	}
}

func (e *Executor) RunGetUpstream(gparams *platform.GlobalParams, uparams *platform.UpstreamParams) {
	e.getUpstream(gparams, uparams, false)
}

func (e *Executor) RunDescribeUpstream(gparams *platform.GlobalParams, uparams *platform.UpstreamParams) {
	e.getUpstream(gparams, uparams, true)
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

func (e *Executor) getUpstreamCrdStatus(name, namespace string, ignoreErr bool) string {
	o, err := e.client.GlueV1().Upstreams(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if ignoreErr {
			return ""
		} else {
			return err.Error()
		}

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

func (e *Executor) updateUpstream(gparams *platform.GlobalParams, uparams *platform.UpstreamParams, isCreate bool) {

	if uparams.Name == "" || uparams.UType == "" {
		Fatal("Both Name and Type of the Upstream must be provided")
	}

	x := upstreamFromArgs(uparams.Name, uparams.UType, uparams.Spec)
	if isCreate {
		_, err := e.client.GlueV1().Upstreams(gparams.Namespace).Create(x)
		if err != nil {
			Fatal(err)
		}
	} else {
		o, err := e.client.GlueV1().Upstreams(gparams.Namespace).Get(uparams.Name, metav1.GetOptions{})
		if err != nil {
			Fatal(err)
		}
		x.ObjectMeta = o.ObjectMeta
		_, err = e.client.GlueV1().Upstreams(gparams.Namespace).Update(x)
		if err != nil {
			Fatal(err)
		}
	}
	err := e.wait(gparams.WaitSec, func(e *Executor) bool {
		s := e.getUpstreamCrdStatus(uparams.Name, gparams.Namespace, true)
		if s != "" {
			fmt.Printf("Upstream Status: %s\n", s)
			return true
		}
		return false
	})
	if err != nil {
		fmt.Println(err)
	} else {
		if isCreate {
			fmt.Println("Upstream created")
		} else {
			fmt.Println("Upstream updated")
		}
	}
}

func (e *Executor) getUpstream(gparams *platform.GlobalParams, uparams *platform.UpstreamParams, isDescribe bool) {
	var w *tabwriter.Writer
	if !isDescribe {
		w = tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "\n NAME\t NAMESPACE\t TYPE\t STATUS")
	}

	if uparams.Name == "" {
		// List
		ll, err := e.client.GlueV1().Upstreams(gparams.Namespace).List(metav1.ListOptions{})
		if err != nil {
			Fatal(err)
		}
		for _, o := range ll.Items {
			e.printUpstream(&o, isDescribe, w)
		}
	} else {
		// Single
		o, err := e.client.GlueV1().Upstreams(gparams.Namespace).Get(uparams.Name, metav1.GetOptions{})
		if err != nil {
			Fatal(err)
		}
		e.printUpstream(o, isDescribe, w)
	}
	if !isDescribe {
		w.Flush()
	}
}

func (e *Executor) printUpstream(o *v1.Upstream, isDescribe bool, w *tabwriter.Writer) {
	if isDescribe {
		x, err := json.MarshalIndent(o, "", "  ")
		if err != nil {
			fmt.Println(o)
		}
		fmt.Println(string(x))
	} else {
		fmt.Fprintf(w, " %s \t %s \t %s \t %s\n", o.Name, o.Namespace, o.Spec.Type, o.Status)
	}
}

func Fatal(x ...interface{}) {
	fmt.Println("\nERROR: ", x)
	os.Exit(1)
}
