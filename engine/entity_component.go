package engine

type EntityComponents struct {
	Entity     EntityToken
	Components []ComponentType
}

func (ecs *EntityComponents) Flatten() []EntityComponent {
	flattened := make([]EntityComponent, len(ecs.Components))
	for i, component := range ecs.Components {
		flattened[i] = EntityComponent{ecs.Entity, component}
	}
	return flattened
}

type EntityComponent struct {
	Entity    EntityToken
	Component ComponentType
}

// used to conveniently call AtomicEntityModify
func ECs(EClist ...EntityComponents) []EntityComponents {
	return EClist
}

// used to conveniently call AtomicEntitiesModify
func EC(
	entity EntityToken, components ...ComponentType) EntityComponents {

	return EntityComponents{entity, components}
}
