package engine

type EntityModificationRequest struct {
	entity     EntityToken
	components []ComponentType
}
