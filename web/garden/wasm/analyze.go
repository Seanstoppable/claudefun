package main

import (
	"bufio"
	"math"
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

// languageFromName maps user-facing language names to canonical names.
var languageFromName = map[string]string{
	"go":         "Go",
	"python":     "Python",
	"javascript": "JavaScript",
	"js":         "JavaScript",
	"typescript": "TypeScript",
	"ts":         "TypeScript",
	"java":       "Java",
	"rust":       "Rust",
	"c":          "C",
	"cpp":        "C++",
	"c++":        "C++",
	"ruby":       "Ruby",
	"rb":         "Ruby",
	"html":       "HTML",
	"css":        "CSS",
}

// analyzeString analyzes a single code string and returns codebase statistics.
func analyzeString(code, language string) *CodebaseStats {
	lang := languageFromName[strings.ToLower(language)]
	if lang == "" {
		lang = language
	}

	fs := &FileStats{
		Path:     "pasted-code",
		Language: lang,
	}

	inBlockComment := false
	scanner := bufio.NewScanner(strings.NewReader(code))

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		fs.Lines++

		if trimmed == "" {
			fs.BlankLines++
			continue
		}

		if inBlockComment {
			fs.Comments++
			if strings.Contains(trimmed, "*/") {
				inBlockComment = false
			}
			countTODOs(trimmed, fs)
			continue
		}

		if strings.HasPrefix(trimmed, "/*") {
			fs.Comments++
			inBlockComment = !strings.Contains(trimmed, "*/")
			countTODOs(trimmed, fs)
			continue
		}

		if isLineComment(trimmed, lang) {
			fs.Comments++
			countTODOs(trimmed, fs)
			continue
		}

		countTODOs(trimmed, fs)
		countFunctions(trimmed, lang, fs)
		countTypes(trimmed, lang, fs)
		countTests(trimmed, lang, false, fs)
		countComplexity(trimmed, fs)
		countImports(trimmed, lang, fs)
	}

	stats := &CodebaseStats{
		RootPath:        "Pasted Code",
		TotalFiles:      1,
		TotalLines:      fs.Lines,
		TotalFuncs:      fs.Functions,
		TotalTypes:      fs.Types,
		TotalTests:      fs.Tests,
		TotalComments:   fs.Comments,
		TotalTODOs:      fs.TODOs,
		TotalComplexity: fs.Complexity,
		Languages:       []Language{{Name: lang, FileCount: 1, LineCount: fs.Lines}},
		Files:           []FileStats{*fs},
	}
	stats.HealthScore = stats.ComputeHealth()
	return stats
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
		return strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#")
	}
}

func countTODOs(line string, fs *FileStats) {
	upper := strings.ToUpper(line)
	for _, marker := range []string{"TODO", "FIXME", "HACK", "XXX", "BUG"} {
		if strings.Contains(upper, marker) {
			fs.TODOs++
			return
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
	case "Ruby":
		if strings.HasPrefix(trimmed, "def ") {
			fs.Functions++
		}
	case "Java":
		if (strings.HasPrefix(trimmed, "public ") || strings.HasPrefix(trimmed, "private ") ||
			strings.HasPrefix(trimmed, "protected ") || strings.HasPrefix(trimmed, "static ") ||
			strings.HasPrefix(trimmed, "void ") || strings.HasPrefix(trimmed, "abstract ")) &&
			strings.Contains(trimmed, "(") && !strings.HasPrefix(trimmed, "import ") &&
			!strings.HasPrefix(trimmed, "package ") && !strings.Contains(trimmed, " class ") &&
			!strings.Contains(trimmed, " interface ") && !strings.Contains(trimmed, " enum ") &&
			!strings.Contains(trimmed, " new ") {
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
		if strings.HasPrefix(trimmed, "it(") || strings.HasPrefix(trimmed, "test(") ||
			strings.HasPrefix(trimmed, "it('") || strings.HasPrefix(trimmed, "test('") ||
			strings.HasPrefix(trimmed, "it(\"") || strings.HasPrefix(trimmed, "test(\"") {
			fs.Tests++
		}
	case "Rust":
		if strings.HasPrefix(trimmed, "#[test]") || strings.HasPrefix(trimmed, "fn test_") {
			fs.Tests++
		}
	case "Java":
		if strings.HasPrefix(trimmed, "@Test") {
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

		return math.Round(score)
}
