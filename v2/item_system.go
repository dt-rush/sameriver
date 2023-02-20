package sameriver

import (
	"fmt"
	"os"

	"encoding/json"
	"io/ioutil"

	"github.com/dt-rush/sameriver/v2/utils"
)

func assureFloatMap(m any) map[string]float64 {
	var properties map[string]float64
	_, ok := m.(map[string]float64)
	if !ok {
		intProperties := m.(map[string]int)
		properties = make(map[string]float64)
		for k, v := range intProperties {
			properties[k] = float64(v)
		}
	} else {
		properties = m.(map[string]float64)
	}
	return properties
}

type ItemSystem struct {
	w               *World
	inventorySystem *InventorySystem `sameriver-system-dependency:"-"`
	spriteSystem    *SpriteSystem    `sameriver-system-dependency:"optional"`

	ItemEntities *UpdatedEntityList
	Archetypes   map[string]*ItemArchetype

	// how long as a time.Time item entities should last on the ground without
	// despawning (nil if not applicable)
	despawn_ms *float64

	// we update Degradations every second, so we need to use an accum of smaller dt_ms
	// values in Update()
	degradation_accum_ms float64

	// whether we will spawn item entities
	spawn bool
	// whether the spawned item entities have a sprite
	sprite bool
	// default box for spawned item entities
	defaultEntityBox Vec2D
}

func NewItemSystem(config map[string]any) *ItemSystem {
	i := &ItemSystem{
		Archetypes: make(map[string]*ItemArchetype),
	}

	if _, ok := config["spawn"]; ok {
		i.spawn = config["spawn"].(bool)
		if _, ok := config["defaultEntityBox"]; ok {
			i.defaultEntityBox = config["defaultEntityBox"].(Vec2D)
		}
	}

	if _, ok := config["sprite"]; ok {
		i.sprite = config["sprite"].(bool)
	}

	if _, ok := config["despawn_ms"]; ok {
		despawn_ms := float64(config["despawn_ms"].(int))
		i.despawn_ms = &despawn_ms
	}

	return i
}

func (i *ItemSystem) registerArchetype(arch *ItemArchetype) {
	i.Archetypes[arch.Name] = arch
}

func (i *ItemSystem) CreateArchetype(spec map[string]any) {

	var name, displayName, flavourText string
	var properties map[string]float64
	var tags []string
	var entity map[string]any

	if _, ok := spec["name"]; ok {
		name = spec["name"].(string)
	} else {
		panic("Must supply \"name\" to CreateArchetype")
	}
	if _, ok := spec["displayName"]; ok {
		displayName = spec["displayName"].(string)
	} else {
		displayName = name
	}
	if _, ok := spec["flavourText"]; ok {
		flavourText = spec["flavourText"].(string)
	} else {
		flavourText = ""
	}
	if _, ok := spec["properties"]; ok {
		properties = assureFloatMap(spec["properties"])
	}
	if _, ok := spec["tags"]; ok {
		tags = spec["tags"].([]string)
	}
	if _, ok := spec["entity"]; ok {
		entity = spec["entity"].(map[string]any)
	}

	tagList := NewTagList()
	tagList.Add(tags...)

	a := &ItemArchetype{
		Name:        name,
		DisplayName: displayName,
		FlavourText: flavourText,
		Properties:  properties,
		Tags:        tagList,
		Entity:      entity,
	}
	i.registerArchetype(a)
}

func (i *ItemSystem) CreateSubArchetype(spec map[string]any) {
	var parent string
	var name, displayName, flavourText string
	var properties map[string]float64
	var tagDiff []string
	var entity map[string]any

	if _, ok := spec["parent"]; ok {
		parent = spec["parent"].(string)
		if _, ok := i.Archetypes[parent]; !ok {
			panic(fmt.Sprintf("Archetype %s not found at time needed", parent))
		}
	} else {
		panic("Must supply \"parent\" to CreateSubArchetype")
	}
	if _, ok := spec["name"]; ok {
		name = spec["name"].(string)
	} else {
		panic("Must supply \"name\" to CreateSubArchetype")
	}
	if _, ok := spec["displayName"]; ok {
		displayName = spec["displayName"].(string)
	} else {
		displayName = i.Archetypes[parent].DisplayName
	}
	if _, ok := spec["flavourText"]; ok {
		flavourText = spec["flavourText"].(string)
	} else {
		flavourText = i.Archetypes[parent].FlavourText
	}
	if _, ok := spec["properties"]; ok {
		properties = assureFloatMap(spec["properties"])
	}
	if _, ok := spec["tagDiff"]; ok {
		tagDiff = spec["tagDiff"].([]string)
	}
	if _, ok := spec["entity"]; ok {
		entity = spec["entity"].(map[string]any)
	}

	a := &ItemArchetype{
		Name:        name,
		DisplayName: displayName,
		FlavourText: flavourText,
	}
	// copy in parent arch properties
	a.Properties = make(map[string]float64)
	for k, v := range i.Archetypes[parent].Properties {
		a.Properties[k] = v
	}
	// supplied properties shadow those of the arch parent
	for k, v := range properties {
		a.Properties[k] = v
	}
	// copy parent arch tags then apply the tagdiff specifications
	tags := i.Archetypes[parent].Tags.CopyOf()
	for _, tSpec := range tagDiff {
		op, t := tSpec[0:1], tSpec[1:]
		if op == "+" {
			tags.Add(t)
		} else if op == "-" {
			tags.Remove(t)
		}
	}
	a.Tags = tags
	// copy parent entity and then shadow with spec entity
	a.Entity = make(map[string]any)
	for k, v := range i.Archetypes[parent].Entity {
		a.Entity[k] = v
	}
	for k, v := range entity {
		a.Entity[k] = v
	}

	i.registerArchetype(a)
}

