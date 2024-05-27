/*
Copyright © 2024 nikkehtine <nikkehtine@int.pl>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/nikkehtine/maison/builder"
	"github.com/nikkehtine/maison/options"

	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the site",
	Long: `Build the whole site.
By default outputs to ./public, unless specified otherwise
in the config. You can also specify the directory/file to build`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("Too many arguments")
		}

		config := options.DefaultConfig
		if len(args) == 1 {
			config.Input = args[0]
		}

		builder := &builder.Builder{
			Config: config,
		}

		var err error
		err = builder.Init(config)
		if err != nil {
			log.Fatal(err)
		}
		err = builder.Build()
		if err != nil {
			log.Fatal(err)
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
