package engine

// Each LogicFunc will started as a goroutine, supplied with the EntityToken
// of the entity it's attached to, a channel on which a stop signal may
// arrive, and a reference to the EntityManager
//
// Through the EntityManager, the goroutine will be able to:
//
// - request an UpdatedEntityList with an arbitrary EntityQuery
// - make a one-time query with EntityManager.RunQuery()
// - get entity component data via EntityManager.Components.Read${Component}()
// - send EntitySpawnRequest messages to the SpawnChannel
// - atomically modify entities using em.AtomicEntit(y|ies)Modify()
//
type EntityLogicFunc func(
	entity EntityToken,
	StopChannel chan bool,
	em *EntityManager)
