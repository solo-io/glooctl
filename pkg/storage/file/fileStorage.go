package file

import (
	"fmt"
	"os"
	"path"

	"github.com/solo-io/gluectl/pkg/storage"
	"k8s.io/apimachinery/pkg/watch"
)

const (
	upstreamPath = "upstream"
	vhostPath    = "vhost"
)

type FileStorage struct {
	root, namespace string
}

func NewFileStorage(root, namespace string) (*FileStorage, error) {

	fullpath := path.Join(root, namespace)
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		err := os.MkdirAll(fullpath, 0777)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &FileStorage{root: root, namespace: namespace}, nil
}

func (c *FileStorage) Create(item storage.Item) (storage.Item, error) {
	/*
		if obj, ok := item.(*gluev1.Upstream); ok {
		} else if obj, ok := item.(*gluev1.VirtualHost); ok {
		}
	*/
	return nil, fmt.Errorf("Unknown Item Type: %t", item)
}

func (c *FileStorage) Update(item storage.Item) (storage.Item, error) {
	return nil, nil
}

func (c *FileStorage) Delete(item storage.Item) error {
	return nil
}

func (c *FileStorage) Get(item storage.Item, getOptions *storage.GetOptions) (storage.Item, error) {
	return nil, nil
}
func (c *FileStorage) List(item storage.Item, listOptions *storage.ListOptions) ([]storage.Item, error) {
	return nil, nil
}
func (c *FileStorage) Watch(item storage.Item, watchOptions *storage.WatchOptions) (watch.Interface, error) {
	return nil, nil
}
