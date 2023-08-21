package command

import (
	"log"

	"github.com/spf13/cobra"
)

var SeedLanguage = &cobra.Command{
	Use:  "language",
	Long: "seed language data",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println(args)

		return nil
	},
}
