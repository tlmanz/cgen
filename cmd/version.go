package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

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
	current := currentVersion()
	fmt.Printf("cgen %s%s\n\n", current, gitSuffix())

	latest, err := latestGitHubVersion()
	if err != nil {
		fmt.Printf("could not check for updates: %v\n", err)
		return
	}

	fmt.Printf("Latest release: %s\n", latest)

	switch {
	case !isRelease(current):
		fmt.Println("Running a development build.")
		fmt.Printf("Install the release with: go install github.com/%s@latest\n", githubRepo)
	case current == latest:
		fmt.Println("You are up to date.")
	default:
		fmt.Println("Update available!")
		fmt.Printf("Run: go install github.com/%s@latest\n", githubRepo)
	}
}

// currentVersion reads the version Go embedded at install time.
// Returns "(devel)" for local builds without a release tag.
func currentVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}
	return "(devel)"
}

// gitSuffix returns " (commithash[, dirty])" from VCS build settings, or "".
func gitSuffix() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	var rev, modified string
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			if len(s.Value) >= 7 {
				rev = s.Value[:7]
			}
		case "vcs.modified":
			if s.Value == "true" {
				modified = ", dirty"
			}
		}
	}
	if rev == "" {
		return ""
	}
	return fmt.Sprintf(" (%s%s)", rev, modified)
}

// isRelease reports whether version is a proper release tag (e.g. v1.2.3).
// Pseudo-versions (v0.0.0-timestamp-commit) and "(devel)" are not releases.
func isRelease(version string) bool {
	return strings.HasPrefix(version, "v") &&
		!strings.Contains(version, "-") &&
		!strings.Contains(version, "+")
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
