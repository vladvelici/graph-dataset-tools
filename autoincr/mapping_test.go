package main

import "testing"

func TestBasicAddLookup(t *testing.T) {
	index := NewMapping()

	addNodes := []int{3, 5, 6, 1, 9}
	givenAllocations := make([]int, len(addNodes))

	for i, nid := range addNodes {
		alloc, existing := index.Node(nid)
		if existing {
			t.Errorf("Node %d exists. It shouldn't.", nid)
		}
		givenAllocations[i] = alloc
	}

	for i, given := range givenAllocations {
		// autoinrement assignment starting from 1
		if given != i+1 {
			t.Errorf("(Autoincr) Allocation for node %d should be %d (found %d).", addNodes[i], i+1, given)
		}

		// reverse lookup
		if g, found := index.Allocation(given); !found || g != addNodes[i] {
			t.Errorf("(reverse lookup) Wrong allocation. Expected (%d, true), found (%d, %b).", addNodes[i], g, found)
		}

		// normal lookup
		if g, found := index.Node(addNodes[i]); !found || g != given {
			t.Errorf("(lookup) Wrong. Expected (%d, true), found (%d, %b).", given, g, found)
		}
	}
}
