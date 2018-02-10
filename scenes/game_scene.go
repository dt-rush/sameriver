/*
 *
 *
 *
 *
**/



package scenes

import (
    "fmt"
    "time"
    "math/rand"

    "github.com/dt-rush/donkeys-qquest/engine"
    "github.com/dt-rush/donkeys-qquest/engine/utils"
    "github.com/dt-rush/donkeys-qquest/engine/components"
    "github.com/dt-rush/donkeys-qquest/engine/systems"

    "github.com/dt-rush/donkeys-qquest/constants"
    "github.com/dt-rush/donkeys-qquest/entities"
    "github.com/dt-rush/donkeys-qquest/logic"

    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/ttf"
)

type GameScene struct {

    // Scene "abstract class members"

    // whether the scene is running
    running bool
    // used to make destroy() idempotent
    destroyed bool
    // the game
    game *engine.Game

    // ECS declarations

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
    // logic component
    logic_component components.LogicComponent

    // systems

    // screenmessage system
    screenmessage_system systems.ScreenMessageSystem
    // collision system
    collision_system systems.CollisionSystem
    // physics system
    physics_system systems.PhysicsSystem
    // logic system
    logic_system systems.LogicSystem

    // score of player in this scene
    score int
    score_font *ttf.Font
    // score
    score_surface *sdl.Surface
    // texture of the above, for Renderer.Copy() in draw()
    score_texture *sdl.Texture
    // score texture screen width
    score_rect sdl.Rect

    // utilities
    // function profiler
    func_profiler utils.FuncProfiler
    // profiling data
    collision_detection_ms_accum int
    collision_detection_count int
    physics_ms_accum int
    physics_count int
    draw_ms_accum int
    draw_count int
    logic_ms_accum int
    logic_count int
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
    // logic component
    s.logic_component.Init (s.entity_manager.NumberOfEntities(), s.game)
    // init systems

    // TODO: determine how to tune capacity here
    // 4 as a nonsense magic number
    s.game_event_system.Init (4)

    s.screenmessage_system.Init (4)

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
        &s.game_event_system,
        &s.logic_component,
        &s.active_component)



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

    // spawn a player

    s.player_id = entities.SpawnPlayer (
        &s.entity_manager,
        &s.active_component,
        &s.position_component,
        &s.velocity_component,
        &s.color_component,
        &s.hitbox_component)


    // spawn a donkey

    s.donkey_id = entities.SpawnDonkey (
        &s.entity_manager,
        &s.active_component,
        &s.position_component,
        &s.velocity_component,
        // &s.color_component,
        &s.hitbox_component,
        &s.sprite_component,
        &s.logic_component)

    // spawn N_FLAMES

    s.N_FLAMES = 4
    for i := 0; i < s.N_FLAMES; i++ {

        corners := [2]int{i % 2, i / 2}

        utils.DebugPrintf ("spawning flame in corner %d, %d\n", corners[0], corners[1])

        initial_position := [2]float64{
            float64 (int (constants.WINDOW_WIDTH - 50) * corners [0] + 25),
            float64 (int (constants.WINDOW_HEIGHT - 50)  * corners [1] + 25),
        }

        entities.SpawnFlame (
            &s.entity_manager,
            &s.active_component,
            &s.position_component,
            &s.velocity_component,
            // &s.color_component,
            &s.hitbox_component,
            &s.sprite_component,
            &s.logic_component,
            initial_position,
        )
    }
}

func (s *GameScene) Init (game *engine.Game) chan bool {

    init_done_signal_chan := make (chan bool)
    s.game = game

    s.score = 0

    go func () {
        s.destroyed = false

        all_components := []engine.Component{
            &s.active_component,
            &s.position_component,
            &s.velocity_component,
            &s.color_component,
            &s.audio_component,
            &s.sprite_component,
            &s.hitbox_component,
            &s.logic_component,
        }
        // 10 is a magic number, they scale dynamically anyway
        s.entity_manager.Init (10, all_components)
        // set up components and system
        s.setup_ECS()
        // set up the collision
        s.add_collision_logic()
        // spawn some entities
        s.spawn_entities()
        // load the score font
        var err error
        if s.score_font , err = ttf.OpenFont ("./assets/test.ttf", 12); err != nil {
            panic(err)
        }
        // set up the score surface/texture
        s.update_score_texture()

        // just to play a little loading screen fun
        time.Sleep (1 * time.Second)
        init_done_signal_chan <- true
    }()
    return init_done_signal_chan
}

