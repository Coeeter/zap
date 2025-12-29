package cmd

import (
	"os"

	"github.com/coeeter/zap/scan"
	"github.com/spf13/cobra"
)

var regexMode bool

var rootCmd = &cobra.Command{
	Use:   "zap [folder-name]",
	Short: "A fast way to search and remove folders",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var targetFolder string
		if len(args) > 0 {
			targetFolder = args[0]
		}

		if targetFolder == "" {
			cmd.Help()
			return nil
		}

		root, err := os.Getwd()
		if err != nil {
			return err
		}

		var results []scan.Result

		if regexMode {
			results, err = scan.FindFoldersGlob(root, targetFolder)
			if err != nil {
				return err
			}
		} else {
			results, err = scan.FindFolders(root, targetFolder)
			if err != nil {
				return err
			}
		}

		for _, result := range results {
			cmd.Println(result.Path)
		}

		if len(results) == 0 {
			cmd.Println("No matching folders found.")
		}

		return nil
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&regexMode, "regex", "r", false, "Enable regex mode for folder name matching")
}

func Execute() error {
	return rootCmd.Execute()
}
