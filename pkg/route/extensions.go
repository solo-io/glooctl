package route

import (
	"strconv"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	core "github.com/solo-io/gloo/pkg/coreplugins/route-extensions"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type extensionHandler func(*types.Struct) error

type extensionPlugin struct {
	name    string
	handler extensionHandler
}

var (
	availableExtensions = []extensionPlugin{
		{"Add request header", addRequestHeader},
		{"Add response header", addResponseHeader},
		{"Remove response header", removeResponseHeader},
		{"Set max retries", setMaxRetries},
		{"Set timeout", setTimeout},
		{"Rewrite host", rewriteHost},
		//{"Transformation", transformation},
	}
)

// addAll adds all the fields in right to left overwriting any keys in left
func addAll(left, right *types.Struct) {
	if left.Fields == nil {
		left.Fields = make(map[string]*types.Value)
	}
	for k, v := range right.Fields {
		left.Fields[k] = v
	}
}

func addRequestHeader(s *types.Struct) error {
	spec, err := core.DecodeRouteExtensions(s)
	if err != nil {
		return errors.Wrap(err, "unable to decode core route extension")
	}
	questions := []*survey.Question{
		{
			Name:     "key",
			Prompt:   &survey.Input{Message: "Please enter HTTP header name:"},
			Validate: survey.Required,
		},
		{
			Name:   "value",
			Prompt: &survey.Input{Message: "Please enter HTTP header value:"},
		},
		{
			Name:   "append",
			Prompt: &survey.Confirm{Message: "Append HTTP header"},
		},
	}
	header := core.HeaderValue{}
	if err := survey.Ask(questions, &header); err != nil {
		return err
	}
	spec.AddRequestHeaders = append(spec.AddRequestHeaders, header)
	addAll(s, core.EncodeRouteExtensionSpec(spec))
	return nil
}

func addResponseHeader(s *types.Struct) error {
	spec, err := core.DecodeRouteExtensions(s)
	if err != nil {
		return errors.Wrap(err, "unable to decode core route extension")
	}
	questions := []*survey.Question{
		{
			Name:     "key",
			Prompt:   &survey.Input{Message: "Please enter HTTP header name:"},
			Validate: survey.Required,
		},
		{
			Name:   "value",
			Prompt: &survey.Input{Message: "Please enter HTTP header value:"},
		},
		{
			Name:   "append",
			Prompt: &survey.Confirm{Message: "Append HTTP header"},
		},
	}
	header := core.HeaderValue{}
	if err := survey.Ask(questions, &header); err != nil {
		return err
	}
	spec.AddResponseHeaders = append(spec.AddResponseHeaders, header)
	addAll(s, core.EncodeRouteExtensionSpec(spec))
	return nil
}

func removeResponseHeader(s *types.Struct) error {
	spec, err := core.DecodeRouteExtensions(s)
	if err != nil {
		return errors.Wrap(err, "unable to decode core route extension")
	}
	var header string
	if err := survey.AskOne(&survey.Input{Message: "Please enter HTTP header name:"}, &header, survey.Required); err != nil {
		return err
	}
	spec.RemoveResponseHeaders = append(spec.RemoveResponseHeaders, header)
	addAll(s, core.EncodeRouteExtensionSpec(spec))
	return nil
}

func setMaxRetries(s *types.Struct) error {
	spec, err := core.DecodeRouteExtensions(s)
	if err != nil {
		return errors.Wrap(err, "unable to decode core route extension")
	}
	var retries int
	survey.AskOne(&survey.Input{Message: "Please enter maximum number of retries:"}, &retries, func(val interface{}) error {
		_, err := strconv.Atoi(val.(string))
		if err != nil {
			return errors.New("maximum number of retries must be a positive integer")
		}
		return nil
	})
	spec.MaxRetries = uint32(retries)
	addAll(s, core.EncodeRouteExtensionSpec(spec))
	return nil
}

func setTimeout(s *types.Struct) error {
	spec, err := core.DecodeRouteExtensions(s)
	if err != nil {
		return errors.Wrap(err, "unable to decode core route extension")
	}
	var timeout int
	survey.AskOne(&survey.Input{Message: "Please enter request timeout in seconds:"}, &timeout, func(val interface{}) error {
		_, err := strconv.Atoi(val.(string))
		if err != nil {
			return errors.New("timeout must be a positive integer")
		}
		return nil
	})
	spec.Timeout = time.Duration(timeout) * time.Second
	addAll(s, core.EncodeRouteExtensionSpec(spec))
	return nil
}

func rewriteHost(s *types.Struct) error {
	spec, err := core.DecodeRouteExtensions(s)
	if err != nil {
		return errors.Wrap(err, "unable to decode core route extension")
	}
	var hostRewrite string
	survey.AskOne(&survey.Input{Message: "Please enter new host name to rewrite host header:"}, &hostRewrite, survey.Required)
	spec.HostRewrite = hostRewrite
	addAll(s, core.EncodeRouteExtensionSpec(spec))
	return nil
}

func transformation(s *types.Struct) error {
	return errors.New("not implemented")
}
