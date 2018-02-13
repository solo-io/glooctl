package crd_test

import (
	"os"
	"path"
	"testing"

	gluev1 "github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/gluectl/pkg/storage"
	"github.com/solo-io/gluectl/pkg/storage/crd"
	"k8s.io/client-go/tools/clientcmd"
)

func TestUpstream(t *testing.T) {

	s, err := GetCrdStorage()
	if err != nil {
		t.Fatal("GetClient failed", err)
	}
	x := &gluev1.Upstream{Name: "testcrd", Type: "aws"}
	_, err = s.Create(x)
	if err != nil {
		t.Fatal("Create failed", err)
	}
	spec := make(map[string]interface{})
	spec["region"] = "us-east-1"
	spec["secret"] = "my secret"
	x.Spec = spec
	_, err = s.Update(x)
	if err != nil {
		t.Fatal("Update failed", err)
	}
	y, err := s.Get(&gluev1.Upstream{Name: "testcrd"}, nil)
	if err != nil {
		t.Fatal("Get failed", err)
	}
	t.Log(y)
	err = s.Delete(y)
	if err != nil {
		t.Fatal("Delete failed", err)
	}
}

func TestList(t *testing.T) {
	s, err := GetCrdStorage()
	if err != nil {
		t.Fatal("GetClient failed", err)
	}
	x, err := s.List(&gluev1.Upstream{}, nil)
	for _, z := range x {
		y, ok := z.(*gluev1.Upstream)
		if !ok {
			t.Fatal("List failed - type assertion")
		}
		t.Log(y)
	}
}

func TestUpstreamWatch(t *testing.T) {
	s, err := GetCrdStorage()
	if err != nil {
		t.Fatal("GetClient failed", err)
	}
	w, err := s.Watch(&gluev1.Upstream{}, nil)
	if err != nil {
		t.Fatal("Watch failed", err)
	}
	go func() {
		for i := 0; i < 3; i++ {
			c := <-w.ResultChan()
			t.Log(c)
		}
	}()
	TestUpstream(t)
	w.Stop()
}

func GetCrdStorage() (storage.Storage, error) {
	kubecfg := path.Join(os.Getenv("HOME"), ".kube/config")
	cfg, err := clientcmd.BuildConfigFromFlags("", kubecfg)
	if err != nil {
		return nil, err
	}
	return crd.NewCrdStorage(cfg, "ant")
}
