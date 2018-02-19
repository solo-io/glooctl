package executor

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	google_protobuf "github.com/gogo/protobuf/types"
	gloov1 "github.com/solo-io/gloo-api/pkg/api/types/v1"
	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/glooctl/platform"
)

type UpstreamExecutor struct {
	store storage.Interface
}

func NewUpstreamExecutor(store storage.Interface) platform.Executor {

	return &UpstreamExecutor{
		store: store,
	}
}

func (e *UpstreamExecutor) RunCreate(gparams *platform.GlobalParams, params interface{}) {
	e.updateUpstream(gparams, getUParams(params), true)
}

func (e *UpstreamExecutor) RunUpdate(gparams *platform.GlobalParams, params interface{}) {
	e.updateUpstream(gparams, getUParams(params), false)
}

func (e *UpstreamExecutor) RunDelete(gparams *platform.GlobalParams, params interface{}) {
	uparams := getUParams(params)
	if uparams.Name == "" {
		Fatal("Name of the Upstream must be provided")
	}
	err := e.store.V1().Upstreams().Delete(uparams.Name)
	if err != nil {
		Fatal(err)
	}
	err = e.wait(gparams.WaitSec, func(e *UpstreamExecutor) bool {
		s := e.getUpstreamStatus(uparams.Name, gparams.Namespace, false)
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

func (e *UpstreamExecutor) RunGet(gparams *platform.GlobalParams, params interface{}) {
	e.getUpstream(gparams, getUParams(params), false)
}

func (e *UpstreamExecutor) RunDescribe(gparams *platform.GlobalParams, params interface{}) {
	e.getUpstream(gparams, getUParams(params), true)
}

func (e *UpstreamExecutor) getUpstreamStatus(name, namespace string, ignoreErr bool) string {
	_, err := e.store.V1().Upstreams().Get(name)
	if err != nil {
		if ignoreErr {
			return ""
		} else {
			return err.Error()
		}
	}
	// TODO: get status
	return "ok"
}

func (e *UpstreamExecutor) wait(w int, cb func(e *UpstreamExecutor) bool) error {
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

func (e *UpstreamExecutor) updateUpstream(gparams *platform.GlobalParams, uparams *platform.UpstreamParams, isCreate bool) {

	if uparams.Name == "" || uparams.UType == "" {
		Fatal("Both Name and Type of the Upstream must be provided")
	}

	x := &gloov1.Upstream{
		Name: uparams.Name,
		Type: uparams.UType,
		Spec: &google_protobuf.Struct{Fields: uparams.Spec},
	}
	if isCreate {
		_, err := e.store.V1().Upstreams().Create(x)
		if err != nil {
			Fatal(err)
		}
	} else {
		_, err := e.store.V1().Upstreams().Update(x)
		if err != nil {
			Fatal(err)
		}
	}
	err := e.wait(gparams.WaitSec, func(e *UpstreamExecutor) bool {
		s := e.getUpstreamStatus(uparams.Name, gparams.Namespace, true)
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

func (e *UpstreamExecutor) getUpstream(gparams *platform.GlobalParams, uparams *platform.UpstreamParams, isDescribe bool) {
	var w *tabwriter.Writer
	if !isDescribe {
		w = tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "\n NAME\t TYPE")
	}

	if uparams.Name == "" {
		// List
		ll, err := e.store.V1().Upstreams().List()
		if err != nil {
			Fatal(err)
		}
		for _, o := range ll {
			e.printUpstream(o, isDescribe, w)
		}
	} else {
		// Single
		o, err := e.store.V1().Upstreams().Get(uparams.Name)
		if err != nil {
			Fatal(err)
		}
		e.printUpstream(o, isDescribe, w)
	}
	if !isDescribe {
		w.Flush()
	}
}

func (e *UpstreamExecutor) printUpstream(o *gloov1.Upstream, isDescribe bool, w *tabwriter.Writer) {
	if isDescribe {
		x, err := json.MarshalIndent(o, "", "  ")
		if err != nil {
			fmt.Println(o)
		}
		fmt.Println(string(x))
	} else {
		fmt.Fprintf(w, " %s \t %s\n", o.Name, o.Type)
	}
}
