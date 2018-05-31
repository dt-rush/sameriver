# Description 

This script was made to profile the performance of sorting entities into a 
spatial hash (useful for restricting the scope of calculations needed in 
physics, drawing, collision, AI, sound), while those entities are being
locked and unlocked by themselves and each other, according to behaviors.

This is probably quite a bit heavier than what we could expect from the actual
gameplay.

First of all, in 

> [latency numbers every programmer should know](https://gist.github.com/jboner/2841832)

, a main memory reference comes in at 100 ns, so during the atomic entity lock
we simulate here,

```
var ATOMIC_MODIFY_DURATION = func() time.Duration {
	return time.Duration(rand.Intn(8000)) * time.Nanosecond
}
```

... we'd be looking at an average of around *40* main memory references 
(all of them cache misses!)

We also would likely not spawn this many entities, nor entities with behavior,
at one time:

```
const N_ENTITIES = 2000
const N_ENTITIES_WITH_BEHAVIOR = 1000
```

The entity amount chosen was based on a quick sketch of the types and numbers
of entities, how many behaviors they'd have, how often they'd run, which was 
done in a notebook. In making the rough calculations I erred on the side of 
a lot of entities. We'd probably want to, if running this many entities, do
some kind of partial deactivate based on distance, fading into a despawn zone.

The full set of constants tested initially:

From `consts.go`:
```
const N_ENTITIES = 2000
const N_ENTITIES_WITH_BEHAVIOR = 1600
const N_BEHAVIORS_PER_ENTITY = 8
const CHANCE_TO_SAFEGET_IN_BEHAVIOR = 0.3
const SAFEGET_DURATION = 4 * time.Microsecond
const CHANCE_TO_LOCK_OTHER_ENTITY = 0.05
const OTHER_ENTITY_LOCK_DURATION = 12 * time.Microsecond

const WORLD_WIDTH = 3200
const WORLD_HEIGHT = 3200

const FPS = 60
const FRAME_SLEEP = (1000 / FPS) * time.Millisecond

var ATOMIC_MODIFY_DURATION = func() time.Duration {
	return time.Duration(rand.Intn(8000)) * time.Nanosecond
}

var BEHAVIOR_POST_SLEEP = func() time.Duration {
	return time.Duration(time.Duration(rand.Intn(700)) * time.Millisecond)
}
```

From `spatial_hash.go`:
```
const GRID = 10
const CELL_WIDTH = WORLD_WIDTH / GRID
const CELL_HEIGHT = WORLD_HEIGHT / GRID
const SPATIAL_HASH_SCAN_GOROUTINES = 12
const SPATIAL_HASH_UPPER_ESTIMATE_ENTITIES_PER_SQUARE = 12
const GRID_GOROUTINES = 12
```

Despite the heavy load, on the following machine specs, admittedly a somewhat
fast CPU but by no means top of the line,

```
Intel(R) Core(TM) i5-3320M CPU @ 2.60GHz
go version go1.8.3 linux/amd64
```

... with the following constants set,

```
```

... we see the ComputeSpatialHash() algorithm taking 0 or 1 ms for each cycle,
and 16 or more every so often:


It's not immediately clear why the spikes occur (it may be a "perfect storm"
situation re: locking.

The below is data from pprof showing the time spent in various branches over a
test lasting about 11 seconds:

```
Type: cpu
Time: May 30, 2018 at 1:29am (EDT)
Duration: 11.49947704s
Total: 7.34s

main.(*SpatialHash).ComputeSpatialHash
/home/anon/go/src/github.com/dt-rush/donkeys-qquest/testing/collision/spatial_hash/spatial_hash.go

  Total:           0       30ms (flat, cum)  0.41%
    128            .          . } 
    129            .          .  
    130            .          . // spawns a certain number of goroutines to iterate through entities, trying 
    131            .          . // to lock them and get their position and send the entity to another goroutine 
    132            .          . // handling the building of the list for that cell 
    133            .          . func (h *SpatialHash) ComputeSpatialHash() { 
    134            .          .  
    135            .          . 	// lock the UpdatedEntityList from modifying itself while we read it 
    136            .          . 	h.spatialEntities.Mutex.Lock() 
    137            .          . 	defer h.spatialEntities.Mutex.Unlock() 
    138            .          .  
    139            .          . 	// mutex / deferred goroutine pattern means that 
    140            .          . 	// the reset of internal state used to check entities 
    141            .          . 	// will be done immediately after the function returns, 
    142            .          . 	// and in the background 
    143            .          . 	h.resetMutex.Lock() 
    144            .          . 	h.resetMutex.Unlock() 
    145            .          . 	defer func() { 
    146            .          . 		go h.ResetAlreadyReadArray() 
    147            .          . 	}() 
    148            .          .  
    149            .          . 	// waitgroup is used to ensure that every entity has been checked and every 
    150            .          . 	// grid has bene populated 
    151            .          . 	wg := sync.WaitGroup{} 
    152            .          . 	// divide the list of entities into a certain number of partitions 
    153            .          . 	partition_size := len(h.spatialEntities.Entities) / 
    154            .          . 		SPATIAL_HASH_SCAN_GOROUTINES 
    155            .          . 	for partition := 0; partition < SPATIAL_HASH_SCAN_GOROUTINES; partition++ { 
    156            .          . 		// the last partition includes the remainder 
    157            .          . 		offset := partition * partition_size 
    158            .          . 		if partition == SPATIAL_HASH_SCAN_GOROUTINES-1 { 
    159            .          . 			partition_size = len(h.spatialEntities.Entities) - offset 
    160            .          . 		} 
    161            .          . 		// spawn the cellSender goroutine for this partition 
    162            .          . 		wg.Add(1) 
    163            .          . 		go h.cellSender(offset, partition_size, &wg) 
    164            .          . 	} 
    165            .          . 	// for each cell, spawn a cellReceiver goroutine 
    166            .          . 	entitiesRemaining := atomic.NewUint32(N_ENTITIES) 
    167            .          . 	for y := 0; y < GRID; y++ { 
    168            .          . 		for x := 0; x < GRID; x++ { 
    169            .          . 			wg.Add(1) 
    170            .       30ms 			go h.cellReceiver(y, x, entitiesRemaining, &wg) 
    171            .          . 		} 
    172            .          . 	} 
    173            .          . 	wg.Wait() 
    174            .          . } 
    175            .          .  
    176            .          . func (h *SpatialHash) String() string { 
    177            .          . 	var buffer bytes.Buffer 
    178            .          . 	buffer.WriteString("[\n") 
    179            .          . 	for y := 0; y < GRID; y++ { 

main.(*SpatialHash).cellReceiver
/home/anon/go/src/github.com/dt-rush/donkeys-qquest/testing/collision/spatial_hash/spatial_hash.go

  Total:        20ms      320ms (flat, cum)  4.36%
     78            .          . 	h.resetMutex.Unlock() 
     79            .          . } 
     80            .          .  
     81            .          . // used to receive EntityPositions and put them into the right cell 
     82            .          . func (h *SpatialHash) cellReceiver( 
     83            .          . 	y int, x int, entitiesRemaining *atomic.Uint32, wg *sync.WaitGroup) { 
     84            .          .  
     85            .       30ms 	for entitiesRemaining.Load() > 0 { 
     86            .          . 		select { 
     87            .      100ms 		case entityPosition := <-h.cellChannels[y][x]: 
     88         20ms       60ms 			h.cells[y][x] = append(h.cells[y][x], entityPosition) 
     89            .       30ms 			entitiesRemaining.Dec() 
     90            .          . 		default: 
     91            .      100ms 			time.Sleep(2 * time.Microsecond) 
     92            .          . 		} 
     93            .          . 	} 
     94            .          . 	wg.Done() 
     95            .          . } 
     96            .          .  
     97            .          . // used to iterate the entities and send them to the right cell's 
     98            .          . // cellReceiver() instance (they are spawned, one for each cell, as goroutines) 
     99            .          . func (h *SpatialHash) cellSender( 
    100            .          . 	offset int, partition_size int, 

main.(*SpatialHash).cellSender
/home/anon/go/src/github.com/dt-rush/donkeys-qquest/testing/collision/spatial_hash/spatial_hash.go

  Total:        30ms      170ms (flat, cum)  2.32%
     96            .          .  
     97            .          . // used to iterate the entities and send them to the right cell's 
     98            .          . // cellReceiver() instance (they are spawned, one for each cell, as goroutines) 
     99            .          . func (h *SpatialHash) cellSender( 
    100            .          . 	offset int, partition_size int, 
    101            .          . 	wg *sync.WaitGroup) { 
    102            .          .  
    103            .          . 	// keep track of how many we've read 
    104            .          . 	n_read := 0 
    105         10ms       10ms 	for i := 0; n_read < partition_size; i = (i + 1) % partition_size { 
    106         10ms       10ms 		entity := h.spatialEntities.Entities[offset+i] 
    107            .          . 		if h.alreadyRead[offset+i] { 
    108            .          . 			continue 
    109            .          . 		} 
    110            .          . 		// attempt the lock 
    111            .          . 		if h.entityTable.attemptLockEntity(entity) { 
    112            .          . 			// if we locked, grab the position and send it to 
    113            .          . 			// the channel 
    114            .          . 			position := h.position.Data[entity.ID] 
    115            .       50ms 			h.entityTable.releaseEntity(entity) 
    116            .          . 			y := position[1] / CELL_HEIGHT 
    117            .          . 			x := position[0] / CELL_WIDTH 
    118            .       80ms 			h.cellChannels[y][x] <- EntityPosition{entity, position} 
    119         10ms       10ms 			h.alreadyRead[offset+i] = true 
    120            .          . 			n_read++ 
    121            .          . 			continue 
    122            .          . 		} 
    123            .          . 		// else, sleep a bit (to prevent hot loops if there are only 
    124            .          . 		// a few entities left and they are all locked) 
    125            .       10ms 		time.Sleep(10 * time.Microsecond) 
    126            .          . 	} 
    127            .          . 	wg.Done() 
    128            .          . } 
    129            .          .  
    130            .          . // spawns a certain number of goroutines to iterate through entities, trying 
    131            .          . // to lock them and get their position and send the entity to another goroutine 
    132            .          . // handling the building of the list for that cell 
    133            .          . func (h *SpatialHash) ComputeSpatialHash() { 

main.Behavior
/home/anon/go/src/github.com/dt-rush/donkeys-qquest/testing/collision/spatial_hash/behavior.go

  Total:        50ms         1s (flat, cum) 13.62%
      6            .          . ) 
      7            .          .  
      8            .          . // imitates a behavior func which will atomically modify the entity, 
      9            .          . // possibly run a safeget on another entity, possibly lock that other 
     10            .          . // entity, and sleep various amounts of time during all this 
     11            .          . func Behavior(t *EntityTable, entity EntityToken) { 
     12            .          . 	for { 
     13            .          . 		// simulating an AtomicEntityModify 
     14         20ms       70ms 		t.lockEntity(entity) 
     15         10ms      400ms 		time.Sleep(ATOMIC_MODIFY_DURATION()) 
     16         10ms       80ms 		if rand.Float32() < CHANCE_TO_SAFEGET_IN_BEHAVIOR { 
     17         10ms      110ms 			time.Sleep(SAFEGET_DURATION) 
     18            .          . 		} 
     19            .       30ms 		if rand.Float32() < CHANCE_TO_LOCK_OTHER_ENTITY { 
     20            .          . 			otherID := rand.Intn(N_ENTITIES) 
     21            .          . 			for otherID == entity.ID { 
     22            .          . 				otherID = rand.Intn(N_ENTITIES) 
     23            .          . 			} 
     24            .          . 			otherEntity := EntityToken{ID: otherID} 
     25            .          . 			t.lockEntity(otherEntity) 
     26            .       10ms 			time.Sleep(OTHER_ENTITY_LOCK_DURATION) 
     27            .          . 			t.releaseEntity(otherEntity) 
     28            .          . 		} 
     29            .       10ms 		t.releaseEntity(entity) 
     30            .      290ms 		time.Sleep(BEHAVIOR_POST_SLEEP()) 
     31            .          . 	} 
     32            .          . } 
```

Here is some representative data for the total time the `ComputeSpatialHash()`
function takes as profiled with good old fashioned `time.Since()`:

```
computing spatial hash took 0 ms
computing spatial hash took 6 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 1 ms
computing spatial hash took 3 ms
computing spatial hash took 1 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 1 ms
computing spatial hash took 1 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 1 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 1 ms
computing spatial hash took 0 ms
computing spatial hash took 1 ms
computing spatial hash took 1 ms
computing spatial hash took 1 ms
computing spatial hash took 1 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 17 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 1 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 0 ms
computing spatial hash took 1 ms
computing spatial hash took 0 ms
computing spatial hash took 2 ms
computing spatial hash took 1 ms
computing spatial hash took 0 ms
```

Data from a longer run of `30` s on the same machine mentioned above shows an
average runtime of around `0.369` ms. If we need a graphics frame every `16 ms`,
we are in a very comfortable position especially considering that the entity
locking heaviness was cranked way up beyond normal expectations.

Since the spatial hash is a prerequesite for physics, collision
detection, and draw, it's important that we have a spatial hash available 
every frame (aka, at `100` fps, every `16` ms). But it's *not* necessary that we 
have the *current* one available every frame. The spikes will cause frame drops
in the synchronous game loop (process keyboard/mouse, compute spatial hash, 
do physics / collision, draw)

We could, if we wanted to, overcome even the effect of these occassional spikes
by using a good old fashioned double-buffering technique with atomic operations 
to change a pointer to the cells 2D-array of lists while the next one is 
building. It would cost twice the space, but mean that the update cycle was never
interrupted, just at best working off of 1-frame-old position data when culling
for graphics or checking for collisions).

