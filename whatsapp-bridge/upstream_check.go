package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type githubCommit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Committer struct {
			Date time.Time `json:"date"`
		} `json:"committer"`
	} `json:"commit"`
}

func readPinnedCommit(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "UPSTREAM_COMMIT=") {
			return strings.TrimSpace(strings.TrimPrefix(line, "UPSTREAM_COMMIT=")), nil
		}
	}
	return "", fmt.Errorf("UPSTREAM_COMMIT not found in %s", path)
}

// CheckUpstream fetches the latest commit from the upstream repo and prints a
// warning if it differs from the SHA pinned in UPSTREAM.md. Soft-fails on
// network errors so it never blocks bridge startup.
func CheckUpstream(upstreamMDPath string, client *http.Client, apiURL string, out io.Writer) {
	pinned, err := readPinnedCommit(upstreamMDPath)
	if err != nil {
		fmt.Fprintf(out, "upstream check: could not read %s: %v\n", upstreamMDPath, err)
		return
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Fprintf(out, "upstream check: bad URL: %v\n", err)
		return
	}
	req.Header.Set("User-Agent", "whatsapp-mcp-fork/1.0")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(out, "upstream check: could not reach GitHub (offline?): %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(out, "upstream check: GitHub API returned %d\n", resp.StatusCode)
		return
	}

	var commit githubCommit
	if err := json.NewDecoder(resp.Body).Decode(&commit); err != nil {
		fmt.Fprintf(out, "upstream check: could not parse response: %v\n", err)
		return
	}

	if commit.SHA == pinned {
		return
	}

	latest := commit.SHA
	if len(latest) > 7 {
		latest = latest[:7]
	}
	pinnedShort := pinned
	if len(pinnedShort) > 7 {
		pinnedShort = pinnedShort[:7]
	}

	fmt.Fprintf(out, "⚠  Upstream whatsapp-mcp has new commits since your last sync.\n")
	fmt.Fprintf(out, "   Latest: %s (%s)\n", latest, commit.Commit.Committer.Date.Format("2006-01-02"))
	fmt.Fprintf(out, "   Pinned: %s\n", pinnedShort)
	fmt.Fprintf(out, "   Review: https://github.com/lharries/whatsapp-mcp/compare/%s...main\n", pinnedShort)
}
