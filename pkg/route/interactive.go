package route

import (
	"fmt"
	"strconv"

	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/storage"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/AlecAivazis/survey.v1/core"
)

func RouteInteractive(sc storage.Interface) (*v1.Route, error) {
	upstreams, err := sc.V1().Upstreams().List()

	if err != nil {
		return nil, errors.Wrap(err, "unable to get upstreams")
	}
	r := &v1.Route{}
	if err := matcher(r); err != nil {
		return nil, err
	}
	if err := destination(r, upstreams); err != nil {
		return nil, err
	}
	if err := extensions(r); err != nil {
		return nil, err
	}

	return r, nil
}

type request struct {
	path    string
	verbs   []string
	headers map[string]string
}

func matcher(r *v1.Route) error {
	prompt := &survey.Select{
		Message: "Please select the type of matcher for the route:",
		Options: []string{"event", "path-prefix", "path-regex", "path-exact"},
		Default: "path-prefix",
	}
	var mType string
	if err := survey.AskOne(prompt, &mType, survey.Required); err != nil {
		return err
	}

	switch mType {
	case "event":
		prompt := &survey.Input{
			Message: "Please enter the event type:",
		}
		var event string
		survey.AskOne(prompt, &event, survey.Required)
		r.Matcher = &v1.Route_EventMatcher{
			EventMatcher: &v1.EventMatcher{EventType: event},
		}
	case "path-prefix":
		request, err := requestMatcher()
		if err != nil {
			return err
		}
		r.Matcher = &v1.Route_RequestMatcher{
			RequestMatcher: &v1.RequestMatcher{
				Path:    &v1.RequestMatcher_PathPrefix{PathPrefix: request.path},
				Verbs:   request.verbs,
				Headers: request.headers,
			},
		}
	case "path-regex":
		request, err := requestMatcher()
		if err != nil {
			return err
		}
		r.Matcher = &v1.Route_RequestMatcher{
			RequestMatcher: &v1.RequestMatcher{
				Path:    &v1.RequestMatcher_PathRegex{PathRegex: request.path},
				Verbs:   request.verbs,
				Headers: request.headers,
			},
		}
	case "path-exact":
		request, err := requestMatcher()
		if err != nil {
			return err
		}
		r.Matcher = &v1.Route_RequestMatcher{
			RequestMatcher: &v1.RequestMatcher{
				Path:    &v1.RequestMatcher_PathExact{PathExact: request.path},
				Verbs:   request.verbs,
				Headers: request.headers,
			},
		}
	}

	if mType != "event" {
		if err := prefixRewrite(r); err != nil {
			return err
		}
	}
	return nil
}

func requestMatcher() (*request, error) {
	core.MarkedOptionIcon = "☑"
	core.UnmarkedOptionIcon = "☐"

	httpMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	var questions = []*survey.Question{
		{
			Name:     "path",
			Prompt:   &survey.Input{Message: "Please enter the path for the matcher:"},
			Validate: survey.Required,
		},
		{
			Name: "verb",
			Prompt: &survey.MultiSelect{
				Message: "Please select all the HTTP requests methods that are applicable:",
				Options: httpMethods,
			},
		},
	}
	answers := struct {
		Path string
		Verb []string
	}{}
	if err := survey.Ask(questions, &answers); err != nil {
		return nil, err
	}
	// handle case for all
	if len(answers.Verb) == 0 || len(answers.Verb) == len(httpMethods) {
		answers.Verb = nil
	}

	// headers
	headers := make(map[string]string)
	hasHeaders := false
	if err := survey.AskOne(&survey.Confirm{Message: "Do you want to set HTTP headers?"}, &hasHeaders, nil); err != nil {
		return nil, err
	}
	for hasHeaders {
		questions = []*survey.Question{
			{
				Name:     "key",
				Prompt:   &survey.Input{Message: "Please enter HTTP header name:"},
				Validate: survey.Required,
			},
			{
				Name:   "value",
				Prompt: &survey.Input{Message: "Please enter HTTP header value:"},
			},
		}
		header := struct {
			Key   string
			Value string
		}{}
		if err := survey.Ask(questions, &header); err != nil {
			return nil, err
		}
		headers[header.Key] = header.Value
		if err := survey.AskOne(&survey.Confirm{Message: "Do you want to more HTTP headers?"}, &hasHeaders, nil); err != nil {
			return nil, err
		}
	}
	return &request{path: answers.Path, verbs: answers.Verb, headers: headers}, nil
}

