package engine

type EntityQuery interface {
	Test(id uint16, em *EntityManager) bool
}

type GenericEntityQuery struct {
	Name     string
	TestFunc func(id uint16, em *EntityManager) bool
}

func (q GenericEntityQuery) Test(id uint16, em *EntityManager) bool {
	return q.TestFunc(id, em)
}

func GenericEntityQueryForTag(tag string) GenericEntityQuery {
	return GenericEntityQuery{
		Name: tag,
		TestFunc: func(id uint16, em *EntityManager) bool {
			return em.EntityHasTag(id, tag)
		}}
}
