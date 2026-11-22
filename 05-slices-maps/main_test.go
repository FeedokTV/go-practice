package slicesmaps

import (
	"maps"
	"slices"
	"testing"
)

func TestSlices(t *testing.T) {

	tests := []struct {
		Name             string
		NewElems         []int
		WantReallocation bool
		WantNewCap       bool
		WantLen          int
		BaseCap          int
		BaseSliceLow     int
		BaseSliceHigh    int
		BaseSliceMax     int
	}{
		{
			Name:             "fits in cap",
			NewElems:         []int{1, 2},
			WantReallocation: false,
			WantNewCap:       false,
			WantLen:          4,
			BaseCap:          4,
			BaseSliceLow:     0,
			BaseSliceHigh:    2,
			BaseSliceMax:     4,
		},
		{
			Name:             "reallocating",
			NewElems:         []int{1, 2, 3},
			WantReallocation: true,
			WantNewCap:       true,
			WantLen:          5,
			BaseCap:          4,
			BaseSliceLow:     0,
			BaseSliceHigh:    2,
			BaseSliceMax:     4,
		},
		{
			Name:             "full slice expression (force reallocation)",
			NewElems:         []int{1, 2},
			WantReallocation: true,
			WantNewCap:       true,
			WantLen:          4,
			BaseCap:          4,
			BaseSliceLow:     0,
			BaseSliceHigh:    2,
			BaseSliceMax:     2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			baseSlice := make([]int, tt.BaseSliceHigh, tt.BaseCap)

			testSlice := baseSlice[tt.BaseSliceLow:tt.BaseSliceHigh:tt.BaseSliceMax]

			newSlice := append(testSlice, tt.NewElems...)

			if !tt.WantReallocation && &newSlice[0] != &baseSlice[0] {
				t.Errorf("unexpected reallocation")
			}

			if tt.WantReallocation && &newSlice[0] == &baseSlice[0] {
				t.Errorf("expected reallocation")
			}

			if len(newSlice) != tt.WantLen {
				t.Errorf("wanted len: %d, got: %d", tt.WantLen, len(newSlice))
			}

			if !tt.WantNewCap && cap(newSlice) != cap(testSlice) {
				t.Errorf("unexpected new cap, got: %d", cap(newSlice))
			}

			if tt.WantNewCap && cap(newSlice) == cap(testSlice) {
				t.Errorf("expect a new cap, got: %d", cap(newSlice))
			}
		})
	}
}

func TestPredictOverwriting(t *testing.T) {
	baseSlice := make([]int, 2, 4)

	firstSlice := append(baseSlice, 1, 2)

	if &firstSlice[0] != &baseSlice[0] {
		t.Error("unexpected assertion")
	}

	secondSlice := append(baseSlice, 10, 11)

	if &secondSlice[0] != &baseSlice[0] {
		t.Error("unexpected assertion")
	}

	if firstSlice[2] != 10 {
		t.Errorf("expected overwriting and elem 10, got: %d", firstSlice[2])
	}

	if len(firstSlice) != 4 {
		t.Errorf("unexpected len of firstSlice want: 4, got %d", len(firstSlice))
	}

	if cap(firstSlice) != 4 {
		t.Errorf("unexpected cap of firstSlice want: 4, got %d", cap(firstSlice))
	}

	if len(secondSlice) != 4 {
		t.Errorf("unexpected len of secondSlice want: 4, got %d", len(secondSlice))
	}

	if cap(secondSlice) != 4 {
		t.Errorf("unexpected cap of secondSlice want: 4, got %d", cap(secondSlice))
	}
}

// Map addressing is bad! Don't do that

func TestMapAddressingBad(t *testing.T) {

	type someStruct struct {
		k int
	}

	s := someStruct{k: 1}

	mappy := map[string]someStruct{"a": s}
	// Maps fields are unaddressable
	// mappy["a"].k = 2 // cannot assign to struct field mappy["a"].k in map

	// To do what we tried above
	temp := mappy["a"]
	temp.k = 2
	mappy["a"] = temp

	if mappy["a"].k != 2 {
		t.Errorf("expected value 2, got: %d", mappy["a"].k)
	}

	anotherMappy := map[string]*someStruct{"a": &s}

	anotherMappy["a"].k = 2 // this one is good! Because go go *(mappy["a"]).k automatically

	if anotherMappy["a"].k != 2 {
		t.Errorf("expected value 2, got: %d", anotherMappy["a"].k)
	}

}

func TestSortingMap(t *testing.T) {
	mappy := map[string]int{"z": 3, "a": 1, "r": 2}

	keys := slices.Sorted(maps.Keys(mappy))

	if len(keys) != len(mappy) {
		t.Errorf("amount of sorted keys are less than in original map")
	}

	for i, k := range keys {
		if mappy[k] != i+1 {
			t.Errorf("unexpected sorting, want: %d, got %d", i+1, mappy[k])
		}
	}
}
