package route

import (
	"strconv"
	"strings"
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
		{"CORS policy", corsPolicy},
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
	err = survey.AskOne(&survey.Input{Message: "Please enter maximum number of retries:"}, &retries, func(val interface{}) error {
		if _, errConvert := strconv.Atoi(val.(string)); errConvert != nil {
			return errors.New("maximum number of retries must be a positive integer")
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "unable to get maximum retries")
	}
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
	err = survey.AskOne(&survey.Input{Message: "Please enter request timeout in seconds:"}, &timeout, func(val interface{}) error {
		_, errConvert := strconv.Atoi(val.(string))
		if errConvert != nil {
			return errors.New("timeout must be a positive integer")
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "unable to get timetout")
	}
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
	err = survey.AskOne(&survey.Input{Message: "Please enter new host name to rewrite host header:"}, &hostRewrite, survey.Required)
	if err != nil {
		return errors.Wrap(err, "unable to get rewrite host header")
	}
	spec.HostRewrite = hostRewrite
	addAll(s, core.EncodeRouteExtensionSpec(spec))
	return nil
}

func corsPolicy(s *types.Struct) error {
	spec, err := core.DecodeRouteExtensions(s)
	if err != nil {
		return errors.Wrap(err, "unable to decode core route extension")
	}

	questions := []*survey.Question{
		{
			Name: "origins",
			Prompt: &survey.Input{
				Message: "Please enter allowed origin (use comma to separate if more than one): ",
			},
			Validate: survey.Required,
		},
		{
			Name:   "methods",
			Prompt: &survey.Input{Message: "Please enter allowed HTTP methods: "},
		},
		{
			Name:   "allowHeaders",
			Prompt: &survey.Input{Message: "Please enter allowed HTTP headers: "},
		},
		{
			Name:   "exposeHeaders",
			Prompt: &survey.Input{Message: "Please enter exposed HTTP headers: "},
		},
		{
			Name:   "maxAge",
			Prompt: &survey.Input{Message: "Please enter maximum age for cache (seconds): "},
			Validate: func(val interface{}) error {
				s, ok := val.(string)
				if !ok || s == "" {
					return nil
				}
				i, err := strconv.Atoi(s)
				if err != nil || i < 0 {
					return errors.New("maximum age must be a positive integer or empty")
				}
				return nil
			},
		},
		{
			Name:   "credentials",
			Prompt: &survey.Confirm{Message: "Allow requests with credentials? "},
		},
	}
	answers := struct {
		Origins       string
		Methods       string
		AllowHeaders  string
		ExposeHeaders string
		MaxAge        int
		Credentials   bool
	}{}
	if err := survey.Ask(questions, &answers); err != nil {
		return err
	}

	origins := strings.Split(answers.Origins, ",")
	for i, s := range origins {
		origins[i] = strings.TrimSpace(s)
	}
	spec.Cors = &core.CorsPolicy{
		AllowOrigin:      origins,
		AllowMethods:     answers.Methods,
		AllowHeaders:     answers.AllowHeaders,
		ExposeHeaders:    answers.ExposeHeaders,
		MaxAge:           time.Duration(answers.MaxAge) * time.Second,
		AllowCredentials: answers.Credentials,
	}

	addAll(s, core.EncodeRouteExtensionSpec(spec))
	return nil
}

func transformation(s *types.Struct) error {
	return errors.New("not implemented")
}
