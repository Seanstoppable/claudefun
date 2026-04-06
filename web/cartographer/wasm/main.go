package main

import (
	"encoding/json"
	"syscall/js"
)

func generateMap(this js.Value, args []js.Value) interface{} {
	placeName := args[0].String()
	width := 40
	height := 25
	svgWidth := 1000
	svgHeight := 700
	if len(args) > 1 {
		width = args[1].Int()
	}
	if len(args) > 2 {
		height = args[2].Int()
	}
	if len(args) > 3 {
		svgWidth = args[3].Int()
	}
	if len(args) > 4 {
		svgHeight = args[4].Int()
	}

	terrain := GenerateTerrain(placeName, width, height)
	lg := NewLandmarkGenerator(placeName)
	landmarks := lg.PlaceLandmarks(width, height, terrain.IsLand)

	svgR := NewSVGMapRenderer(svgWidth, svgHeight)
	svg := svgR.Render(terrain, landmarks, placeName)

	loreGen := NewLoreGenerator(placeName)
	lore := loreGen.Generate(placeName)

	landmarkList := make([]map[string]interface{}, len(landmarks))
	for i, lm := range landmarks {
		landmarkList[i] = map[string]interface{}{
			"name": lm.Name, "type": lm.Type.String(), "symbol": lm.Symbol,
		}
	}

	legendList := make([]map[string]interface{}, len(lore.Legends))
	for i, leg := range lore.Legends {
		legendList[i] = map[string]interface{}{
			"title": leg.Title, "category": leg.Category, "story": leg.Story,
		}
	}

	result := map[string]interface{}{
		"svg":       svg,
		"placeName": placeName,
		"landmarks": landmarkList,
		"lore": map[string]interface{}{
			"motto":        lore.Motto,
			"creationMyth": lore.CreationMyth,
			"legends":      legendList,
			"warning":      lore.Warning,
		},
	}
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		js.Global().Get("console").Call("error", "JSON marshal error:", err.Error())
		return `{"error":"internal encoding error"}`
	}
	return string(jsonBytes)
}

func main() {
	js.Global().Set("generateMap", js.FuncOf(generateMap))
	select {}
}
