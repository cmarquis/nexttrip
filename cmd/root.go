package cmd

import (
	"fmt"
	"os"

	"github.com/cmarquis/nexttrip/providers"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "[ROUTE] [STOP] [DIRECTION]",
	Short: "Gets the next time transit will be at the specified stop",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		provider := providers.DefaultProviders{Sandboxed: false}
		p := provider.GetProvider("metrotransit") //todo could eventually get this from a config
		nt, err := p.GetNextTrip(args[0], args[1], args[2])
		if err != nil {
			panic(err)
		}
		fmt.Printf("%d Minutes\n", nt)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
