package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/coeeter/zap/internal/scan"
	"github.com/coeeter/zap/internal/tui"
	"github.com/spf13/cobra"
)

var searchMode bool

func Execute() error {
	rootCmd := &cobra.Command{
		Use:   "zap [folder-name]",
		Short: "A fast way to search and remove folders",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var targetFolder string
			if len(args) > 0 {
				targetFolder = args[0]
			}

			if targetFolder == "" {
				inputResult, err := tui.RunInput("node_modules")
				if err != nil {
					return err
				}
				if !inputResult.Submitted {
					return nil
				}
				targetFolder = inputResult.Value
			}

			root, err := os.Getwd()
			if err != nil {
				return err
			}

			var results []scan.Result
			if searchMode {
				results, err = scan.FindFoldersGlob(root, targetFolder)
			} else {
				results, err = scan.FindFolders(root, targetFolder)
			}
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Println("No matching folders found.")
				return nil
			}

			tuiResult, err := tui.RunSelector(results)
			if err != nil {
				return err
			}

			if tuiResult.DeleteConfirmed && len(tuiResult.ToDelete) > 0 {
				_, err := tui.RunDelete(tuiResult.ToDelete)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	rootCmd.Flags().BoolVarP(&searchMode, "search", "s", false, "Enable search mode with glob patterns")

	return rootCmd.ExecuteContext(context.Background())
}
