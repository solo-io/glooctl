package util

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/protoutil"
)

type Printer func(interface{}, io.Writer) error

func Print(output, template string, m proto.Message, tblPrn Printer, w io.Writer) error {
	switch strings.ToLower(output) {
	case "yaml":
		return PrintYAML(m, w)
	case "json":
		return PrintJSON(m, w)
	case "template":
		return PrintTemplate(m, template, w)
	default:
		return tblPrn(m, w)
	}
}

func PrintList(output, template string, list interface{}, tblPrn Printer, w io.Writer) error {
	switch strings.ToLower(output) {
	case "yaml":
		return PrintYAMLList(list, w)
	case "json":
		return PrintJSONList(list, w)
	case "template":
		return PrintTemplate(list, template, w)
	default:
		return tblPrn(list, w)
	}
}

func PrintJSON(m proto.Message, w io.Writer) error {
	b, err := protoutil.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "unable to convert to JSON")
	}
	_, err = fmt.Fprintln(w, string(b))
	return err
}

func PrintYAML(m proto.Message, w io.Writer) error {
	jsn, err := protoutil.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "uanble to marshal")
	}
	b, err := yaml.JSONToYAML(jsn)
	if err != nil {
		return errors.Wrap(err, "unable to convert to YAML")
	}
	_, err = fmt.Fprintln(w, string(b))
	return err
}

func PrintJSONList(data interface{}, w io.Writer) error {
	list := reflect.ValueOf(data)
	_, err := fmt.Fprintln(w, "[")
	if err != nil {
		return errors.Wrap(err, "unable to print JSON list")
	}
	for i := 0; i < list.Len(); i++ {
		v, ok := list.Index(i).Interface().(proto.Message)
		if !ok {
			return errors.New("unable to convert to proto message")
		}
		if i != 0 {
			if _, err := fmt.Fprintln(w, ","); err != nil {
				return errors.Wrap(err, "unable to print JSON list")
			}
		}
		if err := PrintJSON(v, w); err != nil {
			return err
		}
	}
	_, err = fmt.Fprintln(w, "]")
	return err
}

func PrintYAMLList(data interface{}, w io.Writer) error {
	list := reflect.ValueOf(data)
	for i := 0; i < list.Len(); i++ {
		v, ok := list.Index(i).Interface().(proto.Message)
		if !ok {
			return errors.New("unable to convert to proto message")
		}
		if _, err := fmt.Fprintln(w, "---"); err != nil {
			return errors.Wrap(err, "unable to print YAML list")
		}
		if err := PrintYAML(v, w); err != nil {
			return err
		}
	}
	return nil
}

func PrintTemplate(data interface{}, tmpl string, w io.Writer) error {
	t, err := template.New("output").Parse(tmpl)
	if err != nil {
		return errors.Wrap(err, "unable to parse template")
	}
	return t.Execute(w, data)
}
