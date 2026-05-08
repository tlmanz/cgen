package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// Version is the current release version of cgen.
// Update this constant before tagging and pushing a new release.
const Version = "v0.1.0"

const githubRepo = "tlmanz/cgen"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current version and check for updates",
	Run:   runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(_ *cobra.Command, _ []string) {
	fmt.Printf("cgen %s\n\n", Version)

	latest, err := latestGitHubVersion()
	if err != nil {
		fmt.Printf("could not check for updates: %v\n", err)
		return
	}

	fmt.Printf("Latest release: %s\n", latest)

	switch {
	case Version == latest:
		fmt.Println("You are up to date.")
	default:
		fmt.Println("Update available!")
		fmt.Printf("Run: go install github.com/%s@latest\n", githubRepo)
	}
}

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func latestGitHubVersion() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", githubRepo)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("parsing response: %w", err)
	}

	if release.TagName == "" {
		return "", fmt.Errorf("no releases found")
	}

	return release.TagName, nil
}
