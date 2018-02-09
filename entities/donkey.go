package entities

import (
    "math"
    "math/rand"

    "github.com/dt-rush/donkeys-qquest/engine"
    "github.com/dt-rush/donkeys-qquest/engine/components"
    "github.com/dt-rush/donkeys-qquest/constants"
)

func SpawnDonkey (entity_manager *engine.EntityManager,
    active_component *components.ActiveComponent,
    position_component *components.PositionComponent,
    velocity_component *components.VelocityComponent,
    // color_component *components.ColorComponent,
    hitbox_component *components.HitboxComponent,
    sprite_component *components.SpriteComponent,
    logic_component *components.LogicComponent) int {

    donkey_components := []engine.Component{
        engine.Component (active_component),
        engine.Component (position_component),
        engine.Component (velocity_component),
        // engine.Component (color_component),
        engine.Component (hitbox_component),
        engine.Component (sprite_component),
        engine.Component (logic_component),
    }

    donkey_id := entity_manager.SpawnEntity (donkey_components)

    // set component values

    donkey_active := true
    active_component.Set (donkey_id, donkey_active)

    donkey_position := [2]float64 {float64(constants.WINDOW_WIDTH/2) + 40, float64(constants.WINDOW_HEIGHT/2) + 40}
    position_component.Set (donkey_id, donkey_position)

    // donkey_color := uint32 (0xff776622)
    // color_component.Set (donkey_id, donkey_color)

    donkey_hitbox := [2]float64{24, 24}
    hitbox_component.Set (donkey_id, donkey_hitbox)

    donkey_sprite := sprite_component.IndexOf ("donkey.png")
    sprite_component.Set (donkey_id, donkey_sprite)

    // add donkey logic

    donkey_logic_func := DonkeyLogic (donkey_id,
                                        position_component,
                                        velocity_component)
    donkey_logic_unit := components.LogicUnit{"donkey logic",
                                                donkey_logic_func}
    logic_component.Set (donkey_id, donkey_logic_unit)

    // add tag

    entity_manager.TagEntityUnique (donkey_id, "donkey")

    // return id

    return donkey_id
}

func DonkeyLogic (donkey_id int,
    position_component *components.PositionComponent,
    velocity_component *components.VelocityComponent) func (int) {

    // closure state

    // time accumulator(s)
    one_second := engine.CreateTimeAccumulator (1000)
    // donkey_acceleration in a given direction
    donkey_accel := make ([]float64, 2)

    // enclosed function

    return func (dt_ms int) {

        donkey_pos := position_component.Get (donkey_id)
        donkey_vel := velocity_component.Get (donkey_id)

        // every second, set a new random heading (donkeys are finnicky)
        if one_second.Tick (dt_ms) {
            donkey_accel[0] = 2 * (rand.Float64() - 0.5)
            donkey_accel[1] = 2 * (rand.Float64() - 0.5)
        }
        // apply the accel
        donkey_vel[0] += donkey_accel[0]
        donkey_vel[1] += donkey_accel[1]
        // decay the accel
        donkey_accel[0] = 0.9 * donkey_accel[0]
        donkey_accel[1] = 0.9 * donkey_accel[1]
        // add some random circular wobble
        donkey_vel[0] += 1.2 * math.Cos (2 * math.Pi * one_second.Completion())
        donkey_vel[1] += 1.2 * math.Sin (2 * math.Pi * one_second.Completion())

        // donkey experiences acceleration toward center of screen
        center_accel_strength := 2.0
        // x
        if (int32 (donkey_pos[0]) > constants.WINDOW_WIDTH / 2) {donkey_vel[0] -= center_accel_strength}
        if (int32 (donkey_pos[0]) < constants.WINDOW_WIDTH / 2) {donkey_vel[0] += center_accel_strength}
        // y
        if (int32 (donkey_pos[1]) > constants.WINDOW_HEIGHT / 2) {donkey_vel[1] -= center_accel_strength}
        if (int32 (donkey_pos[1]) < constants.WINDOW_HEIGHT / 2) {donkey_vel[1] += center_accel_strength}

        velocity_component.Set (donkey_id, donkey_vel)
    }
}
