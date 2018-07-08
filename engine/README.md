River
===

Game Engine taking advantage of go's concurrency

> **Heraclitus** of Ephesus (/ˌhɛrəˈklaɪtəs/; Greek: Ἡράκλειτος ὁ Ἐφέσιος
> Hērákleitos ho Ephésios; c. 535 – c. 475 BC) was a pre-Socratic Greek
> philosopher, and a native of the city of Ephesus, then part of the 
> Persian Empire. Little is known about his early life and education, but 
> he regarded himself as self-taught and a pioneer of wisdom.
> 
> Heraclitus was famous for his insistence on ever-present change as
> being the fundamental essence of the universe, as stated in the famous
> saying, "No man ever steps in the same river twice" (see panta rhei
> below). This position was complemented by his stark commitment to
> a unity of opposites in the world, stating that "the path up and down
> are one and the same". Through these doctrines Heraclitus characterized all
> existing entities by pairs of contrary properties, whereby no entity may
> ever occupy a single state at a single time. This, along with his cryptic
> utterance that "all entities come to be in accordance with this Logos
> " (literally, "word", "reason", or "account") has been the subject of
> numerous interpretations.

# 1. Technical Details

The engine takes advantage of language features to define concurrently-executing entity behaviour relative to a traditional synchronous game loop ("input, update, draw").

#### 1.a. General engine design

##### 1.a.i. entity component system

The engine is built on an "entity-component-system" architecture, in which:

**Components** are collections of a certain type of data indexed by the ID's of entities. For example, a position component is at bottom a `map[int]([2]int16)`

**Entities** are merely the set of components which their ID's index, and are essentially passed around in the system *as identical with* their ID's.

**Systems** are collections of logic which operate on subsets of components.

There are also some **Managers** which are sort of like the glue holding the engine together, or providing services.

The **`EntityManager`** is a particularly important central part of the engine.

##### 1.a.ii. scenes

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

#### 1.b. Build system

The engine is, without events or components specified by a game using the engine, not of much use. The code-generation system by which the game author can integrate their components and events into the engine will be described in a future modification to this README, when the system has taken more concrete form. This is a TODO.
