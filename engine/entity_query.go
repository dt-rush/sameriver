package engine

type EntityQuery interface {
	Test(entity EntityToken, em *EntityManager) bool
}

type GenericEntityQuery struct {
	Name     string
	TestFunc func(entity EntityToken, em *EntityManager) bool
}

func (q GenericEntityQuery) Test(entity EntityToken, em *EntityManager) bool {
	return q.TestFunc(entity, em)
}

func GenericEntityQueryFromTag(tag string) GenericEntityQuery {
	return GenericEntityQuery{
		Name: tag,
		TestFunc: func(entity EntityToken, em *EntityManager) bool {
			return em.EntityHasTag(entity, tag)
		}}
}
