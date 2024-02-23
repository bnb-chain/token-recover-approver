package cmd

import (
	"fmt"
	"os"

	"github.com/bnb-chain/node/app"
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

var verifyDataFromFullnodeCmd = &cobra.Command{
	Use:   "verify-data-from-fullnode",
	Short: "verify data from fullnode",
	Long:  "verify data from fullnode database",
	Run: func(cmd *cobra.Command, args []string) {
		if len(home) == 0 {
			fmt.Println("home path is required")
			os.Exit(1)
		}

		tool, err := tool.Initialize(cfgFile)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		err = tool.VerifyDataFromFullnode(nodeCtx, home, verifyMerkleRoot)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println("verify data from fullnode successfully!")
	},
}

var (
	migrationFromLocalToSQLConfigPath string

	home             string
	verifyMerkleRoot bool

	nodeCtx = app.ServerContext
)

func init() {
	migrationFromLocalToSQLCmd.Flags().StringVar(&migrationFromLocalToSQLConfigPath, "proof_path", "", "proof file path")
	verifyDataFromFullnodeCmd.Flags().StringVar(&home, "home", app.DefaultNodeHome, "directory for config and data")
	verifyDataFromFullnodeCmd.Flags().BoolVar(&verifyMerkleRoot, "verify_merkle_root", false, "verify merkle root")
	toolCmd.AddCommand(migrationFromLocalToSQLCmd, verifyDataFromFullnodeCmd)
}
