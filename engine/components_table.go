/*
 *
 * Allocates each of the big ass blocks of memory that each component has
 * its data living inside. This is a scientific terminology of game engine
 * design.
 *
 */

// TODO: generate

package engine

type ComponentsTable struct {
	Color    *ColorComponent
	Health   *HealthComponent
	HitBox   *HitBoxComponent
	Logic    *LogicComponent
	Mind     *MindComponent
	Position *PositionComponent
	Sprite   *SpriteComponent
	TagList  *TagListComponent
	Velocity *VelocityComponent
}

func (ct *ComponentsTable) Init(em *EntityManager) {
	ct.allocate()
	ct.linkEntityManager(em)
}

func (ct *ComponentsTable) allocate() {
	ct.Color = &ColorComponent{}
	ct.Health = &HealthComponent{}
	ct.HitBox = &HitBoxComponent{}
	ct.Logic = &LogicComponent{}
	ct.Mind = &MindComponent{}
	ct.Position = &PositionComponent{}
	ct.Sprite = &SpriteComponent{}
	ct.TagList = &TagListComponent{}
	ct.Velocity = &VelocityComponent{}
}

func (ct *ComponentsTable) linkEntityManager(
	em *EntityManager) {

	ct.Color.em = em
	ct.Health.em = em
	ct.HitBox.em = em
	ct.Logic.em = em
	ct.Mind.em = em
	ct.Position.em = em
	ct.Sprite.em = em
	ct.TagList.em = em
	ct.Velocity.em = em
}
