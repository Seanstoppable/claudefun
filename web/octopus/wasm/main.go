package main

import (
	"encoding/json"
	"syscall/js"
)

var advisor = NewAdvisor()

func analyzeMood(this js.Value, args []js.Value) interface{} {
	input := args[0].String()
	analyzer := NewAnalyzer()
	results := analyzer.Analyze(input)
	dominant := results[0]
	info := dominant.Emotion.Info()

	allResults := make([]map[string]interface{}, len(results))
	for i, r := range results {
		ri := r.Emotion.Info()
		allResults[i] = map[string]interface{}{
			"name":       ri.Name,
			"emoji":      ri.Emoji,
			"color":      ri.Color,
			"confidence": r.Confidence,
		}
	}

	shouldAdvise := advisor.ShouldGiveAdvice()
	advice := ""
	if shouldAdvise {
		advice = advisor.GetAdvice(dominant.Emotion)
	}

	result := map[string]interface{}{
		"dominant": map[string]interface{}{
			"name":       info.Name,
			"emoji":      info.Emoji,
			"color":      info.Color,
			"confidence": dominant.Confidence,
		},
		"all":    allResults,
		"eyes":   GetEyes(dominant.Emotion),
		"advice": advice,
	}
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		js.Global().Get("console").Call("error", "JSON marshal error:", err.Error())
		return `{"error":"internal encoding error"}`
	}
	return string(jsonBytes)
}

func getAdvice(this js.Value, args []js.Value) interface{} {
	emotionIdx := args[0].Int()
	emotion := Emotion(emotionIdx)
	return advisor.GetAdvice(emotion)
}

func getAllEmotions(this js.Value, args []js.Value) interface{} {
	emotions := AllEmotions()
	result := make([]map[string]interface{}, len(emotions))
	for i, e := range emotions {
		info := e.Info()
		result[i] = map[string]interface{}{
			"index": int(e),
			"name":  info.Name,
			"emoji": info.Emoji,
			"color": info.Color,
			"eyes":  GetEyes(e),
		}
	}
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		js.Global().Get("console").Call("error", "JSON marshal error:", err.Error())
		return `{"error":"internal encoding error"}`
	}
	return string(jsonBytes)
}

func main() {
	js.Global().Set("analyzeMood", js.FuncOf(analyzeMood))
	js.Global().Set("getAdvice", js.FuncOf(getAdvice))
	js.Global().Set("getAllEmotions", js.FuncOf(getAllEmotions))

	// Keep the Go program running
	select {}
}
