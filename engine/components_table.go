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
	Hitbox   *HitboxComponent
	Logic    *LogicComponent
	Position *PositionComponent
	Sprite   *SpriteComponent
	Velocity *VelocityComponent
}

func AllocateComponentsMemoryBlock() ComponentsTable {
	ct := ComponentsTable{}
	ct.Active = &ActiveComponent{}
	ct.Color = &ColorComponent{}
	ct.Health = &HealthComponent{}
	ct.Hitbox = &HitboxComponent{}
	ct.Logic = &LogicComponent{}
	ct.Position = &PositionComponent{}
	ct.Sprite = &SpriteComponent{}
	ct.Velocity = &VelocityComponent{}
	return ct
}

func (ct *ComponentsTable) LinkEntityManager(em *EntityManager) {
	ct.Active.em = em
	ct.Color.em = em
	ct.Health.em = em
	ct.Hitbox.em = em
	ct.Logic.em = em
	ct.Position.em = em
	ct.Sprite.em = em
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
	if cs.Hitbox != nil {
		ct.Hitbox.Data[id] = *cs.Hitbox
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
	// velocity
	if cs.Velocity != nil {
		ct.Velocity.Data[id] = *cs.Velocity
	}
}
