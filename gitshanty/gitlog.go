package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const fieldSep = "---GIT-SHANTY-SEP---"

// CommitType classifies the nature of a git commit.
type CommitType int

const (
	RegularCommit    CommitType = iota
	MergeCommit
	InitialCommit
	Bugfix
	Feature
	Refactor
	Documentation
	HotfixUrgent
	Revert
	ForceEvent
	BigChange
	TinyChange
	DeleteSpree
	TestCommit
	DependencyUpdate
)

// Commit holds parsed data from a single git log entry.
type Commit struct {
	Hash         string
	ShortHash    string
	Author       string
	AuthorEmail  string
	Date         time.Time
	Message      string
	FilesChanged int
	Insertions   int
	Deletions    int
	IsMerge      bool
	Types        []CommitType
}

// GitAnalyzer reads and classifies commits from a git repository.
type GitAnalyzer struct {
	RepoPath string
}

// NewGitAnalyzer creates a GitAnalyzer rooted at the given repo path.
func NewGitAnalyzer(repoPath string) *GitAnalyzer {
	return &GitAnalyzer{RepoPath: repoPath}
}

// git executes a git command in the analyzer's repo directory.
func (g *GitAnalyzer) git(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = g.RepoPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

// GetCommits returns the last n commits from the repository.
func (g *GitAnalyzer) GetCommits(n int) ([]Commit, error) {
	return g.fetchCommits(n, "")
}

// GetCommitsByAuthor returns the last n commits by the given author.
func (g *GitAnalyzer) GetCommitsByAuthor(author string, n int) ([]Commit, error) {
	return g.fetchCommits(n, author)
}

// fetchCommits is the shared implementation for retrieving and parsing commits.
func (g *GitAnalyzer) fetchCommits(n int, author string) ([]Commit, error) {
	format := strings.Join([]string{"%H", "%h", "%an", "%ae", "%aI", "%s", "%P"}, fieldSep)
	args := []string{
		"log",
		fmt.Sprintf("-n%d", n),
		"--format=" + format,
		"--shortstat",
	}
	if author != "" {
		if strings.HasPrefix(author, "-") {
			return nil, fmt.Errorf("invalid author name: must not start with '-'")
		}
		args = append(args, "--author="+author)
	}

	raw, err := g.git(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to read git log: %w", err)
	}
	if raw == "" {
		return nil, nil
	}

	return g.parseLogOutput(raw)
}

// parseLogOutput splits the combined log+shortstat output into Commits.
func (g *GitAnalyzer) parseLogOutput(raw string) ([]Commit, error) {
	lines := strings.Split(raw, "\n")
	var commits []Commit

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// A commit line contains our separator.
		if !strings.Contains(line, fieldSep) {
			continue
		}

		c, err := g.parseCommitLine(line)
		if err != nil {
			return nil, err
		}

		// The next non-empty line (if any) may be the --shortstat summary.
		if i+1 < len(lines) {
			next := strings.TrimSpace(lines[i+1])
			if next != "" && !strings.Contains(next, fieldSep) {
				g.parseShortStat(next, &c)
				i++ // consume the stat line
			}
		}

		c.Types = g.ClassifyCommit(&c)
		commits = append(commits, c)
	}

	return commits, nil
}

// parseCommitLine parses a single formatted log line into a Commit.
func (g *GitAnalyzer) parseCommitLine(line string) (Commit, error) {
	parts := strings.Split(line, fieldSep)
	if len(parts) < 7 {
		return Commit{}, fmt.Errorf("unexpected log format: %q", line)
	}

	date, err := time.Parse(time.RFC3339, parts[4])
	if err != nil {
		return Commit{}, fmt.Errorf("parsing date %q: %w", parts[4], err)
	}

	parents := strings.Fields(parts[6])

	return Commit{
		Hash:        parts[0],
		ShortHash:   parts[1],
		Author:      parts[2],
		AuthorEmail: parts[3],
		Date:        date,
		Message:     parts[5],
		IsMerge:     len(parents) >= 2,
	}, nil
}

// parseShortStat extracts file/insertion/deletion counts from a --shortstat line.
// Example: " 3 files changed, 42 insertions(+), 7 deletions(-)"
func (g *GitAnalyzer) parseShortStat(line string, c *Commit) {
	for _, token := range strings.Split(line, ",") {
		token = strings.TrimSpace(token)
		fields := strings.Fields(token)
		if len(fields) < 2 {
			continue
		}
		num, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		switch {
		case strings.Contains(token, "file"):
			c.FilesChanged = num
		case strings.Contains(token, "insertion"):
			c.Insertions = num
		case strings.Contains(token, "deletion"):
			c.Deletions = num
		}
	}
}

// ClassifyCommit determines all applicable CommitTypes for a commit.
func (g *GitAnalyzer) ClassifyCommit(c *Commit) []CommitType {
	var types []CommitType
	msg := strings.ToLower(c.Message)

	if c.IsMerge {
		types = append(types, MergeCommit)
	}

	// Initial commit has no parents — detect via message heuristic or
	// single-parent-less hash (IsMerge is false and message suggests it).
	if strings.Contains(msg, "initial commit") || strings.Contains(msg, "first commit") {
		types = append(types, InitialCommit)
	}

	if containsAny(msg, "fix", "bug", "patch", "hotfix") {
		types = append(types, Bugfix)
	}

	if containsAny(msg, "feat", "add", "new", "implement") {
		types = append(types, Feature)
	}

	if containsAny(msg, "refactor", "rename", "move", "clean", "restructure") {
		types = append(types, Refactor)
	}

	if containsAny(msg, "doc", "readme", "comment") {
		types = append(types, Documentation)
	}

	if containsAny(msg, "urgent", "critical", "emergency", "asap") {
		types = append(types, HotfixUrgent)
	}

	if strings.Contains(msg, "revert") {
		types = append(types, Revert)
	}

	if containsAny(msg, "force push", "force-push", "forcepush", "--force") {
		types = append(types, ForceEvent)
	}

	if c.FilesChanged >= 10 {
		types = append(types, BigChange)
	}

	if c.FilesChanged == 1 && (c.Insertions+c.Deletions) <= 5 {
		types = append(types, TinyChange)
	}

	if c.Deletions > c.Insertions && c.Deletions > 0 {
		types = append(types, DeleteSpree)
	}

	if strings.Contains(msg, "test") {
		types = append(types, TestCommit)
	}

	if containsAny(msg, "dependency", "upgrade", "bump") ||
		(strings.Contains(msg, "update") && strings.Contains(msg, "version")) {
		types = append(types, DependencyUpdate)
	}

	if len(types) == 0 {
		types = append(types, RegularCommit)
	}

	return types
}

// RepoName extracts a human-friendly repository name from the remote origin
// URL, falling back to the directory name.
func (g *GitAnalyzer) RepoName() string {
	remote, err := g.git("remote", "get-url", "origin")
	if err == nil && remote != "" {
		// Handle SSH: git@host:org/repo.git  or HTTPS: https://host/org/repo.git
		name := remote
		if idx := strings.LastIndex(name, "/"); idx >= 0 {
			name = name[idx+1:]
		} else if idx := strings.LastIndex(name, ":"); idx >= 0 {
			name = name[idx+1:]
		}
		name = strings.TrimSuffix(name, ".git")
		if name != "" {
			return name
		}
	}

	return filepath.Base(g.RepoPath)
}

// containsAny returns true if s contains any of the given substrings.
func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
