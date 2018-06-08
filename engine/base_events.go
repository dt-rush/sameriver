package engine

type CollisionEvent struct {
	EntityA EntityToken
	EntityB EntityToken
}

type GenericEvent struct {
	Type int
	Data interface{}
}
