package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a resource",
	Long:  `Create a resource`,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete resources by name",
	Long:  `Delete resources by name`,
}

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "show details of specific resource",
	Long:  `Show details of specific resource`,
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "display one or many resources",
	Long:  `Display one or many resources`,
}
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update",
	Long:  `Update`,
}

func init() {
	rootCmd.AddCommand(createCmd)
	CreateGlobalFlags(createCmd, true)
	rootCmd.AddCommand(deleteCmd)
	CreateGlobalFlags(deleteCmd, true)
	rootCmd.AddCommand(describeCmd)
	CreateGlobalFlags(describeCmd, false)
	rootCmd.AddCommand(getCmd)
	CreateGlobalFlags(getCmd, false)
	rootCmd.AddCommand(updateCmd)
	CreateGlobalFlags(updateCmd, true)
}
