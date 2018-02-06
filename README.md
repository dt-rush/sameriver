
donkeys-qquest
===

### What is Donkeys QQuest?

RPG Quest for Donkey Based Plot Items, Action Adventures, Shopkeepers and Dungeons

Will you save the world? Will you rescue Old Man Willis's garden from gophers? What secrets are hidden under the church? Who *are* you and *why are you alive*?

> Inspecting donkey corpse...
>     you found: 
>         donkey pelt x 1
>         donkey ears x 2
>         donkey hooves x 4
>         donkey whiskers x 32
>         gold x 100
>         ruby x 3
>         magic shield x 2
>         health potion x 2

-- ancient proverb



### 1. Technical details

#### a. General engine design

###### game logic

The *game logic* is built on an "entity-component-system" architecture, in which:

**Components** are collections of a certain type of data indexed by the ID's of entities. For example, a position component is at bottom a map[int]\([2]float64\)

**Entities** are merely the set of components which their ID's index, and are essentially passed around in the system *as identical with* their ID's.

**Systems** are collections of logic which operate on subsets of components. 
 
###### display logic
 
The *display logic* is built on a "scene-based" architecture, in which:

**Scenes** are responsible for actually running and displaying game content. They contain various components (in the future, only *references* to components, all of which will be registered and stored with the singleton Game object) and systems needed to support their operation. They are updated each game loop iteration, receiving: 

* keyboard state via a call to a `HandleKeyboardState (keyboard_state []uint8)` method
* delta-time updates via a call to an `Update (dt_ms float64)` method
* possibly^1 a call to a `Draw (window *sdl.Window, renderer *sdl.Renderer)` method.

Currently, only one active scene can exist at a time, and scenes are destroyed as they pass their successor scene to the game loop (they are initialized and loaded in the background while a singleton loading scene will be displayed until the new scene is ready to take over). In the future, it will be possible to push scenes to a stack and pop them off without destroying the underlying prior scene (ie. for cinematics, menu navigation, "battles", etc.)

^1. It's possible that the game loop will not draw every iteration in order to keep a certain framerate


##### b. Package / folder structure

All subfolders (other than `assets/`) are named for the packages of the source files they contain.

The packages described generally:

* engine

The game engine structs and functions, providing the abstract backbone on which the actual game logic is built. Can and should be abstractly separated in future from the other folders / packages which are the actual donkeys-qquest content (a game *buiilt* on the engine)

* engine/components

Definitions for the different types of components

* engine/systems

Definitions of systems

* constants

Definitions of game constants (should not be used for constant *content*, which should probably be in .json files in `assets/`)

* logic

Definitions of game logic (functions, really), to be loaded into various systems / scenes

* main

Holds main.go, the executable entrypoint, from whichh the game engine is initialized, and the game loop run (TODO: migrate even game loop logic into engine)

* scenes

Definitions of scenes 

* utils

Holds utilites for development, debugging, profiling, content-generation



