package engine

type EntityComponent struct {
	entity     EntityToken
	components []ComponentType
}
