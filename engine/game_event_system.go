/**
  *
  *
  *
  *
**/



// Game event system is a bunch of channels used for systems to interact outside of reading
// and writing to components, using pub-sub pattern and direct, instantaneous sends

// This class creates and manages the channels

package engine



import (
    "sync"
)



type GameEventSystem struct {
    channels [](chan GameEvent)
    subscribers map[GameEvent]([](chan GameEvent))

    subscribe_mutex *sync.Mutex
}

func (s *GameEventSystem) Init (capacity int) {
    // be ready for about n^2 channels for n subscribers
    s.channels = make ([](chan GameEvent), capacity * capacity)
    s.subscribers = make (map[GameEvent]([](chan GameEvent)), capacity)

    s.subscribe_mutex = &sync.Mutex{}
}

func (s *GameEventSystem) Subscribe (e GameEvent) chan GameEvent {
    s.subscribe_mutex.Lock()
    // create subscribers array if DNE
    _, ok := s.subscribers [e]
    if ! ok {
        s.subscribers [e] = make ([](chan GameEvent), 0)
    }
    ch := make (chan GameEvent)
    s.subscribers [e] = append (s.subscribers [e], ch)
    s.channels = append (s.channels, ch)
    s.subscribe_mutex.Unlock()
    return ch
}

func (s *GameEventSystem) Publish (e GameEvent) {
    if DEBUG_GAME_EVENTS {
        Logger.Println (e)
    }
    // send e to all streams listening for GameEvent
    for _, ch := range s.subscribers [e] {
        go func (ch chan GameEvent, e GameEvent) {
            ch <- e
        }(ch, e)
    }
}

func (s *GameEventSystem) NumberOfActiveChannels () int {
    return len (s.channels)
}

func (s *GameEventSystem) GetChannel () chan GameEvent {
    new_channel := make (chan GameEvent)
    s.channels = append (s.channels, new_channel)
    return new_channel
}




