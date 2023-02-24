package sameriver

import (
	"math/rand"
	"testing"
	"time"
)

func TestIDGeneratorUnique(t *testing.T) {
	IDGen := NewIDGenerator()
	IDs := make(map[int]bool)
	for i := 0; i < 1024*1024; i++ {
		ID := IDGen.Next()
		if _, ok := IDs[ID]; ok {
			t.Fatal("produced same ID already produced")
		}
		IDs[ID] = true
	}
}

func TestIDGeneratorUniqueRemoval(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	IDGen := NewIDGenerator()
	IDs := make(map[int]bool)
	for i := 0; i < 1024*1024; i++ {
		if i > 1024 && i%2 == 0 {
			ID := rand.Intn(len(IDGen.universe))
			IDGen.Free(ID)
			delete(IDs, ID)
		} else {
			ID := IDGen.Next()
			if _, ok := IDs[ID]; ok {
				t.Fatal("produced same ID already produced")
			}
			IDs[ID] = true
		}
	}
}
