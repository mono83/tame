package cmd

import (
	"fmt"
	"github.com/mono83/tame/user"
	"github.com/spf13/cobra"
	"sort"
)

var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "Displays list of user agents",
	Run: func(cmd *cobra.Command, args []string) {
		sort.Strings(user.CommonUserAgents)
		for _, ua := range user.CommonUserAgents {
			fmt.Println(ua)
		}
	},
}
