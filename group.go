package main

import (
	"fmt"
)

// Group represents a collection of devices
type Group struct {
	ID      int
	Name    string
	Devices []int
}

// Groups represents a collection of groups
type Groups struct {
	Groups map[int]*Group // Keyed on group id
	dl     DeviceLibrary
}

func createGroups(dl DeviceLibrary) *Groups {
	return &Groups{
		Groups: make(map[int]*Group),
		dl:     dl}
}

// getGroup returns the group as a slice (instead of map which is original).
// returns empty slice if group don't exit.
func (g *Groups) getGroups() []*Group {
	var result []*Group
	for _, group := range g.Groups {
		result = append(result, group)
	}
	return result
}

func (g *Groups) add(newGroup *Group) error {
	_, ok := g.Groups[newGroup.ID]
	if ok == true {
		return fmt.Errorf("Group with id %d already exists", newGroup.ID)
	}
	g.Groups[newGroup.ID] = newGroup
	return nil
}

// get returns nil if group don't exists
func (g *Groups) get(groupID int) *Group {
	group, ok := g.Groups[groupID]
	if ok == false {
		return nil
	}
	return group
}

// trunOn turns on all devices in group that supports OnOff. Returns error
// if group don't exists.
func (g *Groups) turnOn(groupID int) error {
	group, ok := g.Groups[groupID]
	if ok == false {
		return fmt.Errorf("Group with id %d does not exist", groupID)
	}
	for _, id := range group.Devices {
		if g.dl.SupportsOnOff(id) {
			g.dl.TurnOn(id)
		}
	}
	return nil
}

// trunOff turns off all devices in group that supports OnOff. Returns error
// if group don't exists.
func (g *Groups) turnOff(groupID int) error {
	group, ok := g.Groups[groupID]
	if ok == false {
		return fmt.Errorf("Group with id %d does not exist", groupID)
	}
	for _, id := range group.Devices {
		if g.dl.SupportsOnOff(id) {
			g.dl.TurnOff(id)
		}
	}
	return nil
}

// dim dims all devices in group that supports Dim. Returns error
// if group don't exists.
func (g *Groups) dim(groupID int, level byte) error {
	group, ok := g.Groups[groupID]
	if ok == false {
		return fmt.Errorf("Group with id %d does not exist", groupID)
	}
	for _, id := range group.Devices {
		if g.dl.SupportsDim(id) {
			g.dl.Dim(id, level)
		}
	}
	return nil
}
