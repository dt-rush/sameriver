/**
  * 
  * 
  * 
  * 
**/



package scenes

import (
	"fmt"
	"time"
	"math"
	"math/rand"

	"github.com/dt-rush/donkeys-qquest/engine"
	"github.com/dt-rush/donkeys-qquest/components"
	"github.com/dt-rush/donkeys-qquest/systems"
	"github.com/dt-rush/donkeys-qquest/utils"
	"github.com/dt-rush/donkeys-qquest/constants"
	"github.com/dt-rush/donkeys-qquest/logic"
	
	"github.com/veandco/go-sdl2/sdl"
//	"github.com/veandco/go-sdl2/img"
)

type GameScene struct {

	// TODO separate
	// Scene "abstract class members"

	// whether the scene is running
	running bool
	// used to make destroy() idempotent
	destroyed bool
	// the game
	game *engine.Game


	

	// TODO preserve
	// data specific to this scene
	// (all the below)

	// ECS (TM) declarations
	
	// entity manager
	entity_manager engine.EntityManager
	// gameevent system
	game_event_system engine.GameEventSystem

	// special keys into the entity array
	player_id int
	donkey_id int
	N_FLAMES int

	// components
	
	// active component
	active_component components.ActiveComponent
	// sprite component
	sprite_component components.SpriteComponent
	// color component
	color_component components.ColorComponent
	// audio component
	audio_component components.AudioComponent
	// position component
	position_component components.PositionComponent
	// velocity component
	velocity_component components.VelocityComponent
	// hitbox component
	hitbox_component components.HitboxComponent

	// actual (TM) systems
	
	// screenmessage system
	screenmessage_system systems.ScreenmessageSystem
	// collision system
	collision_system systems.CollisionSystem
	// physics system
	physics_system systems.PhysicsSystem
	// logic system
	logic_system systems.LogicSystem


	// utilities

	// (blocking) function profiler
	func_profiler utils.FuncProfiler

	// profiling data
	collision_detection_ms_accum float64
	collision_detection_count int
	physics_ms_accum float64
	physics_count int
	draw_ms_accum float64
	draw_count int
	logic_ms_accum float64
	logic_count int
}



func (s GameScene) CreateScene () engine.Scene {
	return engine.Scene (&s)
}



func (s *GameScene) setup_ECS() {

	// ECS (TM)

	// init components
	
	// active component
	s.active_component.Init (s.entity_manager.NumberOfEntities(), s.game)
	// sprite component
	s.sprite_component.Init (128, s.game)
	// color component
	s.color_component.Init (s.entity_manager.NumberOfEntities(), s.game)
	// audio component
	s.audio_component.Init (s.entity_manager.NumberOfEntities(), s.game)
	// position component
	s.position_component.Init (s.entity_manager.NumberOfEntities(), s.game)
	// velocity component
	s.velocity_component.Init (s.entity_manager.NumberOfEntities(), s.game)
	// hitbox component
	s.hitbox_component.Init (s.entity_manager.NumberOfEntities(), s.game)

	// init systems
	
	// TODO: determine how to tune capacity here
	// 4 as a nonsense magic number
	s.game_event_system.Init (4)


	// TODO, turn into a system (TM) which can spawn entities
	// YES< FCUK IT, SCREENMESSAGES ARE ENTITIES,
	// SO WILL BE MENUS 
	// YOU CAN STILL BACK THEM WITH SYSTEMS AND LOGIC FUNCS
	
//	s.screenmessage_system.init (4)


	// init (TM) systems
	
	s.collision_system.Init (&s.entity_manager,
		&s.active_component,
		&s.position_component,
		&s.hitbox_component,
		&s.game_event_system)
	
	s.physics_system.Init (&s.entity_manager,
		&s.active_component,
		&s.position_component,
		&s.velocity_component)

	s.logic_system.Init (&s.entity_manager,
		&s.game_event_system)


	
	// utilities
	s.func_profiler.Init (4)
	
}

func (s *GameScene) add_collision_logic () {

	// load and add collision logic from logic package
	// exported variable CollisionLogicCollection

	for name, l := range logic.CollisionLogicCollection {
		id := s.collision_system.AddCollisionLogic (name, l)
		s.collision_system.SetCollisionLogicActiveState (id, true)
	}

}



