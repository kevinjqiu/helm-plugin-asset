package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	plugin "github.com/kevinjqiu/helm-plugin-asset/pkg"
)

var (
	valuesFiles plugin.ValuesOverrideFiles
	assetDir string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "helm-plugin-asset",
	Short: "A helm plugin to manage chart assets",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		assets, err := plugin.NewAssets(assetDir, valuesFiles)
		if err != nil {
			panic(err)
		}

		rendered, err := assets.Render()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v", rendered)
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.Flags().VarP(&valuesFiles, "values", "f", "Values override files")
	RootCmd.Flags().StringVar(&assetDir, "asset-dir", "d", "The parent directory of the assets")
}
