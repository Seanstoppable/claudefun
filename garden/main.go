package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	svgFile := flag.String("svg", "", "export SVG to file")
	svgWidth := flag.Int("svg-width", 800, "SVG width")
	svgHeight := flag.Int("svg-height", 600, "SVG height")
	detailed := flag.Bool("detailed", false, "show per-file breakdown")
	jsonOut := flag.Bool("json", false, "output raw stats as JSON")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "🌱 Code Garden — visualize your codebase as a garden")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Usage: code-garden [flags] [path]")
		fmt.Fprintln(os.Stderr)
		flag.PrintDefaults()
	}

	flag.Parse()

	dir := "."
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	}

	// Resolve to absolute path for display
	absDir, err := filepath.Abs(dir)
	if err != nil {
		absDir = dir
	}

	// Check the directory exists
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		fmt.Println("🌱 Point me at a codebase and I'll grow you a garden!")
		fmt.Println("  Usage: code-garden [path]")
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n  Error: %v\n", err)
		}
		os.Exit(1)
	}

	// Header
	fmt.Println()
	fmt.Println("╭──────────────────────────────────────────────╮")
	fmt.Println("│          🌱  CODE GARDEN  🌱                 │")
	fmt.Println("╰──────────────────────────────────────────────╯")
	fmt.Println()
	fmt.Printf("  Analyzing: %s...\n", absDir)

	// Analyze
	stats, err := AnalyzeDirectory(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error analyzing directory: %v\n", err)
		os.Exit(1)
	}

	if stats.TotalFiles == 0 {
		fmt.Println("  No source files found. Try pointing at a directory with code!")
		os.Exit(0)
	}

	langNames := make([]string, 0, len(stats.Languages))
	for _, l := range stats.Languages {
		langNames = append(langNames, l.Name)
	}
	fmt.Printf("  Found %s files across %d languages\n\n",
		formatNumber(stats.TotalFiles), len(stats.Languages))

	// JSON output mode
	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(stats); err != nil {
			fmt.Fprintf(os.Stderr, "  Error encoding JSON: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Build garden
	garden := BuildGarden(stats)

	// SVG export
	if *svgFile != "" {
		svgRenderer := NewSVGGardenRenderer(*svgWidth, *svgHeight)
		if err := svgRenderer.RenderToFile(garden, *svgFile); err != nil {
			fmt.Fprintf(os.Stderr, "  Error writing SVG: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("  🎨 SVG exported to %s\n\n", *svgFile)
	}

	// Terminal rendering
	termRenderer := NewTermGardenRenderer()

	fmt.Println(termRenderer.Render(garden))
	fmt.Println()
	fmt.Println(termRenderer.RenderStats(garden, stats))

	// Detailed per-file breakdown
	if *detailed {
		fmt.Println()
		fmt.Println("  📂 Per-file breakdown:")
		fmt.Println("  " + repeatStr("─", 60))
		for _, f := range stats.Files {
			rel, err := filepath.Rel(dir, f.Path)
			if err != nil {
				rel = f.Path
			}
			fmt.Printf("  %-40s %4d lines  %2d funcs  %2d tests\n",
				rel, f.Lines, f.Functions, f.Tests)
		}
		fmt.Println("  " + repeatStr("─", 60))
	}

	// Legend
	fmt.Println()
	fmt.Println(termRenderer.RenderLegend())
}

func repeatStr(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

