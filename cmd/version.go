package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bnb-chain/airdrop-service/internal/version"
)

var (
	versionCmd = &cobra.Command{
		Use:          "version",
		Short:        "Print version information",
		SilenceUsage: true,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(version.Version())
		},
	}
)