func (s *GameScene) spawn_entities() {
		// spawn some entities

	// spawn entities by entity_manager.spawn_entity ([]Component)).
	// This allocates component data for each entity

	// spawn a player
	
	player_components := []engine.Component{engine.Component(&s.active_component),
		engine.Component(&s.position_component),
		engine.Component(&s.velocity_component),
		engine.Component(&s.hitbox_component)}
	// engine.Component(&s.logic_component),
	// engine.Component(&s.color_component),
	// engine.Component(&s.sprite_component),
	
	s.player_id = s.entity_manager.SpawnEntity (player_components)
	// TODO refactor into a function on entity_manager accepting
	// a diff of component-values (map?)
	player_active := true
	s.active_component.Set (s.player_id, player_active)
	player_position := [2]float64 {float64(constants.WINDOW_WIDTH/2), float64(constants.WINDOW_HEIGHT/2)}
	s.position_component.Set (s.player_id, player_position)
	player_color := uint32 (0xff00AACC)
	s.color_component.Set (s.player_id, player_color)
	player_hitbox := [2]float64{20, 20}
	s.hitbox_component.Set (s.player_id, player_hitbox)
	// add tag
	s.entity_manager.TagEntityUnique (s.player_id, "player")

	


	// spawn a donkey
	
	donkey_components := []engine.Component{engine.Component(&s.active_component),
		engine.Component(&s.position_component),
		engine.Component(&s.velocity_component),
		engine.Component(&s.hitbox_component)}
	// engine.Component(&s.logic_component),
	// engine.Component(&s.color_component),
	// engine.Component(&s.sprite_component)
	
	s.donkey_id = s.entity_manager.SpawnEntity (donkey_components)
	// init the entity's component values
	// TODO refactor into a function on entity_manager accepting
	// a diff of component-values (map?)
	donkey_active := true
	s.active_component.Set (s.donkey_id, donkey_active)
	donkey_position := [2]float64 {float64(constants.WINDOW_WIDTH/2) + 40, float64(constants.WINDOW_HEIGHT/2) + 40}
	s.position_component.Set (s.donkey_id, donkey_position)
	donkey_color := uint32 (0xff776622)
	s.color_component.Set (s.donkey_id, donkey_color)
	donkey_hitbox := [2]float64{24, 24}
	s.hitbox_component.Set (s.donkey_id, donkey_hitbox)
	// add donkey logic
	s.logic_system.RunLogic ((func (donkey_id int,
		position_component *components.PositionComponent,
		velocity_component *components.VelocityComponent) (func (float64)) {
			// closure state
			// time accumulator
			var dt_accum float64 = 0
			// acceleration in a given direction
			accel := make ([]float64, 2)
			// enclosed function
			return func (dt float64) {
				dt_accum += dt
				donkey_pos := position_component.Get (donkey_id)
				donkey_vel := velocity_component.Get (donkey_id)

				if dt_accum > 1000 {
					dt_accum -= 1000
					accel[0] = rand.Float64() - 0.5
					accel[1] = rand.Float64() - 0.5
				}
				accel[0] = 0.9 * accel[0]
				accel[1] = 0.9 * accel[1]

				donkey_vel[0] += 0.8 * math.Cos (2 * math.Pi * dt_accum / 1000.0)
				donkey_vel[1] += 0.8 * math.Sin (2 * math.Pi * dt_accum / 1000.0)
				
				donkey_vel[0] += accel[0]
				donkey_vel[1] += accel[1]
				
				// donkey experiences acceleration toward center of screen
				// x
				if (int32 (donkey_pos[0]) > constants.WINDOW_WIDTH / 2) {donkey_vel[0] -= .1}
				if (int32 (donkey_pos[0]) < constants.WINDOW_WIDTH / 2) {donkey_vel[0] += .1}
				// y
				if (int32 (donkey_pos[1]) > constants.WINDOW_HEIGHT / 2) {donkey_vel[1] -= .1}
				if (int32 (donkey_pos[1]) < constants.WINDOW_HEIGHT / 2) {donkey_vel[1] += .1}

				// TODO replace statements like this with defer statements or at least
				// find a way to directly update map values, to avoid this weird
				// modify and replace pattern
				velocity_component.Set (donkey_id, donkey_vel)
				
			}
		})(s.donkey_id, &s.position_component, &s.velocity_component))
	// add tag
	s.entity_manager.TagEntityUnique (s.donkey_id, "donkey")

	

	// spawn N_FLAMES
	
	s.N_FLAMES = 3
	
	for i := 0; i < s.N_FLAMES; i++ {
		
		flame_components := []engine.Component{engine.Component(&s.active_component),
			engine.Component(&s.position_component),
			engine.Component(&s.velocity_component),
			engine.Component(&s.hitbox_component),
			// engine.Component(&s.logic_component),
			// engine.Component(&s.color_component),
			engine.Component(&s.sprite_component)}
		
		// init the entity's component values
		// TODO refactor into a function on entity_manager accepting
		// a diff of component-values (map?)
		flame_id := s.entity_manager.SpawnEntity (flame_components)
		flame_active := true
		s.active_component.Set (flame_id, flame_active)
		flame_position := [2]float64 {rand.Float64() * float64 (constants.WINDOW_WIDTH - 20) + 20,
			rand.Float64() * float64 (constants.WINDOW_HEIGHT - 20) + 20}
		s.position_component.Set (flame_id, flame_position)
		flame_color := uint32 (0xffccaa33)
		s.color_component.Set (flame_id, flame_color)
		flame_hitbox := [2]float64{50, 50}
		s.hitbox_component.Set (flame_id, flame_hitbox)
		flame_sprite := s.sprite_component.IndexOf ("flame.png")
		s.sprite_component.Set (flame_id, flame_sprite)
		// add flame logic
		s.logic_system.RunLogic ((func (flame_id int,
			position_component *components.PositionComponent,
			velocity_component *components.VelocityComponent,
			sprite_component *components.SpriteComponent) (func (float64)) {
			// closure state
			// time accumulator
			var dt_accum float64 = 0
			// time accumulator for sprite changing
			var sprite_dt_accum float64 = 0
			// current heading
			var heading []float64 = make ([]float64, 2)
			
			// enclosed function
			return func (dt float64) {
				dt_accum += dt
				sprite_dt_accum += dt
				
				flame_vel := velocity_component.Get (flame_id)
				
				for dt_accum > 200 {
					// if we hit 1000 ms, reset the counter
					dt_accum -= 200
					// ... and trigger change of heading
					heading [0] = 100 * (rand.Float64()*2 - 1)
					heading [1] = 100 * (rand.Float64()*2 - 1)
					flame_vel [0] = heading [0]
					flame_vel [1] = heading [1]
					
				}

				for sprite_dt_accum > 500 {
					sprite_dt_accum -= 500
					// toggle flip horizontal
					current_flip := s.sprite_component.GetFlip (flame_id)
					var new_flip sdl.RendererFlip
					if current_flip == sdl.FLIP_NONE {
						new_flip = sdl.FLIP_HORIZONTAL
					} else {
						new_flip = sdl.FLIP_NONE
					}
					s.sprite_component.SetFlip (flame_id, new_flip)
				}

				// TODO replace statements like this with defer statements or at least
				// find a way to directly update map values, to avoid this weird
				// modify and replace pattern
				s.velocity_component.Set (flame_id, flame_vel)
			}
		})(flame_id, &s.position_component, &s.velocity_component, &s.sprite_component))
		// apply tags to flame
		s.entity_manager.TagEntity (flame_id, "flame")
	}
}






