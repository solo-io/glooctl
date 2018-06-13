package virtualservice

import (
	"github.com/solo-io/gloo/pkg/bootstrap/configstorage"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/storage/file"
	"github.com/solo-io/gloo/pkg/storage"
	"fmt"
)

func configure(opts *bootstrap.Options, filename string, overwriteExisting, deleteExisting bool) error {
	if filename == "" {
		return errors.Errorf("file not provided")
	}

	var cfg v1.Config
	if err := file.ReadFileInto(filename, &cfg); err != nil {
		return errors.Wrap(err, "reading file as v1.Config failed")
	}

	gloo, err := configstorage.Bootstrap(*opts)
	if err != nil {
		return errors.Wrap(err, "unable to create storage client")
	}

	return apply(gloo, &cfg, overwriteExisting, deleteExisting)
}

func apply(gloo storage.Interface, cfg *v1.Config, overwriteExisting, deleteExisting bool) error {
	actualUpstreams, err := gloo.V1().Upstreams().List()
	if err != nil {
		return err
	}
	actualVirtualServices, err := gloo.V1().VirtualServices().List()
	if err != nil {
		return err
	}
	if err := syncObjects(gloo, usToConfigObj(cfg.Upstreams), usToConfigObj(actualUpstreams), overwriteExisting, deleteExisting); err  != nil {
		return errors.Wrap(err, "syncing upstreams")
	}
	if err := syncObjects(gloo, vsToConfigObj(cfg.VirtualServices), vsToConfigObj(actualVirtualServices), overwriteExisting, deleteExisting); err  != nil {
		return errors.Wrap(err, "syncing virtualServices")
	}
	return nil
}

func usToConfigObj(upstreams []*v1.Upstream) []v1.ConfigObject {
	var items []v1.ConfigObject
	for _, us := range upstreams {
		items = append(items, us)
	}
	return items
}

func vsToConfigObj(virtualServices []*v1.VirtualService) []v1.ConfigObject {
	var items []v1.ConfigObject
	for _, us := range virtualServices {
		items = append(items, us)
	}
	return items
}


func syncObjects(gloo storage.Interface, desired, actual []v1.ConfigObject, overwriteExisting, deleteExisting bool) error {
	var (
		toCreate []v1.ConfigObject
		toUpdate []v1.ConfigObject
	)
	for _, desiredItem := range desired {
		var update bool
		for i, actualItem := range actual {
			if desiredItem.GetName() == actualItem.GetName() {
				if !overwriteExisting {
					return errors.Errorf("found existing config object with the name %v." +
						" to overwrite, enable overwriting with -w", desiredItem.GetName())
				}
				// set metadata if it's nil
				if desiredItem.GetMetadata() == nil {
					desiredItem.SetMetadata(&v1.Metadata{})
				}
				meta := desiredItem.GetMetadata()
				meta.ResourceVersion = actualItem.GetMetadata().ResourceVersion
				desiredItem.SetMetadata(meta)
				update = true
				toUpdate = append(toUpdate, desiredItem)
				// remove it from the list we match against
				actual = append(actual[:i], actual[i+1:]...)
				break
			}
		}
		if !update {
			// desired was not found, mark for creation
			toCreate = append(toCreate, desiredItem)
		}
	}
	for _, item := range toCreate {
		switch item := item.(type) {
		case *v1.VirtualService:
			if _, err := gloo.V1().VirtualServices().Create(item); err != nil {
				return fmt.Errorf("failed to create virtualService %s: %v", item.GetName(), err)
			}
		case *v1.Upstream:
			if _, err := gloo.V1().Upstreams().Create(item); err != nil {
				return fmt.Errorf("failed to create upstream %s: %v", item.GetName(), err)
			}
		default:
			return errors.Errorf("unsupported object type %v", item.GetName())
		}
		fmt.Printf("%v created\n", item.GetName())
	}
	for _, item := range toUpdate {
		switch item := item.(type) {
		case *v1.VirtualService:
			if _, err := gloo.V1().VirtualServices().Update(item); err != nil {
				return fmt.Errorf("failed to create virtualService %s: %v", item.GetName(), err)
			}
		case *v1.Upstream:
			if _, err := gloo.V1().Upstreams().Update(item); err != nil {
				return fmt.Errorf("failed to create upstream %s: %v", item.GetName(), err)
			}
		default:
			return errors.Errorf("unsupported object type %v", item.GetName())
		}
		fmt.Printf("%v updated\n", item.GetName())
	}
	if !deleteExisting {
		return nil
	}
	// only remaining are no longer desired, delete em!
	for _, item := range actual {
		switch item := item.(type) {
		case *v1.VirtualService:
			if err := gloo.V1().VirtualServices().Delete(item.GetName()); err != nil {
				return fmt.Errorf("failed to create virtualService %s: %v", item.GetName(), err)
			}
		case *v1.Upstream:
			if err := gloo.V1().Upstreams().Delete(item.GetName()); err != nil {
				return fmt.Errorf("failed to create upstream %s: %v", item.GetName(), err)
			}
		default:
			return errors.Errorf("unsupported object type %v", item.GetName())
		}
		fmt.Printf("%v deleted\n", item.GetName())
	}
	return nil
}
