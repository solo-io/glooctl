package util

import (
	"os"
	"path/filepath"
	"time"

	"log"

	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo-storage/crd"
	"github.com/solo-io/gloo-storage/file"
	"k8s.io/client-go/tools/clientcmd"
)

type StorageOptions struct {
	GlooConfigDir string
	SecretDir     string
	KubeConfig    string
	Namespace     string
	SyncPeriod    int
}

func GetStorageClient(opts *StorageOptions) (storage.Interface, error) {
	syncPeriod := time.Duration(opts.SyncPeriod) * time.Second
	if opts.GlooConfigDir != "" {
		log.Printf("Using file-based storage for gloo. Gloo must be configured to use file storage with config dir %v", opts.GlooConfigDir)
		return file.NewStorage(opts.GlooConfigDir, syncPeriod)
	}

	kubeConfig := opts.KubeConfig
	if kubeConfig == "" && HomeDir() != "" {
		kubeConfig = filepath.Join(HomeDir(), ".kube", "config")
	}
	kubeClient, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, err
	}
	return crd.NewStorage(kubeClient, opts.Namespace, syncPeriod)
}

func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
