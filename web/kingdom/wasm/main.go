package main

import (
	"encoding/json"
	"syscall/js"
	"time"
)

var (
	currentKingdom *Kingdom
	policyGen      *PolicyGenerator
	eventGen       *EventGenerator
	bard           *Bard
	currentPolicy  Policy
)

// marshalJSON converts a Go value to a JS-parseable JSON string.
func marshalJSON(v interface{}) interface{} {
	b, err := json.Marshal(v)
	if err != nil {
		return js.ValueOf(map[string]interface{}{"error": err.Error()})
	}
	return js.ValueOf(string(b))
}

func kingdomStateMap() map[string]interface{} {
	factions := map[string]interface{}{}
	for _, f := range AllFactions {
		factions[string(f)] = currentKingdom.FactionMood[f]
	}

	return map[string]interface{}{
		"name":       currentKingdom.Name,
		"turn":       currentKingdom.Turn,
		"treasury":   currentKingdom.Treasury,
		"population": currentKingdom.Population,
		"happiness":  currentKingdom.Happiness,
		"military":   currentKingdom.Military,
		"culture":    currentKingdom.Culture,
		"food":       currentKingdom.Food,
		"reputation": currentKingdom.Reputation,
		"factions":   factions,
		"rulerTitle": currentKingdom.RulerTitle,
		"isStable":   currentKingdom.IsStable(),
	}
}

func policyMap() map[string]interface{} {
	return map[string]interface{}{
		"question": currentPolicy.Question,
		"optionA":  currentPolicy.OptionA,
		"optionB":  currentPolicy.OptionB,
	}
}

func newGame(_ js.Value, args []js.Value) interface{} {
	name := ""
	if len(args) > 0 {
		name = args[0].String()
	}
	seed := time.Now().UnixNano()
	currentKingdom = NewKingdom(name)
	policyGen = NewPolicyGenerator(seed)
	eventGen = NewEventGenerator(seed)
	bard = NewBard(seed, currentKingdom.Name)
	currentPolicy = policyGen.NextPolicy()

	result := map[string]interface{}{
		"type":      "welcome",
		"narration": bard.NarrateWelcome(currentKingdom.Name),
		"kingdom":   kingdomStateMap(),
		"policy":    policyMap(),
	}
	return marshalJSON(result)
}

func choosePolicy(_ js.Value, args []js.Value) interface{} {
	chooseA := args[0].Bool()

	var effect Effect
	var flavor string
	var chosen string
	if chooseA {
		effect = currentPolicy.EffectA
		flavor = currentPolicy.FlavorA
		chosen = currentPolicy.OptionA
	} else {
		effect = currentPolicy.EffectB
		flavor = currentPolicy.FlavorB
		chosen = currentPolicy.OptionB
	}

	currentKingdom.ApplyEffect(effect)
	narration := bard.NarratePolicy(currentPolicy.Question, chosen, flavor)

	// Check for random event
	event := eventGen.MaybeEvent(currentKingdom)
	var eventData map[string]interface{}
	if event != nil {
		currentKingdom.ApplyEffect(event.Effect)
		eventNarration := bard.NarrateEvent(event.Description)
		eventData = map[string]interface{}{
			"name":        event.Name,
			"description": event.Description,
			"narration":   eventNarration,
		}
	}

	currentKingdom.AdvanceTurn()
	gameOver := currentKingdom.CheckGameOver()

	var gameOverData map[string]interface{}
	if gameOver {
		gameOverData = map[string]interface{}{
			"victory":   currentKingdom.Victory,
			"message":   currentKingdom.GameOverMsg,
			"narration": bard.NarrateGameOver(currentKingdom.Victory, currentKingdom.GameOverMsg),
		}
	} else {
		currentPolicy = policyGen.NextPolicy()
	}

	result := map[string]interface{}{
		"type":      "turn",
		"narration": narration,
		"kingdom":   kingdomStateMap(),
		"turnStart": bard.NarrateTurnStart(currentKingdom.Turn),
	}

	if eventData != nil {
		result["event"] = eventData
	}
	if gameOverData != nil {
		result["gameOver"] = gameOverData
	} else {
		result["policy"] = policyMap()
	}

	return marshalJSON(result)
}

func getState(_ js.Value, _ []js.Value) interface{} {
	if currentKingdom == nil {
		return marshalJSON(map[string]interface{}{"error": "no game in progress"})
	}
	result := map[string]interface{}{
		"kingdom": kingdomStateMap(),
		"policy":  policyMap(),
	}
	return marshalJSON(result)
}

func main() {
	js.Global().Set("newGame", js.FuncOf(newGame))
	js.Global().Set("choosePolicy", js.FuncOf(choosePolicy))
	js.Global().Set("getState", js.FuncOf(getState))

	// Keep the Go program running
	select {}
}
