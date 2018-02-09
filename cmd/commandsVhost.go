package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// vhostCmd represents the vhost command
var vhostCmd = &cobra.Command{
	Use:   "vhost",
	Short: "virtual host",
	Long:  `Virtual Host`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vhost called")
	},
}

var vhostDelCmd = &cobra.Command{
	Use:   "vhost",
	Short: "virtual host",
	Long:  `Virtual Host`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vhostDel called")
	},
}

var vhostDescribeCmd = &cobra.Command{
	Use:   "vhost",
	Short: "virtual host",
	Long:  `Virtual Host`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vhostDescribe called")
	},
}

var vhostGetCmd = &cobra.Command{
	Use:   "vhost",
	Short: "virtual host",
	Long:  `Virtual Host`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vhostGet called")
	},
}

var vhostUpdateCmd = &cobra.Command{
	Use:   "vhost",
	Short: "virtual host",
	Long:  `Virtual Host`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vhostUpdate called")
	},
}

func init() {
	createCmd.AddCommand(vhostCmd)
	deleteCmd.AddCommand(vhostDelCmd)
	describeCmd.AddCommand(vhostDescribeCmd)
	getCmd.AddCommand(vhostGetCmd)
	updateCmd.AddCommand(vhostUpdateCmd)
}
