package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Color settings
var (
	colorHeader = color.New(color.FgGreen, color.Underline)
	colorKey    = color.New(color.FgCyan)
	colorValue  = color.New(color.FgWhite)
)

// TameCmd is root command for tame
var TameCmd = &cobra.Command{
	Use:   "tame",
	Short: "Tame CLI toolset",
}

func init() {
	TameCmd.AddCommand(
		agentsCmd,
		pageCmd,
		feedCmd,
		htmlHeadCmd,
	)
}