func (s *GameScene) Init (game *engine.Game) chan bool {
	
	init_done_sig_chan := make (chan bool)
	s.game = game

	go func () {
		s.destroyed = false
		
		
		all_components := []engine.Component{engine.Component (&s.active_component),
			engine.Component (&s.sprite_component),
			engine.Component (&s.color_component),
			engine.Component (&s.audio_component),
			engine.Component (&s.position_component),
			engine.Component (&s.velocity_component),
			engine.Component (&s.hitbox_component)}
		// 10 is a magic number, they scale dynamically anyway
		s.entity_manager.Init (10, all_components)
		// set up components and system
		s.setup_ECS()
		// set up the collision
		s.add_collision_logic()
		// spawn some entities
		s.spawn_entities()
		
		// just to play a little loading screen fun
		time.Sleep (1 * time.Second)
		init_done_sig_chan <- true
	}()
	return init_done_sig_chan
}

func (s *GameScene) Stop () {
	fmt.Printf ("\n\n\n======== ADVENTURE OVER ========\n")
	fmt.Printf ("================================\n\n\n")
	fmt.Printf ("collision_detection_ms_avg = %.3f ms\n", s.collision_detection_ms_accum / float64 (s.collision_detection_count))
	fmt.Printf ("physics_ms_avg = %.3f ms\n", s.physics_ms_accum / float64 (s.physics_count))
	fmt.Printf ("draw_ms_avg = %.3f ms\n", s.draw_ms_accum / float64 (s.draw_count))
	fmt.Printf ("logic_ms_avg = %.3f ms\n", s.logic_ms_accum / float64 (s.logic_count))
	// set this scene not running
	s.running = false
	// actually ends the game
	s.game.NextSceneChan() <- nil
}

