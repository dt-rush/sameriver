sameriver
===

    ##XxxxxXXx:;;;;;;;;;xo;;;;;;;;;;;;;;:::::::::::ooooooooxxXX#X##X
    ###00OO0OOo:,,,,,,;:OO;;;;;;;;;;;;;;;;;::::::::::ooxoooxx00X####
    ###XOOx0xxxo,,,,,,;xO0o:x;,,,,,,;;;;;;;;;;;;;:::oO0#oxx0x00X####
    XXXX00Oxxxo:,,,,,,:x0XxxOx,,,,,,,,,,,,,;;::;;;o:xO0#0XX0O00X###X
    XXX####xooo;,,,;,ox0#XOxOOo,,,,,,,,,,,oo,xxOOxO0OOXXXX#X00O000XX
    #######Oxxoo,:x0OxOX#Oxxxxx,,,,,,,,;ooOOxxO##X0#00X0X####XOXX#X0
    #######XxxxXoOOX#OOX#XxOxxx,......,oxx00XOO##0X#0X#X#####0O0XX00
    ########0#X#OxO##X0###x0OOo,...,.,oxOO0X#O0##OX0XXOX#####0O0X#X#
    #########000xO0###X###0XOXx;,,ox:xOOOO0#X0XX#0X##XX#####XXXXX0X#
    ###########0xO0###########XxO00#XXX0O00#OOX#0OX#####X#####XXXX##
    ############00X############XX######000X#00OOOO###XOOxxOXXX0OXOOO
    ###XXO#OxXOOOOX#Xxx#0####XX#X#X###XOxxOxxxOxxOxxxxxxxxxxxxxxxxxx
    ###OXX#OxxxxxxxxxxOO0OX##xxxxxxoooxxOxxxxxxxxxx0Xxxxxxxxxxoxxxxx
    OOOOOxO#0OxxxX0xxxOxxxxxoooo;;.   ,:oxxxxxxxxxxOOxxoxxxxxxxoxxxx
    0OOxx00###0X0X#000O0OOxxo:;         .  .;oo:ooxxOOOxxxxxxxxxxxxx
    ##X0OOxxxxxxOOxxxxo;,..      .;o::xxxxxxxxxx:,,..,ooxOOxxxxxxxxx
    #####XX0xo;:::,.....      .   .xxxxxoxOxxxxxo,.,oo:oxxxxxxxOOX0X
    ####Xx::,.   ....      .    .,xxxxxxoox: .oxxo::xxxo;oxxO00X0xxx
    xo;:,..                .,oooxxxxxOxOOxxxxoxxxxxxxoxO0xo;oxxO##OX
            .............,;oxxxxxxxxxOOOOOOOOxxxoooxxxoxxxxO0X0##0xx
    ,. . ...........,,;:ooxxxxOOOOOxxxxOOO0XO0OOOx:;oxOxxOOxO0OO##XX
    o;;,,,,...,,,,;:oooxxxxxxOOOOOOOOOOOOO00XXXX##OxxxxOO0OOO#XXXX##
    ,..,,,,,,,;;;oooooooxxxOOOO0XX00000000X0X0XXXXX#XOxOOO00XX######
    ..,,,,:;;:;:ooooxxxoxxxxOOO0X#XXXXX000000000X0O00OO0XX00X##XX#X#

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

---

### 0. What is it?

A game engine which takes advantage of go's language features to define 
concurrently-executing entity behaviour and world logic relative to a 
traditional synchronous game loop ("input, update, draw").

**NOTE: The engine is in heavy development by a single developer, and so is
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

#### NOTE on `go.uber.org/atomic`

Yes, we use a dependency from the open source go libraries at Uber. As a taxi company, they're awful and exploitative. As a taxi company that employs go programmers, they're still awful and exploitative. Those programmers made an improvement to how go handle's atomic primitives by enforcing that those primitives be accessed atomically if they are meant to be, and that's cool. It improves readability and ensures that atomic primitives are always handled atomically.


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
