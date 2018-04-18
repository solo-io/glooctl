package upstream

import (
	"regexp"
	"strconv"
	"time"

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

type upstreamCreator func(storage.Interface, dependencies.SecretStorage, *v1.Upstream) error

type plugin struct {
	name    string
	creator upstreamCreator
}

var (
	upstreamPlugins = []plugin{
		{"AWS", awsInteractive},
		//{"Azure", azureInteractive},
		//{"Consul", consulInteractive},
		{"Google", googleInteractive},
		//{"GRPC", grpcInteractive},
		//{"Kubernetes", kubeInteractive},
		//{"NATS streaming", natsInteractive},
		//{"REST service", restInteractive},
	}

	// name regex
	nameRegex = regexp.MustCompile(`^[a-z][a-z0-9\-\.]{0,252}$`)
)

func UpstreamInteractive(sc storage.Interface, si dependencies.SecretStorage) (*v1.Upstream, error) {
	upstreams, err := sc.V1().Upstreams().List()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of upstreams")
	}
	u := &v1.Upstream{}
	// name
	if err := survey.AskOne(&survey.Input{Message: "Please enter a name for the upstream:"}, &u.Name, func(val interface{}) error {
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
	}); err != nil {
		return nil, err
	}
	// type
	uType, err := upstreamType()
	if err != nil {
		return nil, err
	}
	// spec, service info
	for _, p := range upstreamPlugins {
		if p.name == uType {
			if err := p.creator(sc, si, u); err != nil {
				return nil, err
			}
		}
	}

	// connection timeout
	var timeout int
	survey.AskOne(&survey.Input{Message: "Please enter connection timeout in seconds:"}, &timeout, func(val interface{}) error {
		_, err := strconv.Atoi(val.(string))
		if err != nil {
			return errors.New("timeout must be a positive integer")
		}
		return nil
	})
	u.ConnectionTimeout = time.Duration(timeout) * time.Second

	// functions (separate separate interactions?)
	// metadata
	return u, nil
}

func upstreamType() (string, error) {
	upstreamTypes := make([]string, len(upstreamPlugins))
	for i, u := range upstreamPlugins {
		upstreamTypes[i] = u.name
	}
	question := &survey.Select{
		Message: "Select the type of upstream to create:",
		Options: upstreamTypes,
	}
	var choice string
	if err := survey.AskOne(question, &choice, survey.Required); err != nil {
		return "", err
	}
	return choice, nil
}

func secretRefs(si dependencies.SecretStorage, filter func(*dependencies.Secret) bool) ([]string, error) {
	secrets, err := si.List()
	if err != nil {
		return nil, err
	}
	var refs []string
	for _, s := range secrets {
		if filter(s) {
			refs = append(refs, s.Ref)
		}
	}
	return refs, nil
}

func awsInteractive(sc storage.Interface, si dependencies.SecretStorage, u *v1.Upstream) error {
	u.Type = aws.UpstreamTypeAws
	regions := make([]string, len(aws.ValidRegions))
	i := 0
	for k, _ := range aws.ValidRegions {
		regions[i] = k
		i++
	}

	var region string
	if err := survey.AskOne(&survey.Select{
		Message: "Please select an AWS region",
		Options: regions,
	}, &region, survey.Required); err != nil {
		return err
	}

	secrets, err := secretRefs(si, isAWSSecret)
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

func azureInteractive(sc storage.Interface, si dependencies.SecretStorage, u *v1.Upstream) error {
	return errors.New("not implemented")
}

func consulInteractive(sc storage.Interface, si dependencies.SecretStorage, u *v1.Upstream) error {
	return errors.New("not implemented")
}

func googleInteractive(sc storage.Interface, si dependencies.SecretStorage, u *v1.Upstream) error {
	u.Type = gfunc.UpstreamTypeGoogle

	regions := make([]string, len(gfunc.ValidRegions))
	i := 0
	for k, _ := range gfunc.ValidRegions {
		regions[i] = k
		i++
	}

	var region string
	if err := survey.AskOne(&survey.Select{
		Message: "Please select Google Cloud Function region",
		Options: regions,
	}, &region, survey.Required); err != nil {
		return err
	}

	var project string
	if err := survey.AskOne(&survey.Input{
		Message: "Please enter the project ID:",
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
		secrets, err := secretRefs(si, isGoogleSecret)
		if err != nil {
			return err
		}
		var secretRef string
		if err := survey.AskOne(&survey.Select{
			Message: "Please select a Google secret for the upstream:",
			Options: secrets,
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

func grpcInteractive(sc storage.Interface, si dependencies.SecretStorage, u *v1.Upstream) error {
	return errors.New("not implemented")
}

func kubeInteractive(sc storage.Interface, si dependencies.SecretStorage, u *v1.Upstream) error {
	return errors.New("not implemented")
}

func natsInteractive(sc storage.Interface, si dependencies.SecretStorage, u *v1.Upstream) error {
	return errors.New("not implemented")
}

func restInteractive(sc storage.Interface, si dependencies.SecretStorage, u *v1.Upstream) error {
	return errors.New("not implemented")
}
