package main

import "testing"

func TestMstSet(t *testing.T) {
	mst := make(Mst)

	mst.Add(4, 3)
	if !mst.Has(3, 4) {
		t.Error("simple. Does not have links from 3 to 4.")
	}
	if !mst.Has(4, 3) {
		t.Error("simple. Does not have links from 4 to 3.")
	}

	mst.Add(3, 10)

	if !mst.Has(10, 3) || !mst.Has(3, 10) {
		t.Error("second. Does not have a link.")
	}

}
