package main

import (
	"fmt"
	"github.com/yugo412/leksikon/cli/command"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:  "leksikon",
	Long: "leksikon is a dictionary service",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("leksikon is a dictionary service")

		return nil
	},
}

func main() {
	var err error

	rootCommand.AddCommand(command.SeedLanguage)
	rootCommand.AddCommand(command.DisplayLanguage)

	err = rootCommand.Execute()
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}
