package cmd

import (
	"github.com/mono83/slf/recievers/ansi"
	"github.com/mono83/slf/wd"
	"github.com/spf13/cobra"
)

// TameCmd is root command for tame
var TameCmd = &cobra.Command{
	Use:   "tame",
	Short: "Tame CLI toolset",
}

func init() {
	var verbose bool

	TameCmd.AddCommand(fetchCmd, httpbinCmd, agentsCmd)
	TameCmd.SilenceUsage = true
	TameCmd.SilenceErrors = true
	TameCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	TameCmd.PersistentFlags().BoolP("quiet", "q", false, "Quiet mode")
	TameCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// Enabling logger
		q, _ := cmd.Flags().GetBool("quiet")
		if verbose && !q {
			wd.AddReceiver(ansi.New(true, false, false))
		}
	}
}
