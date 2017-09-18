package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	plugin "github.com/kevinjqiu/helm-plugin-asset/pkg"
)

var (
	valuesFile string
	assetDir string
	iAmSure bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "helm-plugin-asset",
	Short: "A helm plugin to manage chart assets",
	Long: `Render asset files and compile them into the values override file.
Used with the Dynamic ConfigMap chart.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		assets, err := plugin.NewAssets(assetDir, valuesFile, !iAmSure)
		if err != nil {
			return err
		}

		err = assets.Render()
		if err != nil {
			return err
		}
		fmt.Println("Updated! You can now run:")
		fmt.Printf("    helm install -f %s CHART_NAME\n", valuesFile)
		fmt.Println("or:")
		fmt.Printf("    helm upgrade RELEASE_NAME -f %s CHART_NAME\n", valuesFile)
		return nil
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.Flags().StringVarP(&valuesFile, "values", "f", "", "Values override file")
	RootCmd.Flags().StringVarP(&assetDir, "asset-dir", "d", "", "The parent directory of the assets")
	RootCmd.Flags().BoolVar(&iAmSure, "iamsure", false, "Force update of the values override file")
}
