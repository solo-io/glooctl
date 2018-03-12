package route

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo-storage/file"
)

func setupStorage() (storage.Interface, func(), error) {
	dir, err := ioutil.TempDir("", "glooctl-test")
	if err != nil {
		return nil, func() {}, fmt.Errorf("unable to get temporary directory %q", err)
	}
	cleanup := func() { os.RemoveAll(dir) }
	os.Mkdir(filepath.Join(dir, "virtualhosts"), 0777)
	sc, err := file.NewStorage(dir, time.Second)
	if err != nil {
		return nil, cleanup, fmt.Errorf("unable to get storage for testing: %q\n", err)
	}
	return sc, cleanup, nil
}
func TestCreateWithNoDefault(t *testing.T) {
	sc, cleanup, err := setupStorage()
	if err != nil {
		t.Errorf("unable to setup storage %q", err)
	}
	defer cleanup()
	route, err := fromRouteDetail(&routeDetail{pathPrefix: "/foo", upstream: "upstream"})
	if err != nil {
		t.Errorf("error creating route")
	}
	routes, err := runCreate(sc, "", "", route, false)
	if err != nil {
		t.Errorf("unable to create route %q\n", err)
	}
	if len(routes) != 1 {
		t.Errorf("expected one route but got %d instead", len(routes))
	}
}

func TestCreateWithExistingDefaultOfDifferentName(t *testing.T) {
	sc, cleanup, err := setupStorage()
	if err != nil {
		t.Errorf("unable to setup storage %q", err)
	}
	defer cleanup()

	vhost := &v1.VirtualHost{Name: "mydefault"}
	if _, err = sc.V1().VirtualHosts().Create(vhost); err != nil {
		t.Errorf("unable to create virtual host %q", err)
	}
	route, err := fromRouteDetail(&routeDetail{pathPrefix: "/foo", upstream: "upstream"})
	if err != nil {
		t.Errorf("error creating route")
	}
	routes, err := runCreate(sc, "", "", route, false)
	if err != nil {
		t.Errorf("unable to create route %q\n", err)
	}
	if len(routes) != 1 {
		t.Errorf("expected one route but got %d instead", len(routes))
	}

	// check it is on the existing virtual host
	v, err := sc.V1().VirtualHosts().Get("mydefault")
	if err != nil {
		t.Error("unable to get virtual host to validate", err)
	}
	if len(v.Routes) != 1 {
		t.Error("expecting 1 route got", len(v.Routes))
	}
}

func TestCreateAndSort(t *testing.T) {
	sc, cleanup, err := setupStorage()
	if err != nil {
		t.Errorf("unable to setup storage %q", err)
	}
	defer cleanup()

	prefixRoute, _ := fromRouteDetail(&routeDetail{pathPrefix: "/foo", upstream: "upstream"})
	vhost := &v1.VirtualHost{
		Name:   "default",
		Routes: []*v1.Route{prefixRoute}}
	if _, err = sc.V1().VirtualHosts().Create(vhost); err != nil {
		t.Errorf("unable to create virtual host %q", err)
	}
	newRoute, _ := fromRouteDetail(&routeDetail{pathExact: "/a", upstream: "upstream"})
	runCreate(sc, "", "", newRoute, true)

	// check it is on the existing virtual host
	v, err := sc.V1().VirtualHosts().Get("default")
	if err != nil {
		t.Error("unable to get virtual host to validate", err)
	}
	if len(v.Routes) != 2 {
		t.Error("expecting 2 route got", len(v.Routes))
	}
	if !v.Routes[0].Equal(newRoute) {
		t.Error("route not sorted correctly")
	}
}

func TestCreateWithExistingDomain(t *testing.T) {
	sc, cleanup, err := setupStorage()
	if err != nil {
		t.Errorf("unable to setup storage %q", err)
	}
	defer cleanup()

	vhost := &v1.VirtualHost{Name: "default"}
	if _, err = sc.V1().VirtualHosts().Create(vhost); err != nil {
		t.Errorf("unable to create virtual host %q", err)
	}
	vhost2 := &v1.VirtualHost{Name: "axhixh.com", Domains: []string{"axhixh.com"}}
	if _, err = sc.V1().VirtualHosts().Create(vhost2); err != nil {
		t.Errorf("unable to create virtual host 2 %q", err)
	}
	newRoute, _ := fromRouteDetail(&routeDetail{pathExact: "/a", upstream: "upstream"})
	runCreate(sc, "", "axhixh.com", newRoute, true)

	// check it is on the existing virtual host
	v, err := sc.V1().VirtualHosts().Get("axhixh.com")
	if err != nil {
		t.Error("unable to get virtual host to validate", err)
	}
	if len(v.Routes) != 1 {
		t.Error("expecting 1 route got", len(v.Routes))
	}
}

func TestCreateWithNonExistingDomain(t *testing.T) {
	sc, cleanup, err := setupStorage()
	if err != nil {
		t.Errorf("unable to setup storage %q", err)
	}
	defer cleanup()

	vhost := &v1.VirtualHost{Name: "default"}
	if _, err = sc.V1().VirtualHosts().Create(vhost); err != nil {
		t.Errorf("unable to create virtual host %q", err)
	}
	newRoute, _ := fromRouteDetail(&routeDetail{pathExact: "/a", upstream: "upstream"})
	_, err = runCreate(sc, "", "axhixh.com", newRoute, true)
	if err == nil {
		t.Errorf("should have error saying didn't find a virtual host")
	}
}
