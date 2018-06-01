sameriver
===

> Heraclitus of Ephesus (/ˌhɛrəˈklaɪtəs/;[1] Greek: Ἡράκλειτος ὁ Ἐφέσιος
> Hērákleitos ho Ephésios; c. 535 – c. 475 BC) was a pre-Socratic Greek
> philosopher, and a native of the city of Ephesus,[2] then part of
> the Persian Empire.
>
> Heraclitus was famous for his insistence on ever-present change as
> being the fundamental essence of the universe, as stated in the famous
> saying, "No man ever steps in the same river twice"[6] (see panta rhei
> below). This position was complemented by his stark commitment to
> a unity of opposites in the world, stating that "the path up and down
> are one and the same". Through these doctrines Heraclitus characterized all
> existing entities by pairs of contrary properties, whereby no entity may
> ever occupy a single state at a single time. This, along with his cryptic
> utterance that "all entities come to be in accordance with this Logos
> " (literally, "word", "reason", or "account") has been the subject of
> numerous interpretations.

The engine takes advantage of go's language features to define 
concurrently-executing entity behaviour and world logic relative to a 
traditional synchronous game loop ("input, update, draw").

Note: **The engine is in heavy development by a single developer, and so is
probably not in a very readable state at any given time until this comment is
removed from the README. Comments in code files have not been reviewed
thoroughly. Do not believe their lies.**

### 1. Dependencies

#### 1.a. linux

apt:

* libsdl2{,-mixer,-image,-ttf,-gfx}-dev

pacman:

* sdl2{,_mixer,_image,_ttf,_gfx}

#### 1.b. windows: mingw env packages (install from source)

* SDL2-devel-2.0.5-mingw.tar.gz
* SDL2_image-devel-2.0.1-mingw.tar.gz
* SDL2_mixer-devel-2.0.1-mingw.tar.gz
* SDL2_ttf-devel-2.0.14-mingw.tar.gz


#### 1.c. go packages (go get)

* github.com/veandco/go-sdl2/sdl
* github.com/veandco/go-sdl2/mix
* github.com/veandco/go-sdl2/img
* github.com/veandco/go-sdl2/ttf
* github.com/golang-collections/go-datastructures/bitarray
* go.uber.org/atomic
* github.com/dave/jennifer


### 2. Technical details

See the [wiki](https://github.com/dt-rush/donkeys-qquest/wiki) for diagrams.

#### 2.a. General engine design

##### 2.a.i. entity component system

The engine is built on an "entity-component-system" architecture, in which:

**Components** are collections of a certain type of data indexed by the ID's of entities. For example, a position component is at bottom a `map[int]([2]int16)`

**Entities** are merely the set of components which their ID's index, and are essentially passed around in the system *as identical with* their ID's.

**Systems** are collections of logic which operate on subsets of components.

There are also some **Managers** which are sort of like the glue holding the engine together, or providing services.

The **`EntityManager`** is a particularly important central part of the engine.

##### 2.a.ii. scenes

The engine is also built on a "scene-based" architecture, in which:

**Scenes** are loaded into the `Game` object's loop, and have control over input and display while they're running.

The most important scene, the GameScene, holds an `EntityManager`.

All scenes will be registered and stored with the singleton Game object, and can refer to each other by name (strings).

They are updated each game loop iteration, receiving:

* keyboard state via a call to a `HandleKeyboardState (keyboard_state []uint8)` method
* delta-time updates via a call to an `Update (dt_ms float64)` method
* possibly^1 a call to a `Draw (window *sdl.Window, renderer *sdl.Renderer)` method.

Scenes are initialized and loaded in the background while a singleton loading scene will be displayed until the new scene is ready to take over.

^1. It's possible that the game loop will not draw every iteration in order to keep a certain framerate