func prefixRewrite(r *v1.Route) error {
	var prefix string
	if err := survey.AskOne(&survey.Input{Message: "Please enter the rewrite prefix (leave empty if you don't want rewrite):"}, &prefix, nil); err != nil {
		return err
	}
	r.PrefixRewrite = prefix
	return nil
}

func destination(r *v1.Route, upstreams []*v1.Upstream) error {
	upstreamNames := make([]string, len(upstreams))
	for i, u := range upstreams {
		upstreamNames[i] = u.Name
	}
	var destinations []*Destination
	ask := true
	for ask {
		var name string
		survey.AskOne(&survey.Select{
			Message: "Please select an upstream:",
			Options: upstreamNames,
		}, &name, survey.Required)

		var selectedUpstream *v1.Upstream
		for _, u := range upstreams {
			if u.Name == name {
				selectedUpstream = u
				break
			}
		}
		if len(selectedUpstream.Functions) == 0 {
			destinations = append(destinations, &Destination{Upstream: name})
		} else {
			functionNames := make([]string, len(selectedUpstream.Functions))
			for i, f := range selectedUpstream.Functions {
				functionNames[i] = f.Name
			}
			var fname string
			survey.AskOne(&survey.Select{
				Message: "Please select a function:",
				Options: functionNames,
			}, &fname, survey.Required)
			destinations = append(destinations, &Destination{Upstream: name, Function: fname})
		}
		if err := survey.AskOne(&survey.Confirm{Message: "Does the route have more destinations?"}, &ask, nil); err != nil {
			return err
		}
	}

	switch len(destinations) {
	case 0:
		return errors.New("expected at least one destination for the route")
	case 1:
		r.SingleDestination = toAPIDestination(destinations[0])
	default:
		wd := make([]*v1.WeightedDestination, len(destinations))
		for i, d := range destinations {
			var weight int
			q := fmt.Sprintf("Please enter a weight for destination '%s' >", d)
			survey.AskOne(&survey.Input{Message: q}, &weight, func(val interface{}) error {
				i, err := strconv.Atoi(val.(string))
				if err != nil {
					return errors.New("weight must be an integer greather than 0")
				}
				if i <= 0 {
					return errors.New("weight must be an integer greater than 0")
				}
				return nil
			})
			wd[i] = &v1.WeightedDestination{
				Destination: toAPIDestination(d),
				Weight:      uint32(weight),
			}
		}
		r.MultipleDestinations = wd
	}
	return nil
}

func toAPIDestination(d *Destination) *v1.Destination {
	if d.Function != "" {
		return &v1.Destination{
			DestinationType: &v1.Destination_Function{
				Function: &v1.FunctionDestination{
					UpstreamName: d.Upstream,
					FunctionName: d.Function},
			},
		}
	}

	return &v1.Destination{
		DestinationType: &v1.Destination_Upstream{
			Upstream: &v1.UpstreamDestination{Name: d.Upstream},
		},
	}
}

func extensions(r *v1.Route) error {
	extensionOptions := make([]string, len(availableExtensions)+1)
	extensionOptions[0] = "None"
	for i, e := range availableExtensions {
		extensionOptions[i+1] = e.name
	}
	question := &survey.Select{
		Message: "Select the route extension you want to set:",
		Options: extensionOptions,
	}
	r.Extensions = &types.Struct{}
	var choice string
	for choice != "None" {
		if err := survey.AskOne(question, &choice, survey.Required); err != nil {
			return err
		}
		for _, e := range availableExtensions {
			if e.name == choice {
				if err := e.handler(r.Extensions); err != nil {
					return err
				}
				break
			}
		}
		// additional newlines at the end necessary to make it work with survey
		fmt.Printf("\nCurrent extensions:\n%s\n\n\n", Extension(r))
	}
	return nil
}
