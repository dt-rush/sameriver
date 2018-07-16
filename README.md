sameriver
===

[![Build Status](https://travis-ci.com/dt-rush/sameriver.svg?branch=master)](https://travis-ci.com/dt-rush/sameriver)
[![codecov](https://codecov.io/gh/dt-rush/sameriver/branch/master/graph/badge.svg)](https://codecov.io/gh/dt-rush/sameriver)
[![CodeFactor](https://www.codefactor.io/repository/github/dt-rush/sameriver/badge)](https://www.codefactor.io/repository/github/dt-rush/sameriver)

![](https://tokei.rs/b1/github/dt-rush/sameriver)
![](https://tokei.rs/b1/github/dt-rush/sameriver?category=comments)
![](https://tokei.rs/b1/github/dt-rush/sameriver?category=files)

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
concurrently-executing entity and world logic relative to a 
traditional synchronous game loop ("input, update, draw").

### 1. Development

Run `make` to build and test.

Dependencies can be installed by `make deps`.

### 2. Dependencies

apt:

* `libsdl2{,-mixer,-image,-ttf,-gfx}-dev`

pacman:

* `sdl2{,_mixer,_image,_ttf,_gfx}`

windows: mingw env packages (install from source)

* `SDL2-devel-2.0.5-mingw.tar.gz`
* `SDL2_image-devel-2.0.1-mingw.tar.gz`
* `SDL2_mixer-devel-2.0.1-mingw.tar.gz`
* `SDL2_ttf-devel-2.0.14-mingw.tar.gz`


### 3. Technical details

#### 3.a. General engine design

##### 3.a.i. entity component system

The engine is built on an "entity-component-system" architecture, in which:

**Components** are collections of a certain type of data indexed by the ID's of entities. For example, the velocity component is a `[MAX_ENTITIES]Vec2D`.

**Entities** are merely the set of components indexed by an ID (entities can also be active or inactive)

**Systems** are collections of logic which operate on subsets of entities selected for by an arbitrary query

There are also some **Managers** which are sort of like the glue holding the engine together, or providing services.

The **`EntityManager`** is a particularly important central part of the engine.

##### 3.a.ii. scenes

The engine is also built on a "scene-based" architecture, in which:

**Scenes** are loaded into the `Game` object's loop, and have control over input and display while they're running.

All scenes will be registered and stored with the singleton Game object, and can refer to each other by name (strings).

The currently running scene is updated each game loop iteration, receiving:

* keyboard state via a call to a `HandleKeyboardState (keyboard_state []uint8)` method
* delta-time updates via a call to an `Update (dt_ms float64)` method
* possibly^1 a call to a `Draw (window *sdl.Window, renderer *sdl.Renderer)` method.

Scenes are initialized and loaded in the background while a singleton loading scene will be displayed until the new scene is ready to take over.

^1. It's possible that the game loop will not draw every iteration in order to keep a certain framerate
