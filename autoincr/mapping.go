package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type Mapping struct {
	index       map[int]int `json:"_"`
	allocations []int
}

func NewMapping() *Mapping {
	return &Mapping{
		make(map[int]int),
		make([]int, 0),
	}
}

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
	m.index = make(map[int]int, len(m.allocations))
	for given, id := range m.allocations {
		m.index[id] = given
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
	if given, ok := m.index[id]; ok {
		return given, true
	}
	alloc := len(m.allocations)
	m.index[id] = alloc
	m.allocations = append(m.allocations, id)
	return alloc, false
}

// Returns (given, false) if not allocated. (allocation, true) otherwise.
func (m *Mapping) Allocation(given int) (int, bool) {
	if given >= len(m.allocations) || given < 0 {
		return given, false
	}
	return m.allocations[given], true
}

// Removes a node, given by real mapping.
func (m *Mapping) Remove(id int) bool {
	given, ok := m.index[id]
	if !ok {
		return false
	}
	delete(m.index, id)
	m.allocations[given] = m.allocations[len(m.allocations)-1]
	m.allocations = m.allocations[:len(m.allocations)-1]
	return true
}
