package upstream

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/solo-io/gloo/pkg/coreplugins/service"

	"github.com/solo-io/gloo/pkg/storage/dependencies"

	"github.com/solo-io/gloo/pkg/protoutil"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/plugins/aws"
	"github.com/solo-io/gloo/pkg/plugins/google"
	"github.com/solo-io/gloo/pkg/storage"
	psecret "github.com/solo-io/glooctl/pkg/secret"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type upstreamMatcher func(*v1.Upstream) bool
type upstreamEditor func(storage.Interface, dependencies.SecretStorage, *v1.Upstream) error

type plugin struct {
	name    string
	matcher upstreamMatcher
	editor  upstreamEditor
}

var (
	upstreamPlugins = []plugin{
		{"AWS", typeBasedMatcher(aws.UpstreamTypeAws), awsInteractive},
		{"Google", typeBasedMatcher(gfunc.UpstreamTypeGoogle), googleInteractive},
		{"Service", typeBasedMatcher(service.UpstreamTypeService), serviceInteractive},
	}

	// name regex
	nameRegex = regexp.MustCompile(`^[a-z][a-z0-9\-\.]{0,252}$`)
)

func typeBasedMatcher(t string) upstreamMatcher {
	return func(u *v1.Upstream) bool {
		if u == nil {
			return false
		}
		return u.Type == t
	}
}

func SelectInteractive(sc storage.Interface) (*v1.Upstream, error) {
	existing, err := sc.V1().Upstreams().List()
	if err != nil {
		return nil, err
	}
	if len(existing) == 0 {
		return nil, errors.New("no existing upstreams to update")
	}
	upstreamNames := make([]string, len(existing))
	for i, u := range existing {
		upstreamNames[i] = u.Name
	}

	var selected string
	if err := survey.AskOne(&survey.Select{
		Message: "Please select the upstream to edit:",
		Options: upstreamNames,
	}, &selected, survey.Required); err != nil {
		return nil, err
	}

	for _, u := range existing {
		if u.Name == selected {
			return u, nil
		}
	}
	return nil, errors.New("didn't find selected upstream")
}

// Interactive - create/update upstream interactively
func Interactive(sc storage.Interface, si dependencies.SecretStorage, u *v1.Upstream) error {
	err := askName(sc, u)
	if err != nil {
		return err
	}
	// type
	editor, err := upstreamType(u)
	if err != nil {
		return err
	}

	err = editor(sc, si, u)
	if err != nil {
		return nil
	}

	return askConnectionTimeout(u)
}

func upstreamType(u *v1.Upstream) (upstreamEditor, error) {
	upstreamTypes := make([]string, len(upstreamPlugins))
	for i, u := range upstreamPlugins {
		upstreamTypes[i] = u.name
	}
	if u.Type == "" {
		question := &survey.Select{
			Message: "Select the type of upstream to create:",
			Options: upstreamTypes,
		}
		var choice string
		if err := survey.AskOne(question, &choice, survey.Required); err != nil {
			return nil, err
		}
		for _, t := range upstreamPlugins {
			if choice == t.name {
				return t.editor, nil
			}
		}
		return nil, errors.New("did not find an upstream editor")
	}

	for _, t := range upstreamPlugins {
		if t.matcher(u) {
			return t.editor, nil
		}
	}
	return nil, errors.Errorf("Upstream type %s is not supported by interactive mode", u.Type)
}

func askName(sc storage.Interface, u *v1.Upstream) error {
	if u.Name != "" {
		// we are updating an existing upstream so nothing to do
		return nil
	}
	// new upstream
	upstreams, err := sc.V1().Upstreams().List()
	if err != nil {
		return errors.Wrap(err, "unable to get list of upstreams")
	}
	// name
	return survey.AskOne(&survey.Input{Message: "Please enter a name for the upstream:"}, &u.Name, func(val interface{}) error {
		name, ok := val.(string)
		if !ok {
			return errors.New("expecting a string for name")
		}
		// is unique name
		for _, u := range upstreams {
			if name == u.Name {
				return errors.New("upstream name needs to be unique")
			}
		}
		// check format - https://kubernetes.io/docs/concepts/overview/working-with-objects/names/
		if !nameRegex.MatchString(name) {
			return errors.New("upstream name format only supports letters, '\\' and '.'")
		}
		return nil
	})
}

func askConnectionTimeout(u *v1.Upstream) error {
	var defaultTimeout string
	if u.ConnectionTimeout > 0 {
		defaultTimeout = strconv.FormatInt(int64(u.ConnectionTimeout/time.Second), 10)
	}
	var timeout int
	err := survey.AskOne(&survey.Input{
		Message: "Please enter connection timeout in seconds:",
		Default: defaultTimeout,
	}, &timeout, func(val interface{}) error {
		_, errTimeout := strconv.Atoi(val.(string))
		if errTimeout != nil {
			return errors.New("timeout must be a positive integer")
		}
		return nil
	})
	if err != nil {
		return err
	}
	u.ConnectionTimeout = time.Duration(timeout) * time.Second
	return nil
}

