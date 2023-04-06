package sameriver

import (
	"fmt"
	"strings"
	"testing"
)

func TestEntityFilterDSLLexer(t *testing.T) {
	lex := func(s string) {
		fmt.Println(s)
		var l EntityFilterDSLLexer
		l.Init(strings.NewReader(s))
		for tok := l.Lex(); tok != EOF; tok = l.Lex() {
			fmt.Printf("%s: %s\n", tok, l.TokenText())
		}
		fmt.Println()
	}

	lex(`HasTag(ox) && CanBe(yoked, 1); Closest(mind.field)`)
	lex(`First(HasTag(ox) && CanBe(yoked, 1); Closest(mind.field))`)
	lex(`VillagerOf(self.village)`)
	lex(`!VillagerOf(self.village)`)
	lex(`Is(bb.village1.strongest)`)
	lex(`HasTag(deer); Closest(self)`)
	lex(`HasTags(ox,legendary)`)
	lex(`Closest(self)`)
}

func TestEntityFilterDSLParser(t *testing.T) {
	parser := &EntityFilterDSLParser{}
	ast, err := parser.Parse(`HasTag(ox) && CanBe(yoked, 1); Closest(mind.field)`)
	if err != nil {
		t.Fatalf("Why did the expression return an error? it's valid!")
	}
	Logger.Printf("%s", ast)
}

func TestEntityFilterDSLEvaluator(t *testing.T) {
	w := testingWorld()
	ox := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION: Vec2D{0, 0},
			BOX:      Vec2D{3, 2},
			STATE: map[string]int{
				"yoked": 0,
			},
		},
		"tags": []string{"ox"},
	})

	// Initialize parser and evaluator with your custom function maps
	parser := &EntityFilterDSLParser{}
	evaluator := NewEntityFilterDSLEvaluator(EntityFilterDSLPredicates, EntityFilterDSLSorts)

	expression := "HasTags(ox)"

	// Parse and evaluate the expression
	ast, err := parser.Parse(expression)
	if err != nil {
		t.Fatalf("Failed to parse expression: %s", err)
	}

	resolver := &EntityResolver{e: ox}
	filter, _ := evaluator.Evaluate(ast, resolver)

	// Filter entities using the generated filter function
	result := w.FilterAllEntities(filter)

	// Check if the filtered list contains the expected ox entity
	if len(result) != 1 || result[0] != ox {
		t.Fatalf("Failed to select ox entity: got %v", result)
	}
	Logger.Printf("result of HasTags(ox): %v", result)
}
