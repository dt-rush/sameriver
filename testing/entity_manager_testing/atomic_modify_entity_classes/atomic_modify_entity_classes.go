package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dt-rush/donkeys-qquest/engine"
	"github.com/dt-rush/donkeys-qquest/entity/beetle"
	"github.com/dt-rush/donkeys-qquest/entity/crow"
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

func main() {
	rand.Seed(time.Now().UnixNano())
	var em engine.EntityManager
	em.Init()
	fmt.Println("em.Init() finished.")
	var ev engine.EventBus
	ev.Init()
	fmt.Println("ev.Init() finished.")
	var wl engine.WorldLogicManager
	wl.Init(&em, &ev)
	fmt.Println("wl.Init() finished.")
	fmt.Println("-------------------")

	wl.AddList(
		"healthful",
		engine.NewEntityComponentBitArrayQuery(
			engine.MakeComponentBitArray(
				[]int{engine.HEALTH_COMPONENT})))

	wl.AddLogic(
		"report health",
		2*time.Second,
		func(em *engine.EntityManager,
			ev *engine.EventBus,
			wl *engine.WorldLogicManager) {

			healthfulEntities := wl.GetEntitiesFromList("healthful")
			fmt.Println("===HEALTH REPORT===")
			if len(healthfulEntities) == 0 {
				fmt.Println("No healthful entities")
			}
			for _, entity := range healthfulEntities {
				health, valid := em.Components.Health.SafeGet(entity)
				if valid {
					fmt.Printf("Entity %d has health %d\n", entity.ID, health)
				}
			}
		})

	Crows := crow.RegisterCrows(&em)
	Beetles := beetle.RegisterBeetles(&em)

	em.RequestSpawn(Crows.SpawnRequest([2]int16{0, 0}))
	em.RequestSpawn(Beetles.SpawnRequest([2]int16{30, 30}))
	em.RequestSpawn(Beetles.SpawnRequest([2]int16{0, 30}))
	em.RequestSpawn(Beetles.SpawnRequest([2]int16{30, 0}))

	for {
		var t0 time.Time
		em.Update()
		if PROFILE_EM_UPDATE {
			fmt.Printf("EntityManager.Update() took %d ms\n",
				time.Since(t0).Nanoseconds()/1e6)
		}
		time.Sleep(16 * time.Millisecond)
	}
}