Note that as-implemented, we check *points*, not bounding boxes. But adding
another read is probably not going to do much since the reads of the positions
cost only around `20` ms cumulatively in a `10` second test of `600` frames,
or `33.3` microseconds per frame.

If we scale back the constants to something a bit more reasonable,

```
const N_ENTITIES = 1600
const N_ENTITIES_WITH_BEHAVIOR = 800
const N_BEHAVIORS_PER_ENTITY = 4

var ATOMIC_MODIFY_DURATION = func() time.Duration {
	return time.Duration(rand.Intn(4000)) * time.Nanosecond
}
```

... and modify the grid size since it was a bit too sparse at `32x32` (now
`5x5`, it was probably not as efficient to only have `12` partitions with so
many entities (now we have `32`), and increase the number of entities expected
per cell to `N_ENTITIES / GRID` from a constant (very low) of `12`,

```
const GRID = 5
const SPATIAL_HASH_UPPER_ESTIMATE_ENTITIES_PER_SQUARE = N_ENTITIES / GRID
const GRID_GOROUTINES = 32
```

... we instead get an average of `0.147` ms to compute the spatial hash, with
a frequency table of millisecond run-times looking like (for a `30` second
test):

```
0: 1707
1: 45
2: 14
3: 5
4: 3
5: 8
6: 7
7: 4
8: 4
9: 0
10: 1
11: 0
12: 0
13: 1
14: 0
15: 0
```

After some modifications put in in commit 
`f73fe9efd714dca215e4967b445ae9ba9643465b`, we get an average of `92`
*microseconds* per computation with the following parameters:

```
const N_ENTITIES = 1600
const N_ENTITIES_WITH_BEHAVIOR = 800
const N_BEHAVIORS_PER_ENTITY = 4
const CHANCE_TO_SAFEGET_IN_BEHAVIOR = 0.3
const SAFEGET_DURATION = 4 * time.Microsecond
const CHANCE_TO_LOCK_OTHER_ENTITY = 0.05
const OTHER_ENTITY_LOCK_DURATION = 12 * time.Microsecond

const WORLD_WIDTH = 5760
const WORLD_HEIGHT = 5760

const FPS = 60
const FRAME_SLEEP = (1000 / FPS) * time.Millisecond

var ATOMIC_MODIFY_DURATION = func() time.Duration {
	return time.Duration(rand.Intn(4000)) * time.Nanosecond
}

var BEHAVIOR_POST_SLEEP = func() time.Duration {
	return time.Duration(time.Duration(rand.Intn(700)) * time.Millisecond)
}

const GRID = 12
const SPATIAL_HASH_CELL_WIDTH = WORLD_WIDTH / GRID
const SPATIAL_HASH_CELL_HEIGHT = WORLD_HEIGHT / GRID

nScanPartitions := int(2*math.Pow(math.Log(N_ENTITIES+1), 2) + 1)
```
