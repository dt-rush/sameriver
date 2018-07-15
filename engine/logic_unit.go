package engine

import (
	"github.com/dt-rush/sameriver/engine/utils"
)

type LogicUnit struct {
	Name    string
	F       func()
	Active  bool
	WorldID utils.ID
}