func (s *GameScene) update_score_texture () {
    if s.score_surface != nil {
        s.score_surface.Free()
    }
    if s.score_texture != nil {
        s.score_texture.Destroy()
    }
    // render message ("press space") surface
    score_msg := fmt.Sprintf ("%d", s.score)
    var err error
    s.score_surface, err = s.score_font.RenderUTF8Solid (
        score_msg,
        sdl.Color{255, 255, 255, 255})
    if err != nil {
        panic (err)
    }
    // create the texture
    s.score_texture, err = s.game.Renderer.CreateTextureFromSurface (s.score_surface)
    if err != nil {
        panic (err)
    }
    // set the width of the texture on screen
    s.score_rect = sdl.Rect{
        10,
        10,
        int32 (len (score_msg) * 20),
        20}
}



func (s *GameScene) Stop () {
    utils.DebugPrintf ("\n\n\n======== ADVENTURE OVER ========\n")
    utils.DebugPrintf ("================================\n\n\n")
    utils.DebugPrintf ("collision_detection_ms_avg = %.3f ms\n",
        float64 (s.collision_detection_ms_accum) /
        float64 (s.collision_detection_count))
    utils.DebugPrintf ("physics_ms_avg = %.3f ms\n",
        float64 (s.physics_ms_accum) /
        float64 (s.physics_count))
    utils.DebugPrintf ("draw_ms_avg = %.3f ms\n",
        float64 (s.draw_ms_accum) /
        float64 (s.draw_count))
    utils.DebugPrintf ("logic_ms_avg = %.3f ms\n",
        float64 (s.logic_ms_accum) /
        float64 (s.logic_count))
    // set this scene not running
    s.running = false
    
    // actually ends the game
    // s.game.NextSceneChan <- nil
}

func (s *GameScene) IsRunning () bool {
    return s.running
}



func (s *GameScene) Update (dt_ms int) {

    // TODO: form an array of loaded systems and iterate them all

    s.physics_count++
    s.physics_ms_accum += s.func_profiler.Time (
        func (dt_ms int) (func ()) {
            return func () {
                s.physics_system.Update (dt_ms)
            }
        }(dt_ms))

    s.collision_detection_count++
    s.collision_detection_ms_accum += s.func_profiler.Time (
        func (dt_ms int) (func ()) {
            return func () {
                s.collision_system.Update (dt_ms)
            }
        }(dt_ms))

    s.logic_count++
    s.logic_ms_accum += s.func_profiler.Time (
        func (dt_ms int) (func()) {
            return func () {
                s.logic_system.Update (dt_ms)
            }
        }(dt_ms))
}






func (s *GameScene) Draw (window *sdl.Window, renderer *sdl.Renderer) {

    s.draw_count++
    s.draw_ms_accum += s.func_profiler.Time (func () {

        // draw the score
        renderer.Copy (
            s.score_texture,
            nil,
            &s.score_rect)

        // TODO refactor to go through only entities registered
        // with a draw system to avoid checking EntityHasComponent 
        for _, i := range s.entity_manager.Entities() {

            if ! s.active_component.Get (i) {
                // don't draw inactive entities
                continue
            }

            pos := s.position_component.Get (i)
            // ss_pos == "screen-space pos"
            // note that we're not checking first to see if the entity has a hitbox
            // draw the box such that its center is where the position of the entity is
            box := s.hitbox_component.Get (i)
            ss_pos := make ([]int32, 2)
            ss_pos [0] = int32 (pos [0] - (box [0] / 2))
            ss_pos [1] = constants.WINDOW_HEIGHT - (int32 (pos [1]) + int32 (box [1] / 2))
            entity_screen_rect := sdl.Rect{ss_pos [0],
                ss_pos [1],
                int32 (box [0]),
                int32 (box [1])}


            if s.entity_manager.EntityHasComponent (i, &s.color_component) {
                color := s.color_component.Get (i)
                // extracting color components from
                // uint32 ARGB to uint8 RGBA params
                renderer.SetDrawColor (
                    uint8 ((color & 0x00ff0000) >> 16),
                    uint8 ((color & 0x0000ff00) >> 8),
                    uint8 ((color & 0x000000ff) >> 0),
                    uint8 ((color & 0xff000000) >> 24))
                renderer.FillRect (&entity_screen_rect)
            }


            if s.entity_manager.EntityHasComponent (i, &s.sprite_component) {
                // implement component data in sprite_component for these as well?
                var angle float64 = 0
                var center_p *sdl.Point = nil
                renderer.CopyEx (s.sprite_component.Get (i),
                    nil,
                    &entity_screen_rect,
                    angle,
                    center_p,
                    s.sprite_component.GetFlip (i))
            }


            // paint a little white rect where the corner of the box is
            // renderer.SetDrawColor (255, 255, 255, 255)
            // small_rect := sdl.Rect{ss_pos[0], ss_pos[1], 4, 4}
            // renderer.FillRect (&small_rect)

        }
    })

}




