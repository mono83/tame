package cmd

import (
	"fmt"
	"sort"

	"github.com/mono83/tame/client"
	"github.com/spf13/cobra"
)

var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "Displays list of user agents",
	Run: func(cmd *cobra.Command, args []string) {
		sort.Strings(client.CommonUserAgents)
		for _, ua := range client.CommonUserAgents {
			fmt.Println(ua)
		}
	},
}
