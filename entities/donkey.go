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
    color_component *components.ColorComponent,
    hitbox_component *components.HitboxComponent,
    logic_component *components.LogicComponent) int {

    donkey_components := []engine.Component{
        engine.Component (active_component),
        engine.Component (position_component),
        engine.Component (velocity_component),
        engine.Component (color_component),
        engine.Component (hitbox_component),
        engine.Component (logic_component),
    }

    donkey_id := entity_manager.SpawnEntity (donkey_components)

    // set component values

    donkey_active := true
    active_component.Set (donkey_id, donkey_active)

    donkey_position := [2]float64 {float64(constants.WINDOW_WIDTH/2) + 40, float64(constants.WINDOW_HEIGHT/2) + 40}
    position_component.Set (donkey_id, donkey_position)

    donkey_color := uint32 (0xff776622)
    color_component.Set (donkey_id, donkey_color)

    donkey_hitbox := [2]float64{24, 24}
    hitbox_component.Set (donkey_id, donkey_hitbox)

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
    velocity_component *components.VelocityComponent) func (float64) {

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

        for dt_accum > 1000 {
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

        velocity_component.Set (donkey_id, donkey_vel)
    }
}
