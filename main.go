package main

import (
	"fmt"
	"os"

	"github.com/robandpdx/gh-unlock-source-repo/pkg/logger"

	"go.uber.org/zap"

	"github.com/robandpdx/gh-unlock-source-repo/cmd"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:     "gh-unlock-source-repo --org <org-name> --repo <repo-name>",
		Short:   "Unlock a repository in GitHub",
		Long:    `Unlock the specified migration source repository in GitHub.`,
		Example: `  gh-unlock-source-repo --org my-org --repo my-repo`,
		RunE: func(cmdCobra *cobra.Command, args []string) error {
			return cmd.UnlockRepo().RunE(cmdCobra, args)
		},
	}

	rootCmd.Flags().StringP("org", "o", "", "GitHub organization name")
	rootCmd.Flags().StringP("repo", "r", "", "GitHub repository name")

	errFile := rootCmd.MarkFlagRequired("org")
	if errFile != nil {
		logger.Logger.Error("failed to mark flag as required", zap.Error(errFile))
	}
	errFile = rootCmd.MarkFlagRequired("repo")
	if errFile != nil {
		logger.Logger.Error("failed to mark flag as required", zap.Error(errFile))
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func init() {
	logger.InitLogger()
	defer logger.SyncLogger()

	required := []struct {
		name  string
		value string
	}{
		{"GH_SOURCE_PAT", os.Getenv("GH_SOURCE_PAT")},
	}

	var missing []string

	for _, r := range required {
		if r.value == "" {
			missing = append(missing, r.name)
		}
	}

	if len(missing) > 0 {
		logger.Logger.Error("Missing required environment variables",
			zap.Strings("missing", missing))
		os.Exit(1)
	}
}
