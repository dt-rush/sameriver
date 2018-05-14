/*
 *
 * Allocates each of the big ass blocks of memory that each component has
 * its data living inside. This is a scientific terminology of game engine
 * design.
 *
 */

package engine

import (
	"github.com/dt-rush/donkeys-qquest/engine/component"
)

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
	c.active = &ActiveComponent{}
	c.color = &ColorComponent{}
	c.hitbox = &HitboxComponent{}
	c.logic = &LogicComponent{}
	c.position = &PositionComponent{}
	c.sprite = &SpriteComponent{}
	c.velocity = &VelocityComponent{}
	return c
}
