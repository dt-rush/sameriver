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
