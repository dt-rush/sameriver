package sameriver

type GOAPTemporalGoal struct {
	temporalGoals []*GOAPGoal
}

func NewGOAPTemporalGoal(spec any) *GOAPTemporalGoal {
	tg := &GOAPTemporalGoal{}
	if specmap, single := spec.(map[string]int); single {
		tg.temporalGoals = []*GOAPGoal{NewGOAPGoal(specmap)}
	} else if specarr, temporal := spec.([]any); temporal {
		tg.temporalGoals = make([]*GOAPGoal, 0)
		for i := 0; i < len(specarr); i++ {
			specmapi := specarr[i].(map[string]int)
			tg.temporalGoals = append(tg.temporalGoals, NewGOAPGoal(specmapi))
		}
	}
	return tg
}

func (tg *GOAPTemporalGoal) Parametrize(n int) *GOAPTemporalGoal {
	for i, g := range tg.temporalGoals {
		tg.temporalGoals[i] = g.Parametrize(n)
	}
	return tg
}
