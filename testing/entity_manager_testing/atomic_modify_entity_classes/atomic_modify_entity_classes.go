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
const N_CROWS = 1
const N_BEETLES = 3
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
			fmt.Printf("%d healthful entities\n",
				len(healthfulEntities))
			for _, entity := range healthfulEntities {
				fmt.Printf("trying to get health for %v\n", entity)
				health, valid := em.Components.Health.SafeGet(entity)
				if valid {
					fmt.Printf("Entity %d has health %d\n", entity.ID, health)
				}
			}
		})

	em.RegisterEntityClass(&crow.Crows)
	em.RegisterEntityClass(&beetle.Beetles)

	for i := 0; i < N_CROWS; i++ {
		em.RequestSpawn(crow.SpawnRequest([2]int16{0, 0}))
	}
	for i := 0; i < N_BEETLES; i++ {
		em.RequestSpawn(beetle.SpawnRequest([2]int16{30, 30}))
	}

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
