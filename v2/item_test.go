package sameriver

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"

	"testing"
)

func TestItemFromArchetype(t *testing.T) {
	w := testingWorld()
	i := NewItemSystem(nil)
	w.RegisterSystems(i)

	i.CreateArchetype(map[string]any{
		"name":        "sword_iron",
		"displayName": "iron sword",
		"flavourText": "a good irons word, decently sharp",
		"properties": map[string]int{
			"damage":      3,
			"value":       20,
			"degradation": 0,
			"durability":  5,
		},
		"tags": []string{"weapon"},
	})

	i.CreateSubArchetype(map[string]any{
		"parent":      "sword_iron",
		"name":        "sword_iron_manjushris",
		"displayName": "manjushri's sword",
		"flavourText": "the sword of the legendary bodhisattva Manjushri; it can cut illusion itself",
		"tagDiff":     []string{"+legendary"},
	})

	manjushrisSword := i.CreateItem(map[string]any{
		"archetype": "sword_iron_manjushris",
	})

	Logger.Printf("Created: %s", manjushrisSword.String())

	if !manjushrisSword.Tags.Has("weapon") {
		t.Fatal("did not inherit tags!")
	}

	if manjushrisSword.GetProperty("damage") != 3 {
		t.Fatal("did not inherit properties!")
	}

	manjushrisSword.SetProperty("damage", 108)

	if manjushrisSword.GetProperty("damage") != 108 {
		t.Fatal("Did not set property!")
	}
}

func TestItemSystemLoadArchetypes(t *testing.T) {
	i := NewItemSystem(nil)
	i.LoadArchetypesFile("test_data/basic_archetypes.json")
	Logger.Println(i.Archetypes)
	if len(i.Archetypes) != 3 {
		t.Fatal("Did not load from JSON file!")
	}
	coin := i.Archetypes["coin_copper"]
	if len(coin.Entity) != 2 {
		t.Fatal("Did not load entity map of coin_copper")
	}
}

func TestItemSystemSpawnItemEntity(t *testing.T) {
	w := testingWorld()
	i := NewItemSystem(map[string]any{
		"spawn": true,
	})
	w.RegisterSystems(i)
	i.LoadArchetypesFile("test_data/basic_archetypes.json")
	coin := i.CreateItemSimple("coin_copper")
	coinEntity := i.SpawnItemEntity(Vec2D{10, 10}, coin)
	Logger.Println(coinEntity)
}

func TestItemSystemSpawnItemEntitySprite(t *testing.T) {
	skipCI(t)
	w := testingWorld()
	i := NewItemSystem(map[string]any{
		"spawn":  true,
		"sprite": true,
	})
	windowSpec := WindowSpec{
		Title:      "testing game",
		Width:      100,
		Height:     100,
		Fullscreen: false}
	// in a real game, the scene Init() gets a Game object and creates a new
	// sprite system by passing game.Renderer
	_, renderer := CreateWindowAndRenderer(windowSpec)
	sprites := NewSpriteSystem(renderer)

	w.RegisterSystems(i, sprites)
	i.LoadArchetypesFile("test_data/basic_archetypes.json")
	coin := i.CreateItemSimple("coin_copper")
	coinEntity := i.SpawnItemEntity(Vec2D{10, 10}, coin)
	Logger.Println(coinEntity)
	// draw the entity
	coinPos := coinEntity.GetVec2D("Position")
	coinBox := coinEntity.GetVec2D("Box")
	srcRect := sdl.Rect{0, 0, int32(coinBox.X), int32(coinBox.Y)}
	destRect := sdl.Rect{
		int32(coinPos.X),
		int32(coinPos.Y),
		int32(coinPos.X + coinBox.X),
		int32(coinPos.Y + coinBox.Y),
	}
	coinSprite := coinEntity.GetSprite("Sprite")
	renderer.Copy(coinSprite.Texture, &srcRect, &destRect)
	renderer.Present()
	time.Sleep(200 * time.Millisecond)
}

func TestItemSystemDespawnItemEntity(t *testing.T) {
	skipCI(t)
	w := testingWorld()
	i := NewItemSystem(map[string]any{
		"spawn":      true,
		"despawn_ms": 500,
	})

	w.RegisterSystems(i)
	i.LoadArchetypesFile("test_data/basic_archetypes.json")
	coin := i.CreateItemSimple("coin_copper")
	coinEntity := i.SpawnItemEntity(Vec2D{10, 10}, coin)
	Logger.Println(coinEntity)
	// draw the entity
	w.Update(FRAME_DURATION_INT / 2)
	time.Sleep(600 * time.Millisecond)
	w.Update(FRAME_DURATION_INT / 2)
	if !coinEntity.Despawned {
		t.Fatal("Should have despawned after time!")
	}
}
