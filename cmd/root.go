package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/bnb-chain/token-recover-approver/internal/app"
	"github.com/bnb-chain/token-recover-approver/pkg/util"
)

// Root command
var (
	rootCmd = &cobra.Command{
		Run: func(_ *cobra.Command, _ []string) {
			app, err := app.Initialize(cfgFile)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			util.Launch(app.Start, app.Stop, time.Duration(timeout)*time.Second)
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}
}

// Flags
var (
	timeout uint
	cfgFile string
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().UintVar(&timeout, "timeout", 300, "graceful shutdown timeout (second)")
}
