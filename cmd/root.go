package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tlmanz/cgen/scaffold"
)

var rootCmd = &cobra.Command{
	Use:   "cgen",
	Short: "cgen - scaffold new Go microservices from the Catalyst framework",
	Long: `cgen is a CLI tool that scaffolds new Go microservices based on the
Catalyst clean-architecture framework template.`,
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new Catalyst project",
	Long: `Create a new Go microservice project from the Catalyst framework template.

Example:
  cgen new --module=github.com/myorg/myservice
  cgen new --module=github.com/myorg/myservice --dir=./projects`,
	RunE: runNew,
}

var (
	flagModule string
	flagDir    string
	flagYes    bool
)

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringVar(&flagModule, "module", "", "Go module name (e.g. github.com/myorg/myservice)")
	newCmd.Flags().StringVar(&flagDir, "dir", ".", "Parent directory where the project folder will be created")
	newCmd.Flags().BoolVar(&flagYes, "yes", false, "Skip the confirmation prompt")
	_ = newCmd.MarkFlagRequired("module")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runNew(cmd *cobra.Command, _ []string) error {
	projectDir := flagDir
	if !cmd.Flags().Changed("dir") {
		// No --dir given: default to a subdirectory named after the service.
		parts := strings.Split(strings.TrimRight(flagModule, "/"), "/")
		projectDir = filepath.Join(flagDir, parts[len(parts)-1])
	}
	return scaffold.New(flagModule, projectDir, flagYes)
}
