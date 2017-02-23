package cmd

import (
	"errors"
	"fmt"
	"github.com/mono83/tame/page/recipes/html"
	"github.com/mono83/tame/user"
	"github.com/spf13/cobra"
)

var htmlHeadCmd = &cobra.Command{
	Use:   "head url",
	Short: "Parses HTML head metadata from URL",
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
		var h html.Head
		err = html.HeadRecipe(page, &h)
		if err != nil {
			return err
		}

		fmt.Printf("%+v\n", h)
		return nil
	},
}
