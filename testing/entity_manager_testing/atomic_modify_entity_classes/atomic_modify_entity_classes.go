package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dt-rush/donkeys-qquest/engine"
)

const PROFILE_EM_UPDATE = false
const N_CROWS = 20
const N_BEETLES = 16
const VERBOSE = true

func verboseLog(s string, params ...interface{}) {
	if VERBOSE {
		fmt.Printf(s, params...)
	}
}

type CrowEntityClass struct {
	crows   engine.UpdatedEntityList
	beetles engine.UpdatedEntityList
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
			// 50% of the time
			beetleList.Mutex.Lock()
			fmt.Printf("%s: %v\n", beetleList.Name, beetleList.Entities)
			if len(beetleList.Entities) > 0 {
				beetle := beetleList.Entities[0]
				verboseLog("Crow %d notices a tasty beetle: Beetle %d\n",
					crowID, beetle.ID)
				if rand.Intn(2) == 0 {
					verboseLog("Crow %d decides to eat the Beetle %d.\n",
						crowID, beetle.ID)
					didEat := em.AtomicEntityModify(
						beetle, func(e *engine.EntityModification) {
							e.Type = engine.ENTITY_STATE_MODIFICATION
							e.Modification = engine.ENTITY_DESPAWN
						})
					if didEat {
						verboseLog("Crow %d ate the delicious Beetle %d.\n",
							crowID, beetle.ID)
					} else {
						verboseLog("Crow %d couldn't eat Beetle %d, it was "+
							"gone by the time it got to it!\n",
							crowID, beetle.ID)
					}
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
