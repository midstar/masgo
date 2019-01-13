package main

import (
	"fmt"
	"strconv"
	"strings"
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

// parseGroup creates a group from the configuration syntax:
//   GROUP <groupID> "<name>" <ID1> [<ID2> .. <IDn>]
func parseGroup(configStr string) (*Group, error) {
	sections := strings.Split(configStr, "\"")
	if len(sections) != 3 {
		return nil, fmt.Errorf("Group name must be within \"")
	}
	words := strings.Fields(sections[0])
	if len(words) != 2 {
		return nil, fmt.Errorf("Group id is missing")
	}
	if words[0] != "GROUP" {
		return nil, fmt.Errorf("This is not a group configuration. Expected GROUP got %s", words[0])
	}
	groupID, err := strconv.Atoi(words[1])
	if err != nil {
		return nil, fmt.Errorf("invalid group id %s", words[1])
	}
	groupName := sections[1]
	deviceIDStrs := strings.Fields(sections[2])
	nbrDevices := len(deviceIDStrs)
	devices := make([]int, nbrDevices)
	for i := 0; i < nbrDevices; i++ {
		deviceIDStr := deviceIDStrs[i]
		devices[i], err = strconv.Atoi(deviceIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid device id %s", deviceIDStr)
		}
	}
	result := Group{
		ID:      groupID,
		Name:    groupName,
		Devices: devices}
	return &result, nil
}

func (g *Group) toConfigStr() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("GROUP %d \"%s\"", g.ID, g.Name))
	for _, id := range g.Devices {
		builder.WriteString(fmt.Sprintf(" %d", id))
	}
	return builder.String()
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
