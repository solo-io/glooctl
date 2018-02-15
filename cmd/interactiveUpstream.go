package cmd

import (
	"fmt"
	"strconv"
)

func InteractiveModeUpstream(cmd string) {
	if !interactive {
		return
	}
	switch cmd {
	case "create":
		fallthrough
	case "update":
		getNameAndNamespace(true)
		getUpstreamType()
		getUpstreamSpec()
	case "delete":
		getNameAndNamespace(true)
	case "get":
		fallthrough
	case "describe":
		getNameAndNamespace(false)
	default:
	}
}

func getUpstreamType() {
	tt := make([]string, 0, 10)
	for n := range specs {
		//		fmt.Print(n, ",")
		tt = append(tt, n)
	}

	for {
		fmt.Print("Valid Types: ", tt)
		fmt.Println()
		t := getString("Upstream Type", uparams.UType, true)
		if IsUpstreamTypeValid(t) {
			uparams.UType = *t
			return
		}
	}
}

func getUpstreamSpec() {
	//fmt.Println("DEBUG", uparams.UType, specs[uparams.UType])
	for n, t := range specs[uparams.UType] {
		switch t.(type) {
		case *string:
			s := getString(n, *t.(*string), false)
			specs[uparams.UType][n] = s
		case *int:
			for {
				s := getString(n, strconv.Itoa(*t.(*int)), false)
				i, err := strconv.Atoi(*s)
				if err == nil {
					specs[uparams.UType][n] = i
					break
				}
				fmt.Printf("%s requires integer. Please try again...\n", n)
			}
		case *bool:
			for {
				s := getString(n, strconv.FormatBool(*t.(*bool)), false)
				i, err := strconv.ParseBool(*s)
				if err == nil {
					specs[uparams.UType][n] = i
					break
				}
				fmt.Printf("%s requires boolean (true/false). Please try again...\n", n)
			}
		default:
			fmt.Printf("Unknown parameter type: %t\n", t)
		}
	}
}
