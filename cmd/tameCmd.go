package cmd

import (
	"github.com/spf13/cobra"
)

// TameCmd is root command for tame
var TameCmd = &cobra.Command{
	Use:   "tame",
	Short: "Tame CLI toolset",
}

func init() {
	TameCmd.AddCommand(
		agentsCmd,
		mitmCmd,
		pageCmd,
		feedCmd,
		htmlHeadCmd,
	)
}