func (i *ItemSystem) CreateItem(spec map[string]any) *Item {
	// destructure params anymap
	var archetype string
	var properties map[string]float64
	var tags []string
	var degradationRate float64
	if _, ok := spec["archetype"]; ok {
		archetype = spec["archetype"].(string)
		if _, ok := i.Archetypes[archetype]; !ok {
			panic(fmt.Sprintf("Trying to create item of archetype that isn't created yet: %s", archetype))
		}
	} else {
		panic("Must specify \"archetype\" in CreateItem()")
	}
	if _, ok := spec["properties"]; ok {
		properties = assureFloatMap(spec["properties"])
	}
	if _, ok := spec["tags"]; ok {
		tags = spec["tags"].([]string)
	}
	if _, ok := spec["degradationRate"]; ok {
		degradationRate = spec["degradationRate"].(float64)
	}

	// shadow the archetypes properties with the params properties
	mergedProperties := make(map[string]float64)
	for k, v := range i.Archetypes[archetype].Properties {
		mergedProperties[k] = v
	}
	for k, v := range properties {
		mergedProperties[k] = v
	}

	tagList := NewTagList()
	tagList.Add(tags...)
	tagList.MergeIn(i.Archetypes[archetype].Tags)

	item := &Item{
		sys:        i,
		Archetype:  archetype,
		Properties: mergedProperties,
		Tags:       tagList,
		Count:      1,
	}

	if item.Tags.Has("perishable") {
		item.degradationRate = degradationRate
		item.Degradations = []float64{0}
	}

	item.reevaluateDisplayStr()

	return item
}

func (i *ItemSystem) CreateStack(n int, spec map[string]any) *Item {
	item := i.CreateItem(spec)
	// make the single item a stack
	stack := item
	stack.Count = n
	stack.reevaluateDisplayStr()
	if item.Tags.Has("perishable") {
		stack.Degradations = make([]float64, n)
		for i, _ := range stack.Degradations {
			stack.Degradations[i] = 0
		}
	}
	return stack
}

func (i *ItemSystem) CreateItemSimple(archetype string) *Item {
	return i.CreateItem(map[string]any{
		"archetype": archetype,
	})
}

func (i *ItemSystem) CreateStackSimple(n int, archetype string) *Item {
	return i.CreateStack(n, map[string]any{
		"archetype": archetype,
	})
}

func (i *ItemSystem) SpawnItemEntity(pos Vec2D, item *Item) *Entity {
	entityBox := i.defaultEntityBox
	arch := item.GetArchetype()
	if _, ok := arch.Entity["box"]; ok {
		box := arch.Entity["box"].([2]float64)
		entityBox = Vec2D{box[0], box[1]}
	}
	// the item entity will have an inventory containing the item
	// (and any unstacked rotting items)
	inventory := NewInventory()
	inventory.Credit(item)
	components := map[string]any{
		"Vec2D,Position":    pos,
		"Vec2D,Box":         entityBox,
		"Generic,Inventory": inventory,
	}
	if i.sprite {
		if i.spriteSystem == nil {
			panic("Trying to create entity with sprite=true while spriteSystem was not registered")
		}
		components["Sprite,Sprite"] = i.spriteSystem.GetSprite(arch.Entity["sprite"].(string))
	}
	if i.despawn_ms != nil {
		components["TimeAccumulator,DespawnTimer"] = utils.NewTimeAccumulator(*i.despawn_ms)
	}

	return i.w.Spawn(map[string]any{
		"components": components,
		"tags":       []string{"item"},
	})
}

func (i *ItemSystem) LoadArchetypesFile(filename string) {
	Logger.Printf("Loading item archetypes from %s...", filename)
	jsonFile, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("Trying to open %s - doesn't exist", filename))
	}
	contents, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	i.LoadArchetypesJSON(contents)
}

