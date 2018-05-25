package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dt-rush/donkeys-qquest/engine"
)

const PROFILE_EM_UPDATE = false
const N_CROWS = 2
const N_BEETLES = 1
const VERBOSE = true

func verboseLog(s string, params ...interface{}) {
	if VERBOSE {
		fmt.Printf(s, params...)
	}
}

func spawnCrows(crowClass *engine.EntityClass, em *engine.EntityManager) {
	for i := 0; i < N_CROWS; i++ {
		spawnCrow(crowClass, em)
	}
}

func spawnCrow(crowClass *engine.EntityClass, em *engine.EntityManager) {

	logic := engine.NewLogicUnit(
		"crow logic",
		crowClass.GenerateLogicFunc(
			[]engine.Behavior{CrowEatBeetleBehavior}))

	position := [2]int16{0, 0}

	spawnRequest := engine.EntitySpawnRequest{
		Components: engine.ComponentSet{
			Position: &position,
			Logic:    &logic,
		},
		Tags: []string{"bird", "crow"}}
	em.RequestSpawn(spawnRequest)
}

var CrowEatBeetleBehavior = engine.Behavior{
	Sleep: 500 * time.Millisecond,
	Func: func(crow engine.EntityToken,
		crowClass *engine.EntityClass,
		em *engine.EntityManager) {

		time.Sleep(
			time.Duration(rand.Intn(200)) * time.Millisecond)
		// the crow CAW's periodically
		verboseLog("Crow %d says, CAW!\n", crow.ID)
		// the crow examines the list of beetles and tries to eat one
		// 50% of the time
		fmt.Printf("Crow %d sees %d beetles\n",
			crow.ID, crowClass.Lists["beetle"].Length())
		beetle, err := crowClass.Lists["beetle"].GetFirst()
		if err == nil {
			verboseLog("Crow %d notices a tasty beetle: Beetle %d\n",
				crow.ID, beetle.ID)
			if rand.Intn(2) == 0 {
				verboseLog("Crow %d decides to eat the Beetle %d.\n",
					crow.ID, beetle.ID)
				didEat := em.AtomicEntityModify(
					beetle, func(e *engine.EntityModification) {
						e.Type = engine.ENTITY_STATE_MODIFICATION
						e.Modification = engine.ENTITY_DESPAWN
					})
				if didEat {
					verboseLog("Crow %d ate the delicious Beetle %d.\n",
						crow.ID, beetle.ID)
				} else {
					verboseLog("Crow %d couldn't eat Beetle %d, it was "+
						"gone by the time it got to it!\n",
						crow.ID, beetle.ID)
				}
			}
		}
		time.Sleep(
			time.Duration(rand.Intn(1000)) * time.Millisecond)
	}}

func spawnBeetles(beetleClass *engine.EntityClass, em *engine.EntityManager) {
	for i := 0; i < N_BEETLES; i++ {
		spawnBeetle(beetleClass, em)
	}
}

func spawnBeetle(beetleClass *engine.EntityClass, em *engine.EntityManager) {

	logic := engine.NewLogicUnit(
		"beetle logic",
		beetleClass.GenerateLogicFunc(
			[]engine.Behavior{BeetleBeepBehavior}))

	position := [2]int16{0, 0}

	spawnRequest := engine.EntitySpawnRequest{
		Components: engine.ComponentSet{
			Position: &position,
			Logic:    &logic,
		},
		Tags: []string{"insect", "beetle"}}
	em.RequestSpawn(spawnRequest)
}

var BeetleBeepBehavior = engine.Behavior{
	Sleep: 500 * time.Millisecond,
	Func: func(beetle engine.EntityToken,
		beetleClass *engine.EntityClass,
		em *engine.EntityManager) {

		verboseLog("beetle %d says, beep!\n", beetle.ID)
		time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
	}}

func main() {
	rand.Seed(time.Now().UnixNano())
	var t0 = time.Now()
	var em = engine.EntityManager{}
	em.Init()
	fmt.Printf("EntityManager.Init() took %d ms\n",
		time.Since(t0).Nanoseconds()/1e6)
	fmt.Println("Starting entity manager testing")

	var crowClass = engine.NewEntityClass(
		&em,
		"crow",
		[]engine.GenericEntityQuery{
			engine.GenericEntityQueryForTag("crow"),
			engine.GenericEntityQueryForTag("beetle"),
		})

	var beetleClass = engine.NewEntityClass(
		&em,
		"beetle",
		[]engine.GenericEntityQuery{
			engine.GenericEntityQueryForTag("beetle"),
			engine.GenericEntityQueryForTag("crow"),
		})

	spawnCrows(&crowClass, &em)
	spawnBeetles(&beetleClass, &em)

	for {
		if PROFILE_EM_UPDATE {
			t0 = time.Now()
		}
		em.Update()
		if PROFILE_EM_UPDATE {
			fmt.Printf("EntityManager.Update() took %d ms\n",
				time.Since(t0).Nanoseconds()/1e6)
		}
		time.Sleep(16 * time.Millisecond)
	}
}
