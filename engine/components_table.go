/*
 *
 * Allocates each of the big ass blocks of memory that each component has
 * its data living inside. This is a scientific terminology of game engine
 * design.
 *
 */

package engine

type ComponentsTable struct {
	Active   *ActiveComponent
	Color    *ColorComponent
	Health   *HealthComponent
	HitBox   *HitBoxComponent
	Logic    *LogicComponent
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
	ct.Active = &ActiveComponent{}
	ct.Color = &ColorComponent{}
	ct.Health = &HealthComponent{}
	ct.HitBox = &HitBoxComponent{}
	ct.Logic = &LogicComponent{}
	ct.Position = &PositionComponent{}
	ct.Sprite = &SpriteComponent{}
	ct.TagList = &TagListComponent{}
	ct.Velocity = &VelocityComponent{}
}

func (ct *ComponentsTable) linkEntityManager(
	em *EntityManager) {

	ct.Active.em = em
	ct.Color.em = em
	ct.Health.em = em
	ct.HitBox.em = em
	ct.Logic.em = em
	ct.Position.em = em
	ct.Sprite.em = em
	ct.TagList.em = em
	ct.Velocity.em = em
}

// NOTE: this must be called in a context in which the entity lock is preventing
// any reads or writes to the entity, or the gods will have mighty revenge on
// you for your hubris
func (ct *ComponentsTable) ApplyComponentSet(id int, cs ComponentSet) {
	// color
	if cs.Color != nil {
		ct.Color.Data[id] = *cs.Color
	}
	// health
	if cs.Health != nil {
		ct.Health.Data[id] = *cs.Health
	}
	// hitbox
	if cs.HitBox != nil {
		ct.HitBox.Data[id] = *cs.HitBox
	}
	// logic
	if cs.Logic != nil {
		ct.Logic.Data[id] = *cs.Logic
	}
	// position
	if cs.Position != nil {
		ct.Position.Data[id] = *cs.Position
	}
	// sprite
	if cs.Sprite != nil {
		ct.Sprite.Data[id] = *cs.Sprite
	}
	// taglist
	if cs.TagList != nil {
		ct.TagList.Data[id] = *cs.TagList
	}
	// velocity
	if cs.Velocity != nil {
		ct.Velocity.Data[id] = *cs.Velocity
	}
}
