package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeTempUpstreamMD(t *testing.T, sha string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "UPSTREAM.md")
	os.WriteFile(path, []byte("UPSTREAM_COMMIT="+sha+"\n"), 0644)
	return path
}

func makeGitHubServer(sha string, date time.Time) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var c githubCommit
		c.SHA = sha
		c.Commit.Committer.Date = date
		json.NewEncoder(w).Encode(c)
	}))
}

func TestReadPinnedCommit_Found(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "UPSTREAM.md")
	os.WriteFile(path, []byte("# header\n\nUPSTREAM_COMMIT=abc123\nUPSTREAM_DATE=2026-04-21\n"), 0644)

	sha, err := readPinnedCommit(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sha != "abc123" {
		t.Errorf("got %q, want %q", sha, "abc123")
	}
}

func TestReadPinnedCommit_Missing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "UPSTREAM.md")
	os.WriteFile(path, []byte("# no commit here\n"), 0644)

	_, err := readPinnedCommit(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCheckUpstream_UpToDate(t *testing.T) {
	srv := makeGitHubServer("abc123", time.Now())
	defer srv.Close()

	path := writeTempUpstreamMD(t, "abc123")
	var buf bytes.Buffer
	CheckUpstream(path, srv.Client(), srv.URL, &buf)

	if buf.Len() != 0 {
		t.Errorf("expected no output when up to date, got: %q", buf.String())
	}
}

func TestCheckUpstream_NewCommits(t *testing.T) {
	srv := makeGitHubServer("newsha789xyz", time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC))
	defer srv.Close()

	path := writeTempUpstreamMD(t, "oldsha123abc")
	var buf bytes.Buffer
	CheckUpstream(path, srv.Client(), srv.URL, &buf)

	out := buf.String()
	if !strings.Contains(out, "⚠") {
		t.Errorf("expected warning symbol in output, got: %q", out)
	}
	if !strings.Contains(out, "newsha7") {
		t.Errorf("expected latest SHA prefix in output, got: %q", out)
	}
	if !strings.Contains(out, "oldsha1") {
		t.Errorf("expected pinned SHA prefix in output, got: %q", out)
	}
	if !strings.Contains(out, "2026-04-21") {
		t.Errorf("expected date in output, got: %q", out)
	}
}

func TestCheckUpstream_Offline(t *testing.T) {
	path := writeTempUpstreamMD(t, "abc123")
	var buf bytes.Buffer
	// Port 1 is reserved and will refuse connections
	CheckUpstream(path, &http.Client{Timeout: 500 * time.Millisecond}, "http://127.0.0.1:1", &buf)

	out := buf.String()
	if !strings.Contains(out, "upstream check:") {
		t.Errorf("expected soft warning on network failure, got: %q", out)
	}
}
