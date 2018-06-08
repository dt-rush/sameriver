package engine

type OneEntityComponents struct {
	entity     EntityToken
	components []ComponentType
}

func (ecs *OneEntityComponents) Flatten() []EntityComponent {
	flattened := make([]EntityComponent, len(ecs.components))
	for i, component := range ecs.components {
		flattened[i] = EntityComponent{ecs.entity, component}
	}
	return flattened
}

type EntityComponent struct {
	entity    EntityToken
	component ComponentType
}
