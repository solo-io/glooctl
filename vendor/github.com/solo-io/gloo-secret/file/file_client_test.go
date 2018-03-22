package file

import (
	"io/ioutil"
	"os"
	"testing"

	secret "github.com/solo-io/gloo-secret"
)

func TestInvalidDirFails(t *testing.T) {
	_, e := NewClient("doesnotexist")
	if e == nil {
		t.Error("not existing directory should have returned error")
	}

	_, e = NewClient("file_client.go")
	if e == nil {
		t.Error("non directory should have returned error")
	}
}

func TestCreateAndGet(t *testing.T) {
	dir, err := ioutil.TempDir("", "gloo-secret-test")
	if err != nil {
		t.Errorf("unable to get temporary directory %q", err)
	}
	defer os.RemoveAll(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Errorf("Unable to setup client %q", err)
	}
	secrets, err := client.V1().List()
	if err != nil {
		t.Errorf("unable to list with empty repository %q", err)
	}
	if len(secrets) != 0 {
		t.Errorf("was expecting empty repository")
	}

	secret := &secret.Secret{Name: "test.yaml", Data: map[string][]byte{
		"user":     []byte("hello"),
		"password": []byte("secret"),
	}}
	_, err = client.V1().Create(secret)
	if err != nil {
		t.Errorf("unable to create secrets %q", err)
	}

	loaded, err := client.V1().Get("test.yaml")
	if err != nil {
		t.Errorf("unable to get created secret %q", err)
	}
	if secret.Name != loaded.Name {
		t.Errorf("name do not match; expected %s got %s", secret.Name, loaded.Name)
	}

	if len(secret.Data) != len(loaded.Data) {
		t.Errorf("number of entries do not match; expected %d got %d", len(secret.Data), len(loaded.Data))
	}

	for k, v := range secret.Data {
		if string(v) != string(loaded.Data[k]) {
			t.Errorf("value did not match for key %s; expected %s got %s", k, string(v), string(loaded.Data[k]))
		}
	}
}
func TestCreateAndList(t *testing.T) {
	dir, err := ioutil.TempDir("", "gloo-secret-test")
	if err != nil {
		t.Errorf("unable to get temporary directory %q", err)
	}
	defer os.RemoveAll(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Errorf("Unable to setup client %q", err)
	}
	secrets, err := client.V1().List()
	if err != nil {
		t.Errorf("unable to list with empty repository %q", err)
	}
	if len(secrets) != 0 {
		t.Errorf("was expecting empty repository")
	}

	secret := &secret.Secret{Name: "test.yaml", Data: map[string][]byte{
		"user":     []byte("hello"),
		"password": []byte("secret"),
	}}
	_, err = client.V1().Create(secret)
	if err != nil {
		t.Errorf("unable to create secrets %q", err)
	}

	secrets, err = client.V1().List()
	if err != nil {
		t.Errorf("unable to list after creating secret")
	}
	if len(secrets) != 1 {
		t.Errorf("expected one secret")
	}
}

func TestDelete(t *testing.T) {
	dir, err := ioutil.TempDir("", "gloo-secret-test")
	if err != nil {
		t.Errorf("unable to get temporary directory %q", err)
	}
	defer os.RemoveAll(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Errorf("Unable to setup client %q", err)
	}
	secret := &secret.Secret{Name: "test.yaml", Data: map[string][]byte{
		"user":     []byte("hello"),
		"password": []byte("secret"),
	}}
	_, err = client.V1().Create(secret)
	if err != nil {
		t.Errorf("unable to create secrets %q", err)
	}

	_, err = client.V1().Get("test.yaml")
	if err != nil {
		t.Errorf("unable to get after creating secret")
	}

	err = client.V1().Delete("test.yaml")
	if err != nil {
		t.Errorf("unable to delete secret")
	}

	_, err = client.V1().Get("test.yaml")
	if err == nil {
		t.Errorf("shouldn't be able to get deleted secret")
	}
}

func TestUpdate(t *testing.T) {
	dir, err := ioutil.TempDir("", "gloo-secret-test")
	if err != nil {
		t.Errorf("unable to get temporary directory %q", err)
	}
	defer os.RemoveAll(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Errorf("Unable to setup client %q", err)
	}
	secret := &secret.Secret{Name: "test.yaml", Data: map[string][]byte{
		"user":     []byte("hello"),
		"password": []byte("secret"),
	}}
	_, err = client.V1().Create(secret)
	if err != nil {
		t.Errorf("unable to create secrets %q", err)
	}

	existing, err := client.V1().Get("test.yaml")
	if err != nil {
		t.Errorf("unable to get after creating secret")
	}

	existing.Data["domain"] = []byte("axhixh.com")
	_, err = client.V1().Update(existing)
	if err != nil {
		t.Errorf("unable to update secret")
	}

	postUpdate, err := client.V1().Get("test.yaml")
	if err != nil {
		t.Errorf("error getting updated secret")
	}

	if string(postUpdate.Data["domain"]) != "axhixh.com" {
		t.Errorf("updated information not saved")
	}
}