func (s *GameScene) IsRunning () bool {
	return s.running
}



func (s *GameScene) Update (dt_ms float64) {

	// TODO: form an array of loaded systems and iterate them all

	// TODO: finish refactoring s.physics_system_update (dt_ms)
	// to become s.physics_system.update (dt_ms)
	s.physics_count++
	s.physics_ms_accum += s.func_profiler.Time (
		func (dt_ms float64) (func ()) {
			return func () {
				s.physics_system.Update (dt_ms)
			}
		}(dt_ms)) 

	// TODO make a better profiler (see func_profiler.go)
	s.collision_detection_count++
	s.collision_detection_ms_accum += s.func_profiler.Time (
		func (dt_ms float64) (func ()) {
			return func () {
				s.collision_system.Update (dt_ms)
			}
		}(dt_ms)) 



	// update logic system
	s.logic_count++
	s.logic_ms_accum += s.func_profiler.Time (
		func (dt_ms float64) (func()) {
			return func () {
				s.logic_system.Update (dt_ms)
			}
		}(dt_ms))
}






func (s *GameScene) Draw (window *sdl.Window, renderer *sdl.Renderer) {

	s.draw_count++
	s.draw_ms_accum += s.func_profiler.Time (func () {
		
		renderer.SetDrawColor (0, 0, 0, 255)
		renderer.FillRect (&sdl.Rect{0, 0, int32 (constants.WINDOW_WIDTH), int32 (constants.WINDOW_HEIGHT)})

		// TODO refactor to go through only entities registered with a draw system
		// to avoid this index checking
		// draw each entity
		for _, i := range s.entity_manager.Entities() {

			if ! s.active_component.Get (i) {
				// don't draw inactive entities
				continue
			}


			// TODO refactor .has() to be recorded by a component using
			// a map[int]bitarray backing
			// detecting if this entity has a sprite
			has_sprite := s.entity_manager.EntityHasComponent (i, engine.Component (&s.sprite_component))
			
			pos := s.position_component.Get (i)
			// ss_pos == "screen-space pos"
			// note that we're not checking first to see if the entity has a hitbox
			// draw the box such that its center is where the position of the entity is
			box := s.hitbox_component.Get (i)
			ss_pos := make ([]int32, 2)
			ss_pos [0] = int32 (pos [0] - (box [0] / 2))
			ss_pos [1] = constants.WINDOW_HEIGHT - (int32 (pos [1]) + int32 (box [1] / 2))
			screen_rect := sdl.Rect{ss_pos [0],
				ss_pos [1],
				int32 (box [0]),
				int32 (box [1])}

			


			color := s.color_component.Get (i)
			// extracting color components from uint32 ARGB to uint8 RGBA params
			renderer.SetDrawColor (uint8 ((color & 0x00ff0000) >> 16),
				uint8 ((color & 0x0000ff00) >> 8),
				uint8 ((color & 0x000000ff) >> 0),
				uint8 ((color & 0xff000000) >> 24))
			renderer.FillRect (&screen_rect)

			

			if has_sprite {
				// implement component data in sprite_component for these as well?
				var angle float64 = 0
				var center_p *sdl.Point = nil
				renderer.CopyEx (s.sprite_component.Get (i),
					nil,
					&screen_rect,
					angle,
					center_p,
					s.sprite_component.GetFlip (i))
			}
			

			// paint a little white rect where the corner of the box is
			renderer.SetDrawColor (255, 255, 255, 255)
			small_rect := sdl.Rect{ss_pos[0], ss_pos[1], 4, 4}
			renderer.FillRect (&small_rect)

		}
	})

}