func (i *ItemSystem) LoadArchetypesJSON(jsonStr []byte) {
	var archetypes [](map[string]*json.RawMessage)
	err := json.Unmarshal(jsonStr, &archetypes)
	if err != nil {
		panic(err)
	}
	for ix, jsonSpec := range archetypes {
		spec := make(map[string]any)
		if _, ok := jsonSpec["name"]; ok {
			var name string
			json.Unmarshal(*jsonSpec["name"], &name)
			spec["name"] = name
		} else {
			panic(fmt.Sprintf("object at index %d was missing \"name\" property", ix))
		}
		if _, ok := jsonSpec["displayName"]; ok {
			var displayName string
			json.Unmarshal(*jsonSpec["displayName"], &displayName)
			spec["displayName"] = displayName
		}
		if _, ok := jsonSpec["flavourText"]; ok {
			var flavourText string
			json.Unmarshal(*jsonSpec["flavourText"], &flavourText)
			spec["flavourText"] = flavourText
		}
		if _, ok := jsonSpec["properties"]; ok {
			var properties map[string]float64
			json.Unmarshal(*jsonSpec["properties"], &properties)
			spec["properties"] = properties
		}
		if _, ok := jsonSpec["tags"]; ok {
			var tags []string
			json.Unmarshal(*jsonSpec["tags"], &tags)
			spec["tags"] = tags
		}
		if _, ok := jsonSpec["entity"]; ok {
			entity := make(map[string]any)
			var entityMap map[string]*json.RawMessage
			json.Unmarshal(*jsonSpec["entity"], &entityMap)
			if _, ok := entityMap["sprite"]; ok {
				var sprite string
				json.Unmarshal(*entityMap["sprite"], &sprite)
				entity["sprite"] = sprite
			}
			if _, ok := entityMap["box"]; ok {
				var box [2]float64
				json.Unmarshal(*entityMap["box"], &box)
				entity["box"] = box
			}
			spec["entity"] = entity
		}
		i.CreateArchetype(spec)
	}
}

// System funcs

func (i *ItemSystem) GetComponentDeps() []string {
	deps := []string{}
	if i.spawn {
		deps = append(deps, "Generic,Item")
		deps = append(deps, "Vec2D,Position")
		deps = append(deps, "Vec2D,Box")
	}
	if i.sprite {
		deps = append(deps, "Sprite,Sprite")
	}
	if i.despawn_ms != nil {
		deps = append(deps, "TimeAccumulator,DespawnTimer")
	}
	return deps
}

func (i *ItemSystem) LinkWorld(w *World) {
	i.w = w

	i.ItemEntities = w.em.GetSortedUpdatedEntityList(
		EntityFilterFromTag("item"))
}

func (i *ItemSystem) Update(dt_ms float64) {
	// despawn any expired entities
	if i.despawn_ms != nil {
		for _, e := range i.ItemEntities.entities {
			accum := e.GetTimeAccumulator("DespawnTimer")
			if accum.Tick(dt_ms) {
				i.w.Despawn(e)
			}
		}
	}

	// if a second or more has elapsed, update Degradations
	i.degradation_accum_ms += dt_ms
	if i.degradation_accum_ms > 1000 {
		for _, e := range i.inventorySystem.InventoryEntities.entities {
			inv := e.GetGeneric("Inventory").(*Inventory)
			newRotStacks := make([]*Item, 0)
			for _, s := range inv.Stacks {
				if s.Tags.Has("perishable") {
					// we can do this nicely since s.Degradations is always sorted
					// (debit splits evenly along the sorted slice and credit resorts
					// after appending)
					var rotten bool
					var rottenIx int
					for ix, d := range s.Degradations {
						/*
							0.5 degradation per second means
							0.5 / 1000 = 0.0005 degradation per ms
						*/
						d_per_ms := s.degradationRate / 1000
						s.Degradations[ix] = d + d_per_ms*i.degradation_accum_ms
						if s.Degradations[ix] >= 100 {
							rotten = true
							rottenIx = ix
							continue
						}
					}
					if rotten && !s.Tags.Has("rotten") {
						// we use the inv interface even for debit so maybe, if you wanted
						// to put a logger there, you'd see these as debits and credits too
						rottenStack := inv.DebitNWithPreference(s, s.Count-rottenIx, ITEM_MOST_DEGRADED)
						rottenStack.Tags.Add("rotten")
						newRotStacks = append(newRotStacks, rottenStack)

					}
				}
			}
			for _, s := range newRotStacks {
				// we obviously must credit this stack as its now different, and has
				// been removed in quantity from the existing stack in inv
				inv.Credit(s)
			}
		}
		i.degradation_accum_ms = 0
	}

}
