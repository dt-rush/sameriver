package logic

import (
    "github.com/dt-rush/donkeys-qquest/engine"
    "github.com/dt-rush/donkeys-qquest/engine/systems"
    "github.com/dt-rush/donkeys-qquest/constants"
)


// check donkey-player collision

var CollisionLogicCollection = map [string]systems.CollisionLogic{

    "player-donkey": systems.CollisionLogic{
        // NOTE: we have to check whether i = player and j = donkey or
        // i = donkey and j = player, because we don't know
        // who will be i or j in the "handshake" as ID's are added to a bag of ID's which
        // may only come out in a given order by coincidence assuring that, for example,
        // the player were always i and the donkey j, never reaching the donkey first
        // via i to compare collisions with a player on j

        Selector: func (i int,
            j int,
            em *engine.EntityManager) bool {

                player_id := em.GetTagEntityUnique ("player")
                donkey_id := em.GetTagEntityUnique ("donkey")

                return ((i == player_id && j == donkey_id) ||
                    (i == donkey_id && j == player_id))
            },

        EventGenerator: func (i int,
            j int,
            em *engine.EntityManager) engine.GameEvent {

                return constants.GAME_EVENT_DONKEY_CAUGHT
            },
    },

    "player-flame": systems.CollisionLogic{

        Selector: func (i int,
            j int,
            em *engine.EntityManager) bool {
                player_id := em.GetTagEntityUnique ("player")
                i_is_flame := em.EntityHasTag (i, "flame")
                j_is_flame := em.EntityHasTag (j, "flame")

                return ((i == player_id && j_is_flame) ||
                    (j == player_id && i_is_flame))
            },

        EventGenerator: func (i int,
            j int,
            em *engine.EntityManager) engine.GameEvent {

                return constants.GAME_EVENT_FLAME_HIT_PLAYER
            },
    },
}



