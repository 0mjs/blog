package siteinfo

import (
	"os"
	"os/exec"
	"strings"
	"time"
)

// Build identifies the commit the site was rendered from, for the footer stamp.
type Build struct {
	Commit string // short hash; empty when unknown
	URL    string // commit page on GitHub
	Dirty  bool   // uncommitted changes were present
	Time   time.Time
}

// CurrentBuild is resolved once at startup. BUILD_COMMIT overrides git, for
// environments that build from an archive rather than a checkout.
var CurrentBuild = LoadBuild()

const repoURL = "https://github.com/0mjs/blog"

func LoadBuild() Build {
	build := Build{Time: time.Now()}
	commit := strings.TrimSpace(os.Getenv("BUILD_COMMIT"))
	if commit == "" {
		commit = git("rev-parse", "HEAD")
		build.Dirty = commit != "" && git("status", "--porcelain") != ""
	}
	if commit == "" {
		return build
	}
	build.URL = repoURL + "/commit/" + commit
	build.Commit = commit[:min(7, len(commit))]
	return build
}

func git(args ...string) string {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
