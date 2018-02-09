package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/gluectl/platform"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

const (
	ParameterTypeBool   = "bool"
	ParameterTypeInt    = "int"
	ParameterTypeString = "string"
)

var (
	uparams   = &platform.UpstreamParams{}
	specs     = make(map[string]map[string]interface{})
	defSpecs  = make(map[string]map[string]interface{})
	paramDefs = make(map[string][]ParamDefinition)
)

type ParameterType string

type ParamDefinition struct {
	Name         string
	Description  string
	Type         ParameterType
	DefaultValue interface{}
}

func CreateNameParam(cmds ...*cobra.Command) {
	for _, c := range cmds {
		c.PersistentFlags().StringVar(&uparams.Name, "name", "", "upstream name")
	}
}

func CreateTypeParam(cmds ...*cobra.Command) {
	for _, c := range cmds {
		c.PersistentFlags().StringVar(&uparams.UType, "type", "", "upstream type")
	}
}

func CreateSpecParams(cmds ...*cobra.Command) {

	fmt.Println("Reading Spec definitions for Glue Plugins ...")
	readParamsDefinitions()

	for t, m := range paramDefs {
		specs[t] = make(map[string]interface{})
		defSpecs[t] = make(map[string]interface{})
		for _, s := range m {
			name := fmt.Sprintf("spec.%s", s.Name)
			switch s.Type {
			case ParameterTypeString:
				b := s.DefaultValue.(string)
				specs[t][s.Name] = &b
				defSpecs[t][s.Name] = s.DefaultValue.(string)
				for _, cmd := range cmds {
					cmd.PersistentFlags().StringVar(&b, name, b, s.Description)
				}
			case ParameterTypeInt:
				b := s.DefaultValue.(int)
				specs[t][s.Name] = &b
				defSpecs[t][s.Name] = s.DefaultValue.(int)
				for _, cmd := range cmds {
					cmd.PersistentFlags().IntVar(&b, name, b, s.Description)
				}
			case ParameterTypeBool:
				b := s.DefaultValue.(bool)
				specs[t][s.Name] = &b
				defSpecs[t][s.Name] = s.DefaultValue.(bool)
				for _, cmd := range cmds {
					cmd.PersistentFlags().BoolVar(&b, name, b, s.Description)
				}
			default:
			}
		}
	}
}

func GetUpstreamParams() *platform.UpstreamParams {
	if uparams.UType != "" {
		uparams.Spec = specs[uparams.UType]
	}
	return uparams
}

func LoadUpstreamParamsFromFile() {
	if gparams.FileName == "" {
		return
	}

	var config v1.Upstream
	source, err := ioutil.ReadFile(gparams.FileName)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}

	// Merge params
	if uparams.Name == "" {
		uparams.Name = config.Name
	}
	if uparams.UType == "" {
		uparams.UType = string(config.Type)
	}

	if uparams.UType != "" {
		for n, p := range specs[uparams.UType] {
			v := defSpecs[uparams.UType][n]
			eq := false
			switch v.(type) {
			case string:
				s := p.(*string)
				eq = *s == v
			case int:
				s := p.(*int)
				eq = *s == v
			case bool:
				s := p.(*bool)
				eq = *s == v
			default:
			}
			if eq {
				if config.Spec[n] != nil {
					specs[uparams.UType][n] = config.Spec[n]
				}
			}
		}
	}
}

func readParamsDefinitions() {
	// TODO: Actually read params!!!
	paramDefs["aws"] = []ParamDefinition{
		{
			Name:         "region",
			Description:  "aws region",
			Type:         ParameterTypeString,
			DefaultValue: "us-east-1",
		},
		{
			Name:         "secret",
			Description:  "aws secret reference",
			Type:         ParameterTypeString,
			DefaultValue: "",
		},
	}
	paramDefs["kubernetes"] = []ParamDefinition{
		{
			Name:         "servicename",
			Description:  "k8s service name",
			Type:         ParameterTypeString,
			DefaultValue: "",
		},
		{
			Name:         "serviceport",
			Description:  "k8s service port",
			Type:         ParameterTypeInt,
			DefaultValue: 0,
		},
	}
}