func awsInteractive(sc storage.Interface, si dependencies.SecretStorage, u *v1.Upstream) error {
	u.Type = aws.UpstreamTypeAws
	regions := make([]string, len(aws.ValidRegions))
	i := 0
	for k := range aws.ValidRegions {
		regions[i] = k
		i++
	}

	var existingRegion string
	var existingSecretRef string
	if u.Spec != nil {
		spec, err := aws.DecodeUpstreamSpec(u.Spec)
		if err == nil {
			existingRegion = spec.Region
			existingSecretRef = spec.SecretRef
		}
	}

	var region string
	if err := survey.AskOne(&survey.Select{
		Message: "Please select an AWS region",
		Options: regions,
		Default: existingRegion,
	}, &region, survey.Required); err != nil {
		return err
	}

	secrets, err := psecret.SecretRefs(si, isAWSSecret)
	if err != nil {
		return err
	}
	if len(secrets) == 0 {
		return errors.New("unable to find any AWS secret")
	}
	var secretRef string
	if err := survey.AskOne(&survey.Select{
		Message: "Please select an AWS secret for the upstream:",
		Options: secrets,
		Default: existingSecretRef,
	}, &secretRef, survey.Required); err != nil {
		return err
	}
	u.Spec = aws.EncodeUpstreamSpec(aws.UpstreamSpec{Region: region, SecretRef: secretRef})
	return nil
}

func isAWSSecret(s *dependencies.Secret) bool {
	if s.Data == nil {
		return false
	}

	_, first := s.Data[aws.AwsAccessKey]
	_, second := s.Data[aws.AwsSecretKey]
	return first && second
}

func googleInteractive(sc storage.Interface, si dependencies.SecretStorage, u *v1.Upstream) error {
	u.Type = gfunc.UpstreamTypeGoogle

	regions := make([]string, len(gfunc.ValidRegions))
	i := 0
	for k := range gfunc.ValidRegions {
		regions[i] = k
		i++
	}

	var existingRegion string
	var existingProject string

	if u.Spec != nil {
		spec, err := gfunc.DecodeUpstreamSpec(u.Spec)
		if err == nil {
			existingRegion = spec.Region
			existingProject = spec.ProjectId
		}
	}

	var region string
	if err := survey.AskOne(&survey.Select{
		Message: "Please select Google Cloud Function region",
		Options: regions,
		Default: existingRegion,
	}, &region, survey.Required); err != nil {
		return err
	}

	var project string
	if err := survey.AskOne(&survey.Input{
		Message: "Please enter the project ID:",
		Default: existingProject,
	}, &project, survey.Required); err != nil { // add better validation for project id
		return err
	}

	spec, err := protoutil.MarshalStruct(gfunc.UpstreamSpec{Region: region, ProjectId: project})
	if err != nil {
		return err
	}
	u.Spec = spec

	discovery := false
	if err := survey.AskOne(&survey.Confirm{Message: "Do you want to enable function discovery?"}, &discovery, nil); err != nil {
		return err
	}
	if discovery {
		secrets, err := psecret.SecretRefs(si, isGoogleSecret)
		if err != nil {
			return err
		}
		var existingSecretRef string
		if u.Metadata != nil && u.Metadata.Annotations != nil {
			ref, ok := u.Metadata.Annotations[psecret.GoogleAnnotationKey]
			if ok {
				existingSecretRef = ref
			}

		}
		var secretRef string
		if err := survey.AskOne(&survey.Select{
			Message: "Please select a Google secret for the upstream:",
			Options: secrets,
			Default: existingSecretRef,
		}, &secretRef, survey.Required); err != nil {
			return err
		}

		if u.Metadata == nil {
			u.Metadata = &v1.Metadata{}
		}
		if u.Metadata.Annotations == nil {
			u.Metadata.Annotations = make(map[string]string)
		}
		u.Metadata.Annotations[psecret.GoogleAnnotationKey] = secretRef
	}
	return nil
}

func isGoogleSecret(s *dependencies.Secret) bool {
	if s.Data == nil {
		return false
	}
	_, ok := s.Data[psecret.ServiceAccountJsonKeyFile]
	return ok
}

func serviceInteractive(sc storage.Interface, si dependencies.SecretStorage, u *v1.Upstream) error {
	u.Type = service.UpstreamTypeService

	var existingHosts []service.Host
	if u.Spec != nil {
		spec, err := service.DecodeUpstreamSpec(u.Spec)
		if err == nil {
			existingHosts = spec.Hosts
		}
	}

	if len(existingHosts) != 0 {
		printHosts(existingHosts)
		replace := false
		if err := survey.AskOne(&survey.Confirm{Message: "Do you want to replace existing host(s)?"}, &replace, nil); err != nil {
			return err
		}
		if !replace {
			return nil
		}
	}
	var hosts []service.Host
	add := true
	for add {
		questions := []*survey.Question{
			{
				Name:     "addr",
				Prompt:   &survey.Input{Message: "Please enter the service host address:"},
				Validate: survey.Required,
			},
			{
				Name:     "port",
				Prompt:   &survey.Input{Message: "Please enter the service host port:"},
				Validate: validatePort,
			},
		}
		host := service.Host{}
		if err := survey.Ask(questions, &host); err != nil {
			return err
		}
		hosts = append(hosts, host)
		printHosts(hosts)
		if err := survey.AskOne(&survey.Confirm{Message: "Do you want to add more hosts?"}, &add, nil); err != nil {
			return err
		}
	}

	spec, err := protoutil.MarshalStruct(service.UpstreamSpec{Hosts: hosts})
	if err != nil {
		return err
	}
	u.Spec = spec
	return nil
}

func printHosts(list []service.Host) {
	fmt.Println("Hosts")
	for i, h := range list {
		fmt.Printf("%2d: %s:%d\n", (i + 1), h.Addr, h.Port)
	}
	fmt.Printf("\n\n")
}

// validators
func validatePort(val interface{}) error {
	v, ok := val.(string)
	if !ok {
		return errors.New("unable to convert value for validation")
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return errors.Wrap(err, "unable to convert into a number")
	}
	if i <= 0 || i > 65535 {
		return errors.New("invalid port number")
	}
	return nil
}
