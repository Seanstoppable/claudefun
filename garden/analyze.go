package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// Language represents a programming language detected in the codebase.
type Language struct {
	Name      string
	Extension string
	FileCount int
	LineCount int
}

// FileStats holds analysis metrics for a single source file.
type FileStats struct {
	Path       string
	Language   string
	Lines      int
	Functions  int
	Types      int
	Tests      int
	Comments   int
	TODOs      int
	Complexity int
	Imports    int
	BlankLines int
}

// CodebaseStats aggregates metrics across an entire codebase.
type CodebaseStats struct {
	RootPath        string
	TotalFiles      int
	TotalLines      int
	Languages       []Language
	Files           []FileStats
	TotalFuncs      int
	TotalTypes      int
	TotalTests      int
	TotalComments   int
	TotalTODOs      int
	TotalComplexity int
	HealthScore     float64
}

// extensionToLanguage maps file extensions to language names.
var extensionToLanguage = map[string]string{
	".go":   "Go",
	".py":   "Python",
	".js":   "JavaScript",
	".ts":   "TypeScript",
	".jsx":  "JavaScript",
	".tsx":  "TypeScript",
	".java": "Java",
	".rs":   "Rust",
	".c":    "C",
	".cpp":  "C++",
	".h":    "C",
	".hpp":  "C++",
	".rb":   "Ruby",
	".html": "HTML",
	".css":  "CSS",
}

// skipDirs contains directory names to skip during traversal.
var skipDirs = map[string]bool{
	"node_modules": true,
	"vendor":       true,
	".git":         true,
	"__pycache__":  true,
}

