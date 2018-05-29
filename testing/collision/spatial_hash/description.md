# Description 

This script was made to profile the performance of sorting entities into a 
spatial hash (useful for restricting the scope of calculations needed in 
physics, drawing, collision, AI, sound), while those entities are being
locked and unlocked by themselves and each other, according to behaviors.

This is probably quite a bit heavier than what we could expect from the actual
gameplay.

First of all, in [latency numbers every programmer should know](https://gist.github.com/jboner/2841832),
a main memory reference comes in at 100 ns, so during the atomic entity lock
we simulate here, we'd be looking at around *80* main memory references 
(all of them cache misses).

We also would likely not spawn this many entities at one time. The entity amount
chosen was based on a quick sketch of the types and numbers of entities, how
many behaviors they'd have, how often they'd run, which was done in a notebook.
In making the rough calculations I erred on the side of a lot of entities. We'd
probably want to, if running this many entities, do some kind of partial
deactivate based on distance, fading into a despawn zone.

Despite the heavy load, on the following machine specs, admittedly a somewhat
fast CPU but by no means top of the line...

```
Intel(R) Core(TM) i5-3320M CPU @ 2.60GHz
go version go1.8.3 linux/amd64
```

... we saw the DoSpatialHash() function taking between 0 ms and 16 ms for
most every cycle. Since the spatial hash is a prerequesite for physics, collision
detection, and draw, it's important that we have one available every frame
(aka, at 100 fps, every 16 ms). But it's *not* necessary that we have the *current*
one available every frame. We could overcome this ocassional slowness by using a
good old fashioned double-buffering technique with atomic operations to change
a pointer to the spatial hash table.

