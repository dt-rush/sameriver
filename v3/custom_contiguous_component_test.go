package sameriver

import (
	"regexp"
	"testing"
)

type XYZ struct {
	x int
	y int
	z int
}

type XYZComponent struct {
	data []XYZ
}

func (xyz *XYZComponent) Name() string {
	return "XYZ"
}
func (xyz *XYZComponent) AllocateTable(n int) {
	xyz.data = make([]XYZ, n, 2*n)
}
func (xyz *XYZComponent) ExpandTable(n int) {
	extraSpace := make([]XYZ, n)
	xyz.data = append(xyz.data, extraSpace...)
}
func (xyz *XYZComponent) Get(e *Entity) interface{} {
	return xyz.data[e.ID]
}
func (xyz *XYZComponent) Set(e *Entity, x interface{}) {
	xyz.data[e.ID] = x.(XYZ)
}

func TestCCCGetSet(t *testing.T) {
	w := testingWorld()
	xyz := &XYZComponent{}
	w.RegisterCCCs([]CustomContiguousComponent{
		xyz,
	})
	// spawn entity with empty base CS, XYZ custom CS
	e := w.em.Spawn(map[string]any{
		"customComponents": map[string]any{
			"Custom,XYZ": XYZ{x: 1, y: 0, z: 8},
		},
		"customComponentsImpl": map[string]CustomContiguousComponent{
			"XYZ": xyz,
		},
	})

	// get value and check
	firstGet := e.GetCustom("XYZ").(XYZ)
	expected := XYZ{x: 1, y: 0, z: 8}
	if firstGet != expected {
		t.Errorf("Didn't Get() properly")
	}
	// modify local copy
	firstGet.x = 3
	firstGet.y = 3
	firstGet.z = 3
	// set and check
	e.SetCustom("XYZ", firstGet)
	secondGet := e.GetCustom("XYZ").(XYZ)
	expected = XYZ{x: 3, y: 3, z: 3}
	if secondGet != expected {
		t.Errorf("Didn't Set() properly")
	}
}

func TestCCCBitArray(t *testing.T) {
	w := testingWorld()
	xyz := &XYZComponent{}
	w.RegisterCCCs([]CustomContiguousComponent{
		xyz,
	})
	// spawn entity with empty base CS, XYZ custom CS
	e := w.em.Spawn(map[string]any{
		"customComponents": map[string]any{
			"Custom,XYZ": XYZ{x: 1, y: 0, z: 8},
		},
		"customComponentsImpl": map[string]CustomContiguousComponent{
			"XYZ": xyz,
		},
	})

	b := e.ComponentBitArray
	s := w.em.components.BitArrayToString(b)
	Logger.Println(s)
	// TODO: test string for XYZ
	if valid, _ := regexp.MatchString("XYZ", s); !valid {
		t.Errorf("XYZ not in component bit array -> string")
	}
}
