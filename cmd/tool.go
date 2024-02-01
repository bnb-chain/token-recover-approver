package cmd

import (
	"fmt"
	"os"

	"github.com/bnb-chain/token-recover-approver/internal/app/tool"
	"github.com/spf13/cobra"
)

var toolCmd = &cobra.Command{
	Use:   "tool",
	Short: "tool",
	Long: "tool is a command line tool for token-recover-approver, " +
		"it can be used to do some maintenance work, such as migration, etc.",
}

var migrationFromLocalToSQLCmd = &cobra.Command{
	Use:   "migration-from-local-to-sql",
	Short: "migration from local to sql",
	Long:  "migrate data from local to sql store",
	Run: func(cmd *cobra.Command, args []string) {
		if len(migrationFromLocalToSQLConfigPath) == 0 {
			fmt.Println("proof_path is required")
			os.Exit(1)
		}

		tool, err := tool.Initialize(cfgFile)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		err = tool.MigrateDataFromLocalToSQL(migrationFromLocalToSQLConfigPath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

var (
	migrationFromLocalToSQLConfigPath string
)

func init() {
	migrationFromLocalToSQLCmd.Flags().StringVar(&migrationFromLocalToSQLConfigPath, "proof_path", "", "proof file path")
	toolCmd.AddCommand(migrationFromLocalToSQLCmd)
}
