package cmd

import (
	"errors"
	"fmt"

	"github.com/mono83/tame"
	"github.com/mono83/tame/client"
	"github.com/mono83/tame/goquery"
	"github.com/mono83/tame/recipes/dom"
	"github.com/spf13/cobra"
)

var htmlHeadCmd = &cobra.Command{
	Use:   "head url",
	Short: "Parses HTML head metadata from URL",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("URL not provided")
		}
		// Making new client
		cl := client.New()

		// Downloading remote page
		doc, err := goquery.FromDocumentE(cl.Get(args[0]))
		if err != nil {
			return err
		}

		// Parsing
		var h dom.Head
		var og dom.OpenGraph
		if err := tame.Unmarshal(doc, &h, &og); err != nil {
			return err
		}

		fmt.Printf("%+v\n\n%+v\n\n", h, og)
		return nil
	},
}
