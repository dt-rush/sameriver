package sameriver

import (
	"fmt"
	"math/rand"
	"regexp"
	"testing"
	"time"
)

func TestComponentBitArrayToString(t *testing.T) {
	w := testingWorld()
	N_COMPONENT_TYPES := len(w.em.components.names)
	if len(w.em.components.names) <= 1 {
		// there should always be at least GenericTags and GenericLogic
		t.Fatal("why in god's name are there 1 or less component types???")
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 16; i++ {
		c1 := w.em.components.ixs_rev[rand.Intn(N_COMPONENT_TYPES)]
		c2 := c1
		for c2 == c1 {
			c2 = w.em.components.ixs_rev[rand.Intn(N_COMPONENT_TYPES)]
		}
		b := w.em.components.BitArrayFromNames([]string{c1, c2})
		s := w.em.components.BitArrayToString(b)
		found1, _ := regexp.MatchString(c1, s)
		found2, _ := regexp.MatchString(c2, s)
		if !found1 {
			t.Fatal(fmt.Sprintf("name of component %s not found in %s\n",
				c1, s))
		}
		if !found2 {
			t.Fatal(fmt.Sprintf("name of component %s not found in %s\n",
				c2, s))
		}
	}
}
