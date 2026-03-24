package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	n := flag.Int("n", 1, "number of commits to process")
	author := flag.String("author", "", "filter by author name")
	all := flag.Bool("all", false, "generate an epic saga of ALL commits (capped at 20)")
	repo := flag.String("repo", ".", "path to git repo")
	width := flag.Int("width", 60, "output width")
	flag.Parse()

	analyzer := NewGitAnalyzer(*repo)

	var commits []Commit
	var err error

	switch {
	case *all:
		if *author != "" {
			commits, err = analyzer.GetCommitsByAuthor(*author, 20)
		} else {
			commits, err = analyzer.GetCommits(20)
		}
	case *author != "":
		commits, err = analyzer.GetCommitsByAuthor(*author, *n)
	default:
		commits, err = analyzer.GetCommits(*n)
	}

	if err != nil {
		fmt.Println("🏴\u200d☠️ Arr! No git repository found here, matey! Navigate to a repo and try again.")
		os.Exit(1)
	}

	if len(commits) == 0 {
		fmt.Println("🏴\u200d☠️ Arr! No commits found in these waters, matey!")
		os.Exit(0)
	}

	seed := hashToSeed(commits[0].Hash)
	lyrics := NewLyricsGenerator(seed)
	renderer := NewShantyRenderer(*width)

	if *all {
		shanty := lyrics.GenerateEpic(commits)
		fmt.Print(renderer.Render(shanty))
	} else {
		for i, commit := range commits {
			if i > 0 {
				fmt.Println(renderer.RenderDivider())
				fmt.Println()
			}
			shanty := lyrics.GenerateShanty(commit)
			fmt.Print(renderer.Render(shanty))
		}
	}
}

func hashToSeed(hash string) int64 {
	var seed int64
	for _, c := range hash {
		seed = seed*31 + int64(c)
	}
	return seed
}
