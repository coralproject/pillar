package cmd

import (
	"fmt"

	"github.com/coralproject/pillar/pkg/stats"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "stats",
	Short: "Stats is a tool to calculate the statistics of different Coral models",
	Long:  `Stats works as a command line tool to calculate statistics`,
}

func init() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(calculateCmd)
}

//* VERSION COMMAND *//

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Stats",
	Long:  `This is Stats's version.`,
	Run:   version,
}

func version(cmd *cobra.Command, args []string) {
	fmt.Printf("stats version: %v\n", stats.VersionNumber)
}

var calculateCmd = &cobra.Command{
	Use:   "calculate",
	Short: "Calculate stats",
	Long:  `Calculate stats`,
	Run:   calculate,
}

func calculate(cmd *cobra.Command, args []string) {

	stats.Init()
	stats.Calculate()

}
