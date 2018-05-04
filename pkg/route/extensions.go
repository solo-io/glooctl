package route

import (
	"strconv"
	"strings"
	"time"

	"github.com/solo-io/gloo/pkg/plugins/common/transformation"

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
		{"Add a new request header", addRequestHeader},
		{"Add a new response header", addResponseHeader},
		{"CORS policy", corsPolicy},
		{"Remove an existing response header", removeResponseHeader},
		{"Response Transformation", responseTransformation},
		{"Rewrite host", rewriteHost},
		{"Set max retries", setMaxRetries},
		{"Set timeout", setTimeout},
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

// responseTransformation shares the route extension struct with request
// transformation. We are changing just two fields of this struct
// TODO: add request transformation after changes in Gloo
func responseTransformation(s *types.Struct) error {
	spec, err := transformation.DecodeRouteExtension(s)
	if err != nil {
		return errors.Wrap(err, "unable to decode transformation route extension")
	}

	hasResponseParams := false
	err = survey.AskOne(&survey.Confirm{
		Message: "Do you want to define response parameters?",
		Help:    "Define custom parameters that are used by response template",
	}, &hasResponseParams, nil)
	if err != nil {
		return err
	}
	if hasResponseParams {
		path := ""
		err = survey.AskOne(&survey.Input{
			Message: "Please enter path based parameter (leave empty if you don't need one):",
			Help:    "Path based parameter helps you extract parameters from the URL path",
		}, &path, nil)
		if err != nil {
			return err
		}

		responseParam := transformation.Parameters{}
		if path != "" {
			responseParam.Path = &path
		}
		headers, errHeaders := askHeaders("Define a response header based parameter?", "Define another response header based parameter?")
		if errHeaders != nil {
			return errHeaders
		}
		if len(headers) != 0 {
			responseParam.Headers = headers
		}
		spec.ResponseParams = &responseParam
	}

	bodyTemplate := ""
	err = survey.AskOne(&survey.Editor{
		Message: "Please enter the response template for the body (leave empty to not modify the body):",
	}, &bodyTemplate, nil)
	if err != nil {
		return err
	}
	responseTemplate := transformation.Template{}
	responseTemplate.Body = &bodyTemplate

	headers, err := askHeaders("Add a response header?", "Add another response header?")
	if err != nil {
		return err
	}
	if len(headers) != 0 {
		responseTemplate.Header = headers
	}

	spec.ResponseTemplate = &responseTemplate
	// since this is a shared route extension struct with request
	// we should merge individual components
	// FIXME - ashish; deferred. waiting for changes in request transformation
	// changes in Gloo and will add after that
	addAll(s, transformation.EncodeRouteExtension(spec))
	return nil
}

func askHeaders(first, more string) (map[string]string, error) {
	headers := make(map[string]string)
	addHeaders := false
	if err := survey.AskOne(&survey.Confirm{
		Message: first,
	}, &addHeaders, nil); err != nil {
		return nil, err
	}
	for addHeaders {
		questions := []*survey.Question{
			{
				Name:     "key",
				Prompt:   &survey.Input{Message: "Please enter HTTP header name:"},
				Validate: survey.Required,
			},
			{
				Name:     "value",
				Prompt:   &survey.Input{Message: "Please enter HTTP header template:"},
				Validate: survey.Required,
			},
		}
		answers := struct {
			Key   string
			Value string
		}{}
		if err := survey.Ask(questions, &answers); err != nil {
			return nil, err
		}
		headers[answers.Key] = answers.Value

		if err := survey.AskOne(&survey.Confirm{
			Message: more,
		}, &addHeaders, nil); err != nil {
			return nil, err
		}
	}
	return headers, nil
}