// AnalyzeDirectory walks a directory tree and returns aggregated codebase statistics.
func AnalyzeDirectory(path string) (*CodebaseStats, error) {
	stats := &CodebaseStats{RootPath: path}
	langMap := make(map[string]*Language)

	err := filepath.Walk(path, func(fp string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip files we can't read
		}

		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || skipDirs[name] {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(fp))
		lang, ok := extensionToLanguage[ext]
		if !ok {
			return nil
		}

		fs := analyzeFile(fp, lang)
		if fs == nil {
			return nil
		}

		stats.Files = append(stats.Files, *fs)
		stats.TotalFiles++
		stats.TotalLines += fs.Lines
		stats.TotalFuncs += fs.Functions
		stats.TotalTypes += fs.Types
		stats.TotalTests += fs.Tests
		stats.TotalComments += fs.Comments
		stats.TotalTODOs += fs.TODOs
		stats.TotalComplexity += fs.Complexity

		if l, exists := langMap[lang]; exists {
			l.FileCount++
			l.LineCount += fs.Lines
		} else {
			langMap[lang] = &Language{
				Name:      lang,
				Extension: ext,
				FileCount: 1,
				LineCount: fs.Lines,
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	for _, l := range langMap {
		stats.Languages = append(stats.Languages, *l)
	}

	stats.HealthScore = stats.ComputeHealth()
	return stats, nil
}

func analyzeFile(path, lang string) *FileStats {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	fs := &FileStats{
		Path:     path,
		Language: lang,
	}

	// Detect test files by filename for JS/TS
	base := filepath.Base(path)
	isJSTestFile := (lang == "JavaScript" || lang == "TypeScript") &&
		(strings.Contains(base, ".test.") || strings.Contains(base, ".spec."))

	inBlockComment := false
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		fs.Lines++

		if trimmed == "" {
			fs.BlankLines++
			continue
		}

		// Track block comments
		if inBlockComment {
			fs.Comments++
			if strings.Contains(trimmed, "*/") {
				inBlockComment = false
			}
			countTODOs(trimmed, fs)
			continue
		}

		// Detect block comment start
		if strings.HasPrefix(trimmed, "/*") {
			fs.Comments++
			inBlockComment = !strings.Contains(trimmed, "*/")
			countTODOs(trimmed, fs)
			continue
		}

		// Single-line comments
		if isLineComment(trimmed, lang) {
			fs.Comments++
			countTODOs(trimmed, fs)
			continue
		}

		// TODOs can also appear in non-comment lines
		countTODOs(trimmed, fs)

		// Functions
		countFunctions(trimmed, lang, fs)

		// Types
		countTypes(trimmed, lang, fs)

		// Tests
		countTests(trimmed, lang, isJSTestFile, fs)

		// Complexity
		countComplexity(trimmed, fs)

		// Imports
		countImports(trimmed, lang, fs)
	}

	return fs
}

func isLineComment(trimmed, lang string) bool {
	switch lang {
	case "Go", "JavaScript", "TypeScript", "Java", "Rust", "C", "C++":
		return strings.HasPrefix(trimmed, "//")
	case "Python", "Ruby":
		return strings.HasPrefix(trimmed, "#")
	case "HTML":
		return strings.HasPrefix(trimmed, "<!--")
	case "CSS":
		return strings.HasPrefix(trimmed, "/*")
	default:
		return false
	}
}

func countTODOs(line string, fs *FileStats) {
	upper := strings.ToUpper(line)
	for _, marker := range []string{"TODO", "FIXME", "HACK", "XXX", "BUG"} {
		if strings.Contains(upper, marker) {
			fs.TODOs++
			return // count at most one per line
		}
	}
}

func countFunctions(trimmed, lang string, fs *FileStats) {
	switch lang {
	case "Go":
		if strings.HasPrefix(trimmed, "func ") {
			fs.Functions++
		}
	case "Python":
		if strings.HasPrefix(trimmed, "def ") {
			fs.Functions++
		}
	case "JavaScript", "TypeScript":
		if strings.HasPrefix(trimmed, "function ") || strings.Contains(trimmed, "=> {") {
			fs.Functions++
		}
	case "Rust":
		if strings.HasPrefix(trimmed, "fn ") || strings.HasPrefix(trimmed, "pub fn ") ||
			strings.HasPrefix(trimmed, "pub(crate) fn ") || strings.HasPrefix(trimmed, "async fn ") ||
			strings.HasPrefix(trimmed, "pub async fn ") {
			fs.Functions++
		}
	case "Java":
		// Detect Java methods: access modifiers + return type + name + paren
		// e.g. "public void foo(", "private static int bar(", "protected String baz("
		if (strings.HasPrefix(trimmed, "public ") || strings.HasPrefix(trimmed, "private ") ||
			strings.HasPrefix(trimmed, "protected ") || strings.HasPrefix(trimmed, "static ") ||
			strings.HasPrefix(trimmed, "void ") || strings.HasPrefix(trimmed, "abstract ")) &&
			strings.Contains(trimmed, "(") && !strings.HasPrefix(trimmed, "import ") &&
			!strings.HasPrefix(trimmed, "package ") && !strings.Contains(trimmed, " class ") &&
			!strings.Contains(trimmed, " interface ") && !strings.Contains(trimmed, " enum ") &&
			!strings.Contains(trimmed, " new ") {
			fs.Functions++
		}
	case "Ruby":
		if strings.HasPrefix(trimmed, "def ") {
			fs.Functions++
		}
	}
}

func countTypes(trimmed, lang string, fs *FileStats) {
	switch lang {
	case "Go":
		if strings.HasPrefix(trimmed, "type ") {
			fs.Types++
		}
	case "Python", "Java", "JavaScript", "TypeScript", "Ruby":
		if strings.HasPrefix(trimmed, "class ") {
			fs.Types++
		}
		if (lang == "Java" || lang == "TypeScript") && strings.HasPrefix(trimmed, "interface ") {
			fs.Types++
		}
	case "Rust":
		if strings.HasPrefix(trimmed, "struct ") || strings.HasPrefix(trimmed, "pub struct ") {
			fs.Types++
		}
	case "C", "C++":
		if strings.HasPrefix(trimmed, "struct ") || strings.HasPrefix(trimmed, "typedef struct") {
			fs.Types++
		}
		if lang == "C++" && strings.HasPrefix(trimmed, "class ") {
			fs.Types++
		}
	}
}

func countTests(trimmed, lang string, isJSTestFile bool, fs *FileStats) {
	switch lang {
	case "Go":
		if strings.HasPrefix(trimmed, "func Test") {
			fs.Tests++
		}
	case "Python":
		if strings.HasPrefix(trimmed, "def test_") {
			fs.Tests++
		}
	case "JavaScript", "TypeScript":
		if isJSTestFile &&
			(strings.HasPrefix(trimmed, "it(") || strings.HasPrefix(trimmed, "test(") ||
				strings.HasPrefix(trimmed, "it('") || strings.HasPrefix(trimmed, "test('") ||
				strings.HasPrefix(trimmed, "it(\"") || strings.HasPrefix(trimmed, "test(\"")) {
			fs.Tests++
		}
	case "Java":
		if strings.HasPrefix(trimmed, "@Test") ||
			strings.HasPrefix(trimmed, "@ParameterizedTest") ||
			strings.HasPrefix(trimmed, "@RepeatedTest") {
			fs.Tests++
		}
	}
}

func countComplexity(trimmed string, fs *FileStats) {
	keywords := []string{"if ", "for ", "switch ", "case ", "while ", "else ", "elif ", "catch "}
	for _, kw := range keywords {
		if strings.HasPrefix(trimmed, kw) || strings.Contains(trimmed, " "+kw) {
			fs.Complexity++
		}
	}
}

func countImports(trimmed, lang string, fs *FileStats) {
	switch lang {
	case "Go", "Python", "Java":
		if strings.HasPrefix(trimmed, "import ") || trimmed == "import (" {
			fs.Imports++
		}
	case "JavaScript", "TypeScript":
		if strings.HasPrefix(trimmed, "import ") || strings.Contains(trimmed, "require(") {
			fs.Imports++
		}
	case "C", "C++":
		if strings.HasPrefix(trimmed, "#include") {
			fs.Imports++
		}
	case "Ruby":
		if strings.HasPrefix(trimmed, "require ") || strings.HasPrefix(trimmed, "require_relative ") {
			fs.Imports++
		}
	case "Rust":
		if strings.HasPrefix(trimmed, "use ") {
			fs.Imports++
		}
	}
}

// ComputeHealth calculates a 0-100 health score from codebase metrics.
func (cs *CodebaseStats) ComputeHealth() float64 {
	score := 70.0

	if cs.TotalFuncs > 0 {
		testRatio := float64(cs.TotalTests) / float64(cs.TotalFuncs)
		if testRatio > 0.2 {
			score += 10
		}
	}

	if cs.TotalLines > 0 {
		commentRatio := float64(cs.TotalComments) / float64(cs.TotalLines)
		if commentRatio > 0.1 {
			score += 5
		}
	}

	todoPenalty := float64(cs.TotalTODOs) * 5
	if todoPenalty > 20 {
		todoPenalty = 20
	}
	score -= todoPenalty

	if cs.TotalFuncs > 0 {
		avgComplexity := float64(cs.TotalComplexity) / float64(cs.TotalFuncs)
		if avgComplexity > 10 {
			score -= 10
		}
	}

	if len(cs.Languages) > 1 {
		score += 5
	}

	if cs.TotalTests == 0 {
		score -= 15
	}

	if cs.TotalTODOs < 5 {
		score += 10
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}
