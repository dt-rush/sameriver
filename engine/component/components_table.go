/*
 *
 * Allocates each of the big ass blocks of memory that each component has
 * its data living inside. This is a scientific terminology of game engine
 * design.
 *
 */

package component

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

func (t *ComponentsTable) ApplyComponentSet(id uint16, c ComponentSet) {
	// color
	if c.Color != nil {
		t.Color.SafeSet(id, *(c.Color))
	}
	// hitbox
	if c.Hitbox != nil {
		t.Hitbox.SafeSet(id, *(c.Hitbox))
	}
	// logic
	if c.Logic != nil {
		t.Logic.SafeSet(id, *(c.Logic))
	}
	// position
	if c.Position != nil {
		t.Position.SafeSet(id, *(c.Position))
	}
	// sprite
	if c.Sprite != nil {
		t.Sprite.SafeSet(id, *(c.Sprite))
	}
	// velocity
	if c.Velocity != nil {
		t.Velocity.SafeSet(id, *(c.Velocity))
	}
}
