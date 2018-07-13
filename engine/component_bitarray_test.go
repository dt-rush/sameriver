package engine

import (
	"fmt"
	"math/rand"
	"regexp"
	"testing"
	"time"
)

func TestComponentBitArrayToString(t *testing.T) {
	if N_COMPONENT_TYPES <= 1 {
		t.Fatal("why in god's name are there 1 or less component types???")
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 16; i++ {
		c1 := ComponentType(rand.Intn(N_COMPONENT_TYPES))
		c2 := c1
		for c2 == c1 {
			c2 = ComponentType(rand.Intn(N_COMPONENT_TYPES))
		}
		b := MakeComponentBitArray([]ComponentType{c1, c2})
		s := ComponentBitArrayToString(b)
		found1, _ := regexp.MatchString(COMPONENT_NAMES[c1], s)
		found2, _ := regexp.MatchString(COMPONENT_NAMES[c2], s)
		if !found1 {
			t.Fatal(fmt.Sprintf("name of component %s not found in %s\n",
				COMPONENT_NAMES[c1], s))
		}
		if !found2 {
			t.Fatal(fmt.Sprintf("name of component %s not found in %s\n",
				COMPONENT_NAMES[c1], s))
		}
	}
}
