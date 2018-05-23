package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dt-rush/donkeys-qquest/engine"
)

const N_CROWS = 16

func spawnCrows(em *engine.EntityManager) {
	for i := 0; i < N_CROWS; i++ {
		spawnCrow(em)
	}
}

func spawnCrow(em *engine.EntityManager) {
	logic := engine.NewLogicUnit(
		func(crowID uint16,
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
					fmt.Printf("Crow %d says, CAW!\n", crowID)
					time.Sleep(
						time.Duration(rand.Intn(3000)) * time.Millisecond)
				}
			}
		}, "crow logic")
	componentSet := engine.ComponentSet{}
	componentSet.Logic = &logic
	tags := []string{"bird", "crow"}
	spawnRequest := engine.EntitySpawnRequest{
		componentSet,
		tags}
	em.RequestSpawn(spawnRequest)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	t0 := time.Now()
	em := engine.EntityManager{}
	em.Init()
	fmt.Printf("EntityManager.Init() took %d ms\n",
		time.Since(t0).Nanoseconds()/1e6)
	fmt.Println("Starting entity manager testing")

	t0 = time.Now()
	spawnCrows(&em)
	fmt.Printf("spawnCrows() took %d ms\n",
		time.Since(t0).Nanoseconds()/1e6)

	for {
		t0 = time.Now()
		em.Update()
		fmt.Printf("EntityManager.Update() took %d ms\n",
			time.Since(t0).Nanoseconds()/1e6)
		time.Sleep(16 * time.Millisecond)
	}
}
