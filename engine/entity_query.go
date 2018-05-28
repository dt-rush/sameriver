package engine

type EntityQuery struct {
	Name     string
	TestFunc func(entity EntityToken, em *EntityManager) bool
}

func (q EntityQuery) Test(entity EntityToken, em *EntityManager) bool {
	return q.TestFunc(entity, em)
}

func EntityQueryFromTag(tag string) EntityQuery {

	return EntityQuery{
		Name: tag,
		TestFunc: func(entity EntityToken, em *EntityManager) bool {
			return em.EntityHasTag(entity, tag)
		}}
}

func EntityQueryFromComponentBitArray(
	q bitarray.BitArray) EntityQuery {

	return EntityQuery{
		Name: BitArrayToString(q),
		TestFunc: func(entity EntityToken, em *EntityManager) bool {
			// determine if q = q&b
			// that is, if every set bit of q is set in b
			b := em.entityComponentBitArray(entity.ID)
			return q.Equals(q.And(b))
		}}
}
