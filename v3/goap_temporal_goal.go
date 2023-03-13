package sameriver

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
	}
	return tg
}

func (tg *GOAPTemporalGoal) Parametrized(n int) *GOAPTemporalGoal {
	for i, g := range tg.temporalGoals {
		tg.temporalGoals[i] = g.Parametrized(n)
	}
	return tg
}
