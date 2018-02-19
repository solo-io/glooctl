package executor

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	gloov1 "github.com/solo-io/gloo-api/pkg/api/types/v1"
	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/glooctl/platform"
)

type VhostExecutor struct {
	store storage.Interface
}

func NewVhostExecutor(store storage.Interface) platform.Executor {

	return &VhostExecutor{
		store: store,
	}
}

func (e *VhostExecutor) RunCreate(gparams *platform.GlobalParams, params interface{}) {
	e.updateVhost(gparams, getVParams(params), true)
}

func (e *VhostExecutor) RunUpdate(gparams *platform.GlobalParams, params interface{}) {
	e.updateVhost(gparams, getVParams(params), false)
}

func (e *VhostExecutor) RunDelete(gparams *platform.GlobalParams, params interface{}) {
	vparams := getVParams(params)
	if vparams.Name == "" {
		Fatal("Name of the Vhost must be provided")
	}
	err := e.store.V1().VirtualHosts().Delete(vparams.Name)
	if err != nil {
		Fatal(err)
	}
	err = e.wait(gparams.WaitSec, func(e *VhostExecutor) bool {
		s := e.getVhostStatus(vparams.Name, gparams.Namespace, false)
		if s != "" {
			fmt.Printf("Vhost Status: %s\n", s)
			return true
		}
		return false
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Vhost deleted")
	}
}

func (e *VhostExecutor) RunGet(gparams *platform.GlobalParams, params interface{}) {
	e.getVhost(gparams, getVParams(params), false)
}

func (e *VhostExecutor) RunDescribe(gparams *platform.GlobalParams, params interface{}) {
	e.getVhost(gparams, getVParams(params), true)
}

func (e *VhostExecutor) getVhostStatus(name, namespace string, ignoreErr bool) string {
	_, err := e.store.V1().VirtualHosts().Get(name)
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

func (e *VhostExecutor) wait(w int, cb func(e *VhostExecutor) bool) error {
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

func (e *VhostExecutor) updateVhost(gparams *platform.GlobalParams, uparams *platform.VhostParams, isCreate bool) {

	if uparams.Name == "" {
		Fatal("Name of the Vhost must be provided")
	}

	x := &gloov1.VirtualHost{
		Name: uparams.Name,
	}
	if isCreate {
		_, err := e.store.V1().VirtualHosts().Create(x)
		if err != nil {
			Fatal(err)
		}
	} else {
		_, err := e.store.V1().VirtualHosts().Update(x)
		if err != nil {
			Fatal(err)
		}
	}
	err := e.wait(gparams.WaitSec, func(e *VhostExecutor) bool {
		s := e.getVhostStatus(uparams.Name, gparams.Namespace, true)
		if s != "" {
			fmt.Printf("Vhost Status: %s\n", s)
			return true
		}
		return false
	})
	if err != nil {
		fmt.Println(err)
	} else {
		if isCreate {
			fmt.Println("Vhost created")
		} else {
			fmt.Println("Vhost updated")
		}
	}
}

func (e *VhostExecutor) getVhost(gparams *platform.GlobalParams, uparams *platform.VhostParams, isDescribe bool) {
	var w *tabwriter.Writer
	if !isDescribe {
		w = tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "\n NAME\t TYPE")
	}

	if uparams.Name == "" {
		// List
		ll, err := e.store.V1().VirtualHosts().List()
		if err != nil {
			Fatal(err)
		}
		for _, o := range ll {
			e.printVhost(o, isDescribe, w)
		}
	} else {
		// Single
		o, err := e.store.V1().VirtualHosts().Get(uparams.Name)
		if err != nil {
			Fatal(err)
		}
		e.printVhost(o, isDescribe, w)
	}
	if !isDescribe {
		w.Flush()
	}
}

func (e *VhostExecutor) printVhost(o *gloov1.VirtualHost, isDescribe bool, w *tabwriter.Writer) {
	if isDescribe {
		x, err := json.MarshalIndent(o, "", "  ")
		if err != nil {
			fmt.Println(o)
		}
		fmt.Println(string(x))
	} else {
		fmt.Fprintf(w, " %s \t %s\n", o.Name, "boom")
	}
}
