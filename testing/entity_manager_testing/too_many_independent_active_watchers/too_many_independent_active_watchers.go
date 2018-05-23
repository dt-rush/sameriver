package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dt-rush/donkeys-qquest/engine"
)

const PROFILE_EM_UPDATE = true
const N_CROWS = 10
const N_BEETLES = 20
const VERBOSE = false

func verboseLog(s string, params ...interface{}) {
	if VERBOSE {
		fmt.Printf(s, params)
	}
}

func spawnCrows(em *engine.EntityManager) {
	for i := 0; i < N_CROWS; i++ {
		spawnCrow(em)
	}
}

func spawnCrow(em *engine.EntityManager) {
	logic := engine.NewLogicUnit(crowLogic, "crow logic")
	position := [2]int16{0, 0}

	spawnRequest := engine.EntitySpawnRequest{
		Components: engine.ComponentSet{
			Position: &position,
			Logic:    &logic,
		},
		Tags: []string{"bird", "crow"}}
	em.RequestSpawn(spawnRequest)
}

func crowLogic(crowID uint16,
	StopChannel chan bool,
	em *engine.EntityManager) {

	beetleQuery := engine.GenericEntityQuery{
		func(id uint16, em *engine.EntityManager) bool {
			return em.EntityHasTag(id, "beetle")
		}}
	beetleList := em.GetUpdatedActiveEntityList(beetleQuery,
		fmt.Sprintf("crow %d's beetle list", crowID))
logicloop:
	for {
		select {
		case <-StopChannel:
			verboseLog("crow %d logic ending", crowID)
			break logicloop
		default:
			time.Sleep(
				time.Duration(rand.Intn(200)) * time.Millisecond)
			// the crow CAW's periodically
			verboseLog("Crow %d says, CAW!\n", crowID)
			// the crow examines the list of beetles and tries to eat one
			// 10% of the time
			beetleList.Mutex.Lock()
			if len(beetleList.Entities) > 0 {
				verboseLog("Crow %d notices a tasty beetle\n", crowID)
				if rand.Intn(10) == 0 {
					verboseLog("Crow %d decides to eat the beetle it sees.\n",
						crowID)
					beetleID := beetleList.Entities[0]
					em.AtomicEntityModify(
						beetleID, func(e *engine.EntityModification) {

							verboseLog("Crow %d ate a delicious beetle.\n",
								crowID)
							e.Type = engine.ENTITY_STATE_MODIFICATION
							e.Modification = engine.ENTITY_DESPAWN
						})
				}
			}
			beetleList.Mutex.Unlock()
			time.Sleep(
				time.Duration(rand.Intn(3000)) * time.Millisecond)
		}
	}
}

func spawnBeetles(em *engine.EntityManager) {
	for i := 0; i < N_BEETLES; i++ {
		spawnBeetle(em)
	}
}

func spawnBeetle(em *engine.EntityManager) {
	logic := engine.NewLogicUnit(beetleLogic, "beetle logic")
	position := [2]int16{0, 0}

	spawnRequest := engine.EntitySpawnRequest{
		Components: engine.ComponentSet{
			Position: &position,
			Logic:    &logic,
		},
		Tags: []string{"insect", "beetle"}}
	em.RequestSpawn(spawnRequest)
}

func beetleLogic(beetleID uint16,
	StopChannel chan bool,
	em *engine.EntityManager) {
logicloop:
	for {
		select {
		case <-StopChannel:
			break logicloop
		default:
			time.Sleep(
				time.Duration(rand.Intn(200)) * time.Millisecond)
			verboseLog("beetle %d says, beep!\n", beetleID)
			time.Sleep(
				time.Duration(rand.Intn(3000)) * time.Millisecond)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	t0 := time.Now()
	em := engine.EntityManager{}
	em.Init()
	fmt.Printf("EntityManager.Init() took %d ms\n",
		time.Since(t0).Nanoseconds()/1e6)
	fmt.Println("Starting entity manager testing")

	spawnCrows(&em)
	spawnBeetles(&em)

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
