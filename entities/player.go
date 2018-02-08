package entities

import (
    "github.com/dt-rush/donkeys-qquest/engine"
    "github.com/dt-rush/donkeys-qquest/constants"
)

func SpawnPlayer (entity_manager *engine.EntityManager,
                    active_component engine.Component,
                    position_component engine.Component,
                    velocity_component engine.Component,
                    color_component engine.Component,
                    hitbox_component engine.Component) int {

    // register with entity manager

    player_components := []engine.Component{active_component,
        position_component,
        velocity_component,
        color_component,
        hitbox_component}

    player_id := entity_manager.SpawnEntity (player_components)

    // set component values

    player_active := true
    active_component.Set (player_id, player_active)

    player_position := [2]float64 {float64(constants.WINDOW_WIDTH/2), float64(constants.WINDOW_HEIGHT/2)}
    position_component.Set (player_id, player_position)

    player_color := uint32 (0xff00AACC)
    color_component.Set (player_id, player_color)

    player_hitbox := [2]float64{20, 20}
    hitbox_component.Set (player_id, player_hitbox)

    // add tag

    entity_manager.TagEntityUnique (player_id, "player")

    // return id

    return player_id
}
