package route

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"

	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/storage"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/AlecAivazis/survey.v1/core"
)

type SelectionResult struct {
	Selected    []*v1.Route
	NotSelected []*v1.Route
}

func SelectInteractive(routes []*v1.Route, multi bool) (*SelectionResult, error) {
	// group routes by upstream so that selection list is smaller
	routesByUpstream := make(map[string][]*v1.Route)
	for _, r := range routes {
		for _, d := range Destinations(r) {
			routesByUpstream[d.Upstream] = append(routesByUpstream[d.Upstream], r)
		}
	}

	upstreams := make([]string, len(routesByUpstream))
	i := 0
	for k := range routesByUpstream {
		upstreams[i] = k
		i++
	}
	var upstream string
	if err := survey.AskOne(&survey.Select{
		Message: "Please select the upstream for the route:",
		Options: upstreams}, &upstream, survey.Required); err != nil {
		return nil, err
	}

	routesByName := make(map[string]*v1.Route)
	routeOptions := []string{}
	for _, r := range routesByUpstream[upstream] {
		name := toSelectionOption(r)
		routeOptions = append(routeOptions, name)
		routesByName[name] = r
	}
	if multi {
		core.MarkedOptionIcon = "☑"
		core.UnmarkedOptionIcon = "☐"
		var selections []string
		if err := survey.AskOne(&survey.MultiSelect{
			Message: "Please select routes:",
			Options: routeOptions,
		}, &selections, survey.Required); err != nil {
			return nil, err
		}
		selected := make([]*v1.Route, len(selections))
		for i, r := range selections {
			selected[i] = routesByName[r]
		}
		return &SelectionResult{
			Selected:    selected,
			NotSelected: notSelected(routes, selected),
		}, nil
	}

	// single
	var selection string
	if err := survey.AskOne(&survey.Select{
		Message: "Please select a route:",
		Options: routeOptions,
	}, &selection, survey.Required); err != nil {
		return nil, err
	}
	selected := []*v1.Route{routesByName[selection]}
	return &SelectionResult{
		Selected:    selected,
		NotSelected: notSelected(routes, selected),
	}, nil
}

func toSelectionOption(r *v1.Route) string {
	matcher, rType, verb, header := Matcher(r)
	hasExtensions := r.Extensions != nil && len(r.Extensions.Fields) != 0

	return fmt.Sprintf("%s: %s. Verb: %s Headers: %s Has Extensions: %v",
		rType, matcher, verb, header, hasExtensions)
}

// TODO use a map if this is too slow
func notSelected(routes []*v1.Route, selected []*v1.Route) []*v1.Route {
	var filtered []*v1.Route
Route:
	for _, r := range routes {
		for _, s := range selected {
			if r == s {
				continue Route
			}
		}
		filtered = append(filtered, r)
	}
	return filtered
}

func RouteInteractive(sc storage.Interface, r *v1.Route) error {
	upstreams, err := sc.V1().Upstreams().List()
	if err != nil {
		return errors.Wrap(err, "unable to get upstreams")
	}
	if err := matcher(r); err != nil {
		return err
	}
	if err := destination(r, upstreams); err != nil {
		return err
	}
	if err := extensions(r); err != nil {
		return err
	}

	return nil
}

type request struct {
	path    string
	verbs   []string
	headers map[string]string
}

