package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
	s := getString("Name", uparams.Name, isNameRequired)
	uparams.Name = *s
	ns := getString("Namespace", gparams.Namespace, false)
	gparams.Namespace = *ns
}
