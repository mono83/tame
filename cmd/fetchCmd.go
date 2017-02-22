package cmd

import (
	"errors"
	"fmt"
	"github.com/mono83/tame/user"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch url",
	Short: "Fetches remote page",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("URL not provided")
		}

		// Making new user
		u := user.New()
		page, err := u.Get(args[0])
		if err != nil {
			return err
		}

		fmt.Println(page.AsString())
		fmt.Println(page.String())

		return nil
	},
}
