package sameriver

/*
NOTE: non-temporal plans were faster in 097d939, before we treated every plan
implicitly as temporal. We should use the algorithm from 097d939 opportunistically,
creating JIT the zero'd trivial regionoffsets for the so-far-inserted actions
and starting thenceforth to do the full temporalgoal insert and update routine

for example:

let's say we've been opportunistically using insert-before and simpler insert (no regionoffsets)
since no temporal goals have been encountered yet

main g

A.pre q
Y.pre r
B.pre s
C.pre u

A satisfies u
B satisfies u

...

[A Y B C]
remainings [[q] [r] [s] [t] [u] [v]]
path   : [ A   Y   B   C ]
parents: [ C   B   C   * ]

now we want to insert X to satisfy A's pre q, and X has temporal pre [m n]

[X A Y B C]
remainings [[m n] [r] [s] [t] [u] [v]]
path   : [ X   A   Y   B   C ]
parents: [ A   C   B   C   * ]

what should region offsets be?

[[m n] [r] [s] [t] [u] [v]]
[[0 0] [0] [0] [0] [0] [0]]

start to use insertionIx relative to regions, and start to update regions
*/

type GOAPTemporalGoal struct {
	temporalGoals []*GOAPGoal
}

func NewGOAPTemporalGoal(spec any) *GOAPTemporalGoal {
	tg := &GOAPTemporalGoal{}
	if specmap, single := spec.(map[string]int); single {
		tg.temporalGoals = []*GOAPGoal{newGOAPGoal(specmap)}
	} else if specarr, temporal := spec.([]any); temporal {
		tg.temporalGoals = make([]*GOAPGoal, 0)
		for i := 0; i < len(specarr); i++ {
			specmapi := specarr[i].(map[string]int)
			tg.temporalGoals = append(tg.temporalGoals, newGOAPGoal(specmapi))
		}
	} else {
		tg.temporalGoals = []*GOAPGoal{}
	}
	return tg
}

func (tg *GOAPTemporalGoal) Parametrized(n int) *GOAPTemporalGoal {
	for i, g := range tg.temporalGoals {
		tg.temporalGoals[i] = g.Parametrized(n)
	}
	return tg
}