func matcher(r *v1.Route) error {
	oldValues := getMatcherValues(r)
	prompt := &survey.Select{
		Message: "Please select the type of matcher for the route:",
		Options: []string{"event", "path-prefix", "path-regex", "path-exact"},
		Default: oldValues.matcherType,
	}
	var mType string
	if err := survey.AskOne(prompt, &mType, survey.Required); err != nil {
		return err
	}

	switch mType {
	case "event":
		prompt := &survey.Input{
			Message: "Please enter the event type:",
			Default: oldValues.matcher,
		}
		var event string
		survey.AskOne(prompt, &event, survey.Required)
		r.Matcher = &v1.Route_EventMatcher{
			EventMatcher: &v1.EventMatcher{EventType: event},
		}
	case "path-prefix":
		request, err := requestMatcher(oldValues)
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
		request, err := requestMatcher(oldValues)
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
		request, err := requestMatcher(oldValues)
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

type matcherValues struct {
	matcherType string
	matcher     string
	verbs       []string
	headers     map[string]string
}

func getMatcherValues(r *v1.Route) matcherValues {
	switch m := r.GetMatcher().(type) {
	case *v1.Route_EventMatcher:
		return matcherValues{
			matcherType: "event",
			matcher:     m.EventMatcher.EventType,
		}
	case *v1.Route_RequestMatcher:
		switch p := m.RequestMatcher.GetPath().(type) {
		case *v1.RequestMatcher_PathExact:
			return matcherValues{
				matcherType: "path-exact",
				matcher:     p.PathExact,
				verbs:       m.RequestMatcher.Verbs,
				headers:     m.RequestMatcher.Headers,
			}
		case *v1.RequestMatcher_PathPrefix:
			return matcherValues{
				matcherType: "path-prefix",
				matcher:     p.PathPrefix,
				verbs:       m.RequestMatcher.Verbs,
				headers:     m.RequestMatcher.Headers,
			}
		case *v1.RequestMatcher_PathRegex:
			return matcherValues{
				matcherType: "path-regex",
				matcher:     p.PathRegex,
				verbs:       m.RequestMatcher.Verbs,
				headers:     m.RequestMatcher.Headers,
			}
		}
	}
	return matcherValues{
		matcherType: "path-prefix",
		matcher:     "",
	}
}

func requestMatcher(oldValues matcherValues) (*request, error) {
	core.MarkedOptionIcon = "☑"
	core.UnmarkedOptionIcon = "☐"

	httpMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	var questions = []*survey.Question{
		{
			Name: "path",
			Prompt: &survey.Input{
				Message: "Please enter the path for the matcher:",
				Default: oldValues.matcher,
			},
			Validate: survey.Required,
		},
		{
			Name: "verb",
			Prompt: &survey.MultiSelect{
				Message: "Please select all the HTTP requests methods that are applicable:",
				Options: httpMethods,
				Default: oldValues.verbs,
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
	// TODO better editing instead of just replacing
	headers := make(map[string]string)
	hasHeaders := false
	if len(oldValues.headers) != 0 {
		// print headers
		printHeaders(oldValues.headers)
		replaceHeaders := false
		if err := survey.AskOne(&survey.Confirm{Message: "Do you want to replace existing HTTP headers?"}, &replaceHeaders, nil); err != nil {
			return nil, err
		}
		if !replaceHeaders {
			headers = oldValues.headers
			hasHeaders = false
		} else {
			hasHeaders = true
		}
	} else {
		if err := survey.AskOne(&survey.Confirm{Message: "Do you want to set HTTP headers?"}, &hasHeaders, nil); err != nil {
			return nil, err
		}
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

func printHeaders(m map[string]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Header", "Value"})
	for k, v := range m {
		table.Append([]string{k, v})
	}
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
	fmt.Println("\n\n") // keep survey happy
}

func prefixRewrite(r *v1.Route) error {
	if r.PrefixRewrite != "" {
		fmt.Printf("Current prefix rewrite: %s\n\n\n", r.PrefixRewrite)
	}
	var prefix string
	if err := survey.AskOne(&survey.Input{
		Message: "Please enter the rewrite prefix (leave empty if you don't want rewrite):",
	}, &prefix, nil); err != nil {
		return err
	}
	r.PrefixRewrite = prefix
	return nil
}

func destination(r *v1.Route, upstreams []*v1.Upstream) error {
	// check if we want to update destinations
	old := Destinations(r)
	if len(old) != 0 {
		printDestination(old)
		replace := false
		if err := survey.AskOne(&survey.Confirm{Message: "Do you want to replace existing destination?"}, &replace, nil); err != nil {
			return err
		}
		if !replace {
			return nil
		}
	}
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

func printDestination(list []Destination) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Upstream", "Function"})
	for _, d := range list {
		table.Append([]string{d.Upstream, d.Function})
	}
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
	fmt.Println("\n\n") // keep survey happy
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

// TODO better editing of existing extensions instead of just replacing
func extensions(r *v1.Route) error {
	if r.Extensions != nil && len(r.Extensions.Fields) != 0 {
		printExtensions(r)
		replace := false
		if err := survey.AskOne(&survey.Confirm{Message: "Do you want to replace existing extensions?"}, &replace, nil); err != nil {
			return err
		}
		if !replace {
			return nil
		}
	}

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
		printExtensions(r)
	}
	return nil
}

func printExtensions(r *v1.Route) {
	// additional newlines at the end necessary to make it work with survey
	fmt.Printf("\nCurrent extensions:\n%s\n\n\n", Extension(r))
}
