/*
Copyright © 2024 nikkehtine <nikkehtine@int.pl>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/nikkehtine/maison/builder"

	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the site",
	Long: `Build the whole site.
By default outputs to ./public, unless specified otherwise
with the -o flag or in the config.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Maison version " + VERSION)

		if len(args) == 0 {
			log.Fatal("No path specified")
		}

		for _, path := range args {
			builder := new(builder.Builder)
			builder.Init(path)
			builder.Build()
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
