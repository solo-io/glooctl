package cmd

import (
	"fmt"
	"log"

	"github.com/solo-io/gluectl/platform"
	"github.com/spf13/cobra"
)

const (
	ParameterTypeBool   = "bool"
	ParameterTypeInt    = "int"
	ParameterTypeString = "string"
)

type ParameterType string

type ParamDefinition struct {
	Name         string
	Description  string
	Type         ParameterType
	DefaultValue interface{}
}

var (
	gparams         = &platform.GlobalParams{}
	uparams         = &platform.UpstreamParams{}
	specs           = make(map[string]map[string]interface{})
	paramDefs       = make(map[string][]ParamDefinition)
	paramDefsLoaded = false
)

func CreateGlobalFlags(cmd *cobra.Command, isWait bool) {
	cmd.PersistentFlags().StringVarP(&gparams.FileName, "file", "f", "", "file with resource definition")
	cmd.PersistentFlags().StringVarP(&gparams.Namespace, "namespace", "n", "", "resource namespace")

	if isWait {
		cmd.PersistentFlags().IntVarP(&gparams.WaitSec, "wait", "w", 0, "seconds to wait, 0 - return immediately")
	}
}

func GetGlobalFlags() *platform.GlobalParams {
	return gparams
}

func CreateUpstreamParams(cmd *cobra.Command, isSpec bool) {
	cmd.PersistentFlags().StringVar(&uparams.Name, "name", "", "upstream name")
	cmd.PersistentFlags().StringVar(&uparams.UType, "type", "", "upstream type")

	if isSpec {
		if !paramDefsLoaded {
			// TODO: Shouldn't need a lock, but check ...
			log.Println("Reading Spec definitions for Glue Plugins ...")
			readParamsDefinitions()
			paramDefsLoaded = true
		}

		for t, m := range paramDefs {
			specs[t] = make(map[string]interface{})
			for _, s := range m {
				name := fmt.Sprintf("%s.%s", t, s.Name)
				switch s.Type {
				case ParameterTypeString:
					b := s.DefaultValue.(string)
					specs[t][s.Name] = &b
					cmd.PersistentFlags().StringVar(&b, name, b, s.Description)
				case ParameterTypeInt:
					b := s.DefaultValue.(int)
					specs[t][s.Name] = &b
					cmd.PersistentFlags().IntVar(&b, name, b, s.Description)
				case ParameterTypeBool:
					b := s.DefaultValue.(bool)
					specs[t][s.Name] = &b
					cmd.PersistentFlags().BoolVar(&b, name, b, s.Description)
				default:
				}
			}
		}
	}
}

func GetUpstreamParams() *platform.UpstreamParams {
	if uparams.UType == "" {
		log.Fatal("Upstream interface type was not defined")
	}
	uparams.Spec = specs[uparams.UType]
	return uparams
}

func LoadUpstreamParamsFromFile() {
	if gparams.FileName == "" {
		return
	}
	// TODO: Read YAML/JSON definition from file.
	// Before saving the value, check if it was changed from default and if yes - skip it, because it was provided in command line
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
