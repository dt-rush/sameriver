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
	Hitbox   *HitboxComponent
	Logic    *LogicComponent
	Position *PositionComponent
	Sprite   *SpriteComponent
	Velocity *VelocityComponent
}

func AllocateComponentsMemoryBlock() ComponentsTable {
	c := ComponentsTable{}
	c.Active = &ActiveComponent{}
	c.Color = &ColorComponent{}
	c.Hitbox = &HitboxComponent{}
	c.Logic = &LogicComponent{}
	c.Position = &PositionComponent{}
	c.Sprite = &SpriteComponent{}
	c.Velocity = &VelocityComponent{}
	return c
}

func (t *ComponentsTable) LinkEntityLocks(entityLocks *[MAX_ENTITIES]uint32) {
	t.Active.entityLocks = entityLocks
	t.Color.entityLocks = entityLocks
	t.Hitbox.entityLocks = entityLocks
	t.Logic.entityLocks = entityLocks
	t.Position.entityLocks = entityLocks
	t.Sprite.entityLocks = entityLocks
	t.Velocity.entityLocks = entityLocks
}

// NOTE: this must be called in a context in which the entity lock is preventing
// any reads or writes to the entity, or the gods will have mighty revenge on
// you for your hubris
func (t *ComponentsTable) ApplyComponentSet(id uint16, c ComponentSet) {
	// color
	if c.Color != nil {
		t.Color.Data[id] = *c.Color
	}
	// hitbox
	if c.Hitbox != nil {
		t.Hitbox.Data[id] = *c.Hitbox
	}
	// logic
	if c.Logic != nil {
		t.Logic.Data[id] = *c.Logic
	}
	// position
	if c.Position != nil {
		t.Position.Data[id] = *c.Position
	}
	// sprite
	if c.Sprite != nil {
		t.Sprite.Data[id] = *c.Sprite
	}
	// velocity
	if c.Velocity != nil {
		t.Velocity.Data[id] = *c.Velocity
	}
}
