package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type Mapping struct {
	Index       map[int]int `json:"-"`
	Allocations []int
}

// Create an empty mapping.
func NewMapping() *Mapping {
	return &Mapping{
		make(map[int]int),
		[]int{-1},
	}
}

// ReadMapping reads a mapping from a reader.
func ReadMapping(file io.Reader) (*Mapping, error) {
	raw, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var m Mapping
	err = json.Unmarshal(raw, &m)
	if err != nil {
		return nil, err
	}

	// make index from allocations
	m.Index = make(map[int]int, len(m.Allocations))
	for given, id := range m.Allocations {
		m.Index[id] = given
	}

	return &m, nil
}

// Marshal json and write into file.
func (m *Mapping) Write(file io.Writer) error {
	raw, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = file.Write(raw)
	return err
}

// Adds or looks up a node to the index, returning its allocated id.
// Returns (allocatedId, whether it was there before).
func (m *Mapping) Node(id int) (int, bool) {
	if given, ok := m.Index[id]; ok {
		return given, true
	}
	alloc := len(m.Allocations)
	m.Index[id] = alloc
	m.Allocations = append(m.Allocations, id)
	return alloc, false
}

// Returns (given, false) if not allocated. (allocation, true) otherwise.
func (m *Mapping) Allocation(given int) (int, bool) {
	if given >= len(m.Allocations) || given <= 0 {
		return given, false
	}
	return m.Allocations[given], true
}

// Removes a node, given by real mapping.
func (m *Mapping) Remove(id int) bool {
	given, ok := m.Index[id]
	if !ok {
		return false
	}

	delete(m.Index, id)
	m.Allocations[given] = m.Allocations[len(m.Allocations)-1]
	m.Allocations = m.Allocations[:len(m.Allocations)-1]
	return true
}
