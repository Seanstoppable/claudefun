package main

import (
	"encoding/json"
	"syscall/js"
)

func generateConstellation(this js.Value, args []js.Value) interface{} {
	input := args[0].String()
	width := 800
	height := 600
	if len(args) > 1 {
		width = args[1].Int()
	}
	if len(args) > 2 {
		height = args[2].Int()
	}

	sm := GenerateStarMap(input)
	cons := Connect(sm)
	svgR := NewSVGRenderer(width, height)
	svg := svgR.Render(cons)

	mg := NewMythologyGenerator(sm.Seed)
	myth := mg.Generate(sm.ConstellationName(), len(sm.Stars))

	result := map[string]interface{}{
		"svg":   svg,
		"name":  sm.ConstellationName(),
		"stars": len(sm.Stars),
		"edges": len(cons.Edges),
		"mythology": map[string]interface{}{
			"culture":     myth.Culture,
			"story":       myth.Story,
			"moral":       myth.Moral,
			"bestViewing": myth.BestViewing,
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
	js.Global().Set("generateConstellation", js.FuncOf(generateConstellation))
	select {}
}
