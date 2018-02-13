package storage

import "k8s.io/apimachinery/pkg/watch"

const (
	Create = iota
	Update
	Delete
)

type WatchOperation int

type Item interface{}

type GetOptions struct{}
type ListOptions struct{}
type WatchOptions struct{}

type Storage interface {
	Create(item Item) (Item, error)
	Update(item Item) (Item, error)
	Delete(item Item) error
	Get(item Item, getOptions *GetOptions) (Item, error)
	List(item Item, listOptions *ListOptions) ([]Item, error)
	Watch(item Item, watchOptions *WatchOptions) (watch.Interface, error)
}
