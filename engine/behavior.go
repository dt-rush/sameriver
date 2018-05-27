/*
 * Behavior represents a certain logic which will be run for an entity
 * with a sleep time between each run. The entity is inherently a member
 * of an entity class, even if it's a singleton.
 *
 */

package engine

import (
	"time"
)

// the type of a function run
type BehaviorFunc func(
	e EntityToken,
	c *EntityClass,
	em *EntityManager)

type Behavior struct {
	Name string
	// a constant amount of time to sleep after each time Func is run
	Sleep time.Duration
	// the function this behaviour represents (run when running is 0)
	Func BehaviorFunc
	// used atomically as a lock to determine whether to run the Func
	running uint32
}
