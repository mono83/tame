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

		colorHeader.Println("HTML Header")
		colorKey.Print("Title        : ")
		colorValue.Println(h.Title)
		colorKey.Print("Engine       : ")
		colorValue.Println(h.Engine)
		colorKey.Print("Canonical URL: ")
		colorValue.Println(h.URLCanonical)
		colorKey.Print("Keywords     : ")
		colorValue.Println(h.KeywordsCS())
		colorKey.Print("Description  : ")
		colorValue.Println(h.Description)
		fmt.Println()

		colorHeader.Println("Open Graph Data")
		colorKey.Print("Site       : ")
		colorValue.Println(og.SiteName)
		colorKey.Print("Title      : ")
		colorValue.Println(og.Title)
		colorKey.Print("Type       : ")
		colorValue.Println(og.Type)
		colorKey.Print("Locale     : ")
		colorValue.Println(og.Locale)
		colorKey.Print("URL        : ")
		colorValue.Println(og.URL)
		colorKey.Print("Description: ")
		colorValue.Println(og.Description)
		if len(og.Images) > 0 {
			colorKey.Println("Images")
			for _, i := range og.Images {
				colorKey.Print(" * ")
				colorValue.Println(i.String())
			}
		}
		fmt.Println()

		return nil
	},
}
