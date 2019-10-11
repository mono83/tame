package cmd

import (
	"errors"
	"fmt"

	"github.com/mono83/tame"
	"github.com/mono83/tame/client"
	"github.com/mono83/tame/recipes/feed"
	"github.com/spf13/cobra"
)

var feedCmd = &cobra.Command{
	Use:     "feed url",
	Aliases: []string{"rss"},
	Short:   "Downloads and parses feed (RSS or Atom)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("URL not provided")
		}
		// Making new client
		cl := client.New()

		// Downloading remote document
		doc, err := cl.Get(args[0])
		if err != nil {
			return err
		}

		// Parsing
		var f feed.Feed
		if err := tame.Unmarshal(doc, &f); err != nil {
			return err
		}

		colorKey.Print("Feed title      : ")
		colorValue.Println(f.Title)
		colorKey.Print("Feed link       : ")
		colorValue.Println(f.Link)
		colorKey.Print("Feed description: ")
		colorValue.Println(f.Description)

		for _, i := range f.Items {
			fmt.Println(" * ", i.Title)
			fmt.Println("   ", i.Link)
			if len(i.Tags) > 0 {
				fmt.Print("   ")
				for _, tag := range i.Tags {
					fmt.Print(" ")
					fmt.Print("[" + tag + "]")
				}
				fmt.Print("\n")
			}
		}

		return nil
	},
}
