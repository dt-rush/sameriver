/*
 * EntityClass is used to group a set of entities by the EntityManager,
 * and to allow classes of entities (eg. "bear", "crow") to define themselves
 *
 */

package engine

type EntityClass interface {
	Name() string
	DefaultSpawnRequest(position [2]int32) EntitySpawnRequest
}
