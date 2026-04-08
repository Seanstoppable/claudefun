package main

import (
	"encoding/json"
	"syscall/js"
)

func analyzeCode(this js.Value, args []js.Value) interface{} {
	code := args[0].String()
	language := "go"
	if len(args) > 1 {
		language = args[1].String()
	}

	stats := analyzeString(code, language)
	garden := BuildGarden(stats)
	report := generateReport(garden, stats)

	// Build grid as array of rows, each row is array of cells
	grid := make([][]map[string]interface{}, garden.Height)
	for y := 0; y < garden.Height; y++ {
		row := make([]map[string]interface{}, garden.Width)
		for x := 0; x < garden.Width; x++ {
			plot := garden.Plots[y][x]
			sym := plot.Element.Symbol()
			color := plot.Element.Color()
			if plot.Element == Grass {
				sym = plot.Plant.Symbol()
				color = plot.Plant.Color()
			}
			row[x] = map[string]interface{}{
				"symbol": sym,
				"color":  color,
				"label":  plot.Label,
			}
		}
		grid[y] = row
	}

	result := map[string]interface{}{
		"title":       garden.Title,
		"width":       garden.Width,
		"height":      garden.Height,
		"weather":     garden.Weather,
		"season":      garden.Season,
		"healthScore": garden.HealthScore,
		"grid":        grid,
		"stats": map[string]interface{}{
			"trees":       garden.Stats.Trees,
			"flowers":     garden.Stats.Flowers,
			"fences":      garden.Stats.Fences,
			"bugs":        garden.Stats.Bugs,
			"butterflies": garden.Stats.Butterflies,
			"rocks":       garden.Stats.Rocks,
			"weeds":       garden.Stats.Weeds,
		},
		"code": map[string]interface{}{
			"totalFiles":   stats.TotalFiles,
			"totalLines":   stats.TotalLines,
			"totalFuncs":   stats.TotalFuncs,
			"totalTypes":   stats.TotalTypes,
			"totalTests":   stats.TotalTests,
			"totalComments": stats.TotalComments,
			"totalTODOs":   stats.TotalTODOs,
			"complexity":   stats.TotalComplexity,
		},
		"report": report,
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		js.Global().Get("console").Call("error", "JSON marshal error:", err.Error())
		return `{"error":"internal encoding error"}`
	}
	return string(jsonBytes)
}

func main() {
	js.Global().Set("analyzeCode", js.FuncOf(analyzeCode))
	select {}
}
