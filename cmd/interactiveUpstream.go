package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
		getType()
		getSpec()
	case "delete":
		getNameAndNamespace(true)
	case "get":
		fallthrough
	case "describe":
		getNameAndNamespace(false)
	default:
	}
}

func getString(prompt, defstr string, isRequired bool) *string {
	for {
		fmt.Printf("%s [%s]: ", prompt, defstr)

		reader := bufio.NewReader(os.Stdin)
		s, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		s = strings.Trim(s, " \n")
		//		fmt.Scanln(&s)
		if s == "" {
			s = defstr
		}
		if !isRequired || s != "" {
			return &s
		}
		fmt.Printf("%s cannot be empty. Please try again...\n", prompt)
	}
}

func getNameAndNamespace(isNameRequired bool) {
	s := getString("Upstream Name", uparams.Name, isNameRequired)
	uparams.Name = *s
	ns := getString("Namespace", gparams.Namespace, false)
	gparams.Namespace = *ns
}

func getType() {
	tt := make([]string, 0, 10)
	for n, _ := range specs {
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

func getSpec() {
	fmt.Println("DEBUG", uparams.UType, specs[uparams.UType])
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
