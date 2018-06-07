Rationale
===

In order to prevent deadlocks, we need to introduce lock ordering and hierarchy,
in which:

1. a code path can only acquire locks on lower levels than those it already has

and

2. locks on a given level must be acquired "at once", meaning that a list is 
defined of locks on the level to acquire, and we traverse it *in order*, waiting
for any locks which are held when we reach them

If these two rules are followed, deadlock cannot occur, only waits.

Inventory of existing locks
===

At present, the following are the locks which exist

TODO: this is out of date slightly after latest changes, revisit

	./entity_manager.go:	spawnMutex sync.Mutex
	./event_channel.go:	channelSendLock sync.Mutex
	./signal_send_rate_limiter.go:	mutex sync.Mutex
	./entity_logic_table.go:	mutex            sync.RWMutex
	./updated_entity_list.go:	Mutex sync.RWMutex
	./world_logic.go:	mutex  sync.RWMutex
	./logic_table.go:	mutex      sync.RWMutex
	./active_entity_list_collection.go:	mutex    sync.RWMutex
	./tag_table.go:	mutex           sync.RWMutex
	./entity_class_table.go:	mutex   sync.RWMutex
	./event_bus.go:	mutex sync.RWMutex
	./entity_query_watcher_list.go:	mutex sync.RWMutex
	./entity_table.go:	IDMutex sync.RWMutex
	./entity_table.go:	locks [MAX_ENTITIES]atomic.Uint32
	./rate_limiter.go:	mutex sync.RWMutex
	./rate_limiter.go:	mutex sync.RWMutex

Analysis of each lock
===

We will list each lock and the functions in which it becomes locked, also
listing the exposed functions which can call that function


#### `EntityManager.spawnMutex`

functions:

* `EntityManager.processSpawnChannel()` via `EntityManager.Update()`
* `EntityManager.despawnAll()` via `EntityManager.DespawnAll()`


#### `EventChannel.channelSendLock`

functions:

* `EventChannel.Send()`
* `EventChannel.DrainChannel()`


#### `SignalSendRateLimiter.mutex`

functions:

* `SignalSendRateLimiter.Do()`


#### `EntityLogicTable.mutex` (RW)

functions:

* `EntityLogicTable.setLogic() via EntityManager.Spawn()`
* `EntityLogicTable.getLogic() via EntityManager.setActiveState() via EntityManager.Activate(), EntityManager.Deactivate()` (R)
* `EntityLogicTable.deleteLogic() via EntityManager.despawnInternal() via EntityManager.DespawnAll(), EntityManager.Despawn()`


#### `UpdatedEntityList.Mutex` (RW)

functions:
* `UpdatedEntityList.Length()` (R)
* `UpdatedEntityList.FirstEntity()` (R)
* `UpdatedEntityList.RandomEntity()` (R)
* `UpdatedEntityList.actOnSignal() via UpdatedEntitylist.start() via ActiveEntityListCollection.GetUpdatedEntityList()`
 

#### ``

functions:

* ` `
* ` `


#### 

functions:

* ` `
* ` `


#### 

functions:

* ` `
* ` `


#### 

functions:

* ` `
* ` `


#### 

functions:

* ` `
* ` `


#### 

functions:

* ` `
* ` `


#### 

functions:

* ` `
* ` `


#### 

functions:

* ` `
* ` `


#### 

functions:

* ` `
* ` `


#### 

functions:

* ` `
* ` `


#### 

functions:

* ` `
* ` `


#### 

functions:

* ` `
* ` `


#### 

functions:

* ` `
* ` `


#### 

functions:

* ` `
* ` `


#### 

functions:

* ` `
* ` `


#### 

functions:

* ` `
* ` `

