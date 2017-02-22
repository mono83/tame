package cmd

import (
	"errors"
	"fmt"
	"github.com/mono83/tame/page/recipes/feed"
	"github.com/mono83/tame/user"
	"github.com/spf13/cobra"
)

var feedCmd = &cobra.Command{
	Use:   "feed url",
	Short: "Downloads and parses feed (RSS or Atom)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("URL not provided")
		}
		// Making new user
		u := user.New()

		// Downloading remote page
		page, err := u.Get(args[0])
		if err != nil {
			return err
		}

		// Parsing
		var f feed.Feed
		err = feed.Recipe(page, &f)
		if err != nil {
			return err
		}

		fmt.Println("Feed title:      ", f.Title)
		fmt.Println("Feed link:       ", f.Link)
		fmt.Println("Feed description:", f.Description)
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
