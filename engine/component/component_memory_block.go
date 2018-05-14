/*
 *
 * Allocates each of the big ass blocks of memory that each component has
 * its data living inside. This is a scientific terminology of game engine
 * design.
 *
 */

package component

type ComponentsTable struct {
	active   *ActiveComponent
	color    *ColorComponent
	hitbox   *HitboxComponent
	logic    *LogicComponent
	position *PositionComponent
	sprite   *SpriteComponent
	velocity *VelocityComponent
}

func AllocateComponentsMemoryBlock() ComponentsTable {
	c := ComponentsTable{}
	// allocation is done by the fact that we're instantiating structs
	// whose data members are static-sized arrays of component data
	// (eg. [MAX_ENTITIES][2]uint16, for position_component)
	c.active = &ActiveComponent{}
	c.color = &ColorComponent{}
	c.hitbox = &HitboxComponent{}
	c.logic = &LogicComponent{}
	c.position = &PositionComponent{}
	c.sprite = &SpriteComponent{}
	c.velocity = &VelocityComponent{}
	return c
}
