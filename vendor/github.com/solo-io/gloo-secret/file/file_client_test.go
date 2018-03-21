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

	secret := &secret.Secret{Name: "test", Data: map[string][]byte{
		"user":     []byte("hello"),
		"password": []byte("secret"),
	}}
	created, err := client.V1().Create(secret)
	if err != nil {
		t.Errorf("unable to create secrets %q", err)
	}
	if created.ResourceVersion == "" {
		t.Errorf("expected non empty resource version")
	}

	secrets, err = client.V1().List()
	if err != nil {
		t.Errorf("unable to list after creating secret")
	}
	if len(secrets) != 1 {
		t.Errorf("expected one secret")
	}
}

func TestNonEmptyRepository(t *testing.T) {
	dir, err := ioutil.TempDir("", "gloo-secret-test")
	if err != nil {
		t.Errorf("unable to get temporary directory %q", err)
	}
	defer os.RemoveAll(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Errorf("Unable to setup client %q", err)
	}
	secret := &secret.Secret{Name: "test", Data: map[string][]byte{
		"user":     []byte("hello"),
		"password": []byte("secret"),
	}}
	_, err = client.V1().Create(secret)
	if err != nil {
		t.Errorf("unable to create secrets %q", err)
	}

	client, err = NewClient(dir)
	secrets, err := client.V1().List()
	if err != nil {
		t.Errorf("unable to list with non empty repository %q", err)
	}
	if len(secrets) != 1 {
		t.Errorf("was expecting non empty repository")
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
	secret := &secret.Secret{Name: "test", Data: map[string][]byte{
		"user":     []byte("hello"),
		"password": []byte("secret"),
	}}
	_, err = client.V1().Create(secret)
	if err != nil {
		t.Errorf("unable to create secrets %q", err)
	}

	_, err = client.V1().Get("test")
	if err != nil {
		t.Errorf("unable to get after creating secret")
	}

	err = client.V1().Delete("test")
	if err != nil {
		t.Errorf("unable to delete secret")
	}

	_, err = client.V1().Get("test")
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
	secret := &secret.Secret{Name: "test", Data: map[string][]byte{
		"user":     []byte("hello"),
		"password": []byte("secret"),
	}}
	_, err = client.V1().Create(secret)
	if err != nil {
		t.Errorf("unable to create secrets %q", err)
	}

	existing, err := client.V1().Get("test")
	if err != nil {
		t.Errorf("unable to get after creating secret")
	}

	existing.Data["domain"] = []byte("axhixh.com")
	_, err = client.V1().Update(existing)
	if err != nil {
		t.Errorf("unable to update secret")
	}

	postUpdate, err := client.V1().Get("test")
	if err != nil {
		t.Errorf("error getting updated secret")
	}

	if string(postUpdate.Data["domain"]) != "axhixh.com" {
		t.Errorf("updated information not saved")
	}
}
