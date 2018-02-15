package cmd

import "fmt"

func InteractiveModeVhost(cmd string) {
	if !interactive {
		//return
		fmt.Println("Currently VirtualHost can only be configured in the interactive mode")
	}
	switch cmd {
	case "create":
		fallthrough
	case "update":
		getNameAndNamespace(true)
	case "delete":
		getNameAndNamespace(true)
	case "get":
		fallthrough
	case "describe":
		getNameAndNamespace(false)
	default:
	}
}

func getRoutes() {
	for i := 0; ; i++ {
		cont := getString(fmt.Sprintf("Configuring Route %d. Continue?", i), "yes", true)
		if *cont != "yes" {
			break
		}
		getMatcher()
		getWeightedDestinations()
		_ = getString("Prefix Rewrite", "", false)
	}
}

func getMatcher() {

}

func getWeightedDestinations() {

}

func getSingleDestination() {

}