func (s *GameScene) HandleKeyboardState (keyboard_state []uint8) {

    k := keyboard_state

    // get player v0
    player_v := s.velocity_component.Get (s.player_id)
    // get player v1
    vx := 300 * float64 (
        int8 (k [sdl.SCANCODE_D]) -
            int8 (k [sdl.SCANCODE_A]))
    vy := 300 * float64 (
        int8 (k [sdl.SCANCODE_W]) -
            int8 (k [sdl.SCANCODE_S]))
    // shift v0 to v1
    player_v[0] = vx
    player_v[1] = vy
    // set v in the map (for some reason doesn't modify if we just modify
    // player_v, you have to actually put(), in effect)
    s.velocity_component.Set (s.player_id, player_v)

}

func (s *GameScene) HandleKeyboardEvent (keyboard_event *sdl.KeyboardEvent) {
    // null implementation
}

func (s *GameScene) Destroy() {
    // destroy resources claimed
    if ! s.destroyed {
        // using sdl.Do to avoid an issue described in comments
        // in menuscene.go
        sdl.Do (func () {
            s.score_surface.Free()
            s.score_texture.Destroy()
            s.score_font.Close()
            s.sprite_component.Destroy()
        })

        s.destroyed = true
    }
}

func (s *GameScene) SceneLogic () {

    donkey_caught_react := func () {

        donkey_caught_chan := s.game_event_system.Subscribe (constants.GAME_EVENT_DONKEY_CAUGHT)
        for _ = range (donkey_caught_chan) {
            utils.DebugPrintln ("\tYOU CAUGHT A DONKEY")

            s.score += 1
            s.update_score_texture()

            // TODO: expand to actual inventory system
            PRINT_DONKEY_INVENTORY := true
            if PRINT_DONKEY_INVENTORY {
                inventory := []string{"1 x donkey fur", "2 x donkey ears", "3 x donkey whiskers", "4 x donkey meats"}
                for _, item := range (inventory) {
                    utils.DebugPrintf ("\t\t* %s\n", item)
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
                donkey_pos [0] = rand.Float64() * 
                    float64 (constants.WINDOW_WIDTH - 20) + 20
                donkey_pos [1] = rand.Float64() * 
                    float64 (constants.WINDOW_HEIGHT - 20) + 20
                s.active_component.Set (donkey_id, true)
            }()
        }
    }

    flame_hit_player_react := func () {

        flame_hit_player_chan := s.game_event_system.Subscribe (
            constants.GAME_EVENT_FLAME_HIT_PLAYER)

        for _ = range (flame_hit_player_chan) {
            utils.DebugPrintln ("\tYOU DIED BY FALLING IN A FIRE")
            game_over_scene := GameOverScene{}
            s.game.NextScene = &game_over_scene
            s.Stop()
            // stop listening for these events by breaking 
            break
        }
    }

    go donkey_caught_react()
    go flame_hit_player_react()

}

func (s *GameScene) Run () {

    // any scene-specific routines can be spawned in here
    utils.DebugPrintf ("\n\n\n=====================================\n")
    utils.DebugPrintf ("======== ADVENTURE BEGINNING ========\n\n\n\n")

    // spawn scene logic goroutine
    go s.SceneLogic()

    s.running = true

    utils.DebugPrintln ("GameScene.Run() completed.")

}



func (s *GameScene) Name () string {
    return "game scene"
}
