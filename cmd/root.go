/*
Copyright © 2024 nikkehtine <nikkehtine@int.pl>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const VERSION = "0.0.1"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "maison",
	Short: "Static site generator for Markdown",
	Long: `Maison is a static site generator for Markdown,
as well as a simple web server for hosting your generated content.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if rootCmd.Flags().Lookup("version").Value.String() == "true" {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.maison.yaml)")
	rootCmd.Flags().BoolP("version", "v", false, "print version")
}