func (s *GameScene) HandleKeyboardState (keyboard_state []uint8) {

	k := keyboard_state

	// get player v0
	player_v := s.velocity_component.Get (s.player_id)
	// get player v1
	vx := 200 * float64 (
		int8 (k [sdl.SCANCODE_D]) -
			int8 (k [sdl.SCANCODE_A]))
	vy := 200 * float64 (
		int8 (k [sdl.SCANCODE_W]) -
			int8 (k [sdl.SCANCODE_S]))
	// shift v0 to v1
	player_v[0] = vx
	player_v[1] = vy
	// set v in the map (for some reason doesn't modify if we just modify
	// player_v, you have to actually put(), in effect)
	s.velocity_component.Set (s.player_id, player_v)

}

func (s *GameScene) Destroy() {
	// destroy resources claimed
	if ! s.destroyed {
		// using sdl.Do to avoid an issue described in comments
		// in menuscene.go
		sdl.Do (func () {
			// TODO: does this actually matter? currently
			// main() exits when gamescene ends, and the only way
			// this scene ends is to return nil, so hey...

			// TODO: philosophical debate
			// it's unclear whether gamescene will ever stop running
			// or just push other scenes on a stack (currently the
			// main() and Game logic don't support this)
			
			// delete any resources held
		})
		
		s.destroyed = true
	}
}

func (s *GameScene) Run () {
	
	// any scene-specific routines can be spawned in here
	fmt.Printf ("\n\n\n=====================================\n")
	fmt.Printf ("======== ADVENTURE BEGINNING ========\n\n\n\n")



	// spawn scene logic goroutine
	scene_logic := func () {

		donkey_caught_react := func () {
			donkey_caught_chan := s.game_event_system.Subscribe (constants.GAME_EVENT_DONKEY_CAUGHT)
			for _ = range (donkey_caught_chan) {
				fmt.Println ("\tYOU CAUGHT A DONKEY")
				// TODO: expand to actual inventory system
				PRINT_DONKEY_INVENTORY := false
				if PRINT_DONKEY_INVENTORY {
					inventory := []string{"1 x donkey fur", "2 x donkey ears", "3 x donkey whiskers", "4 x donkey meats"}
					for _, item := range (inventory) {
						fmt.Printf ("\t\t* %s\n", item)
					}
				}
				// set donkey to respawn
				// TODO: also simultaneously set invisible?
				// TODO (possibly): separate "visible" component from "active"?
				donkey_id := s.entity_manager.GetTagEntityUnique ("donkey")
				s.active_component.Set (donkey_id, false)

				// sleep 5 seconds before respawning the donkey
				go func() {
					time.Sleep (time.Second * 5) // blocking
					donkey_pos := s.position_component.Get (donkey_id)
					donkey_pos [0] = rand.Float64() * float64 (constants.WINDOW_WIDTH - 20) + 20
					donkey_pos [1] = rand.Float64() * float64 (constants.WINDOW_HEIGHT - 20) + 20
					s.active_component.Set (donkey_id, true)
				}()
			}
		}

		flame_hit_player_react := func () {
					
			flame_hit_player_chan := s.game_event_system.Subscribe (constants.GAME_EVENT_FLAME_HIT_PLAYER)

			for _ = range (flame_hit_player_chan) {
				fmt.Println ("\tYOU DIED BY FALLING IN A FIRE")
				s.Stop()
				// s.game.NextSceneChan() <- nil
			}
		}

		go donkey_caught_react()
		go flame_hit_player_react()

	}

	go scene_logic()

	

	s.running = true

}



func (s *GameScene) Name () string {
	return "game scene"
}
