package main

import (
	"fmt"
	"time"
)

// Weekday represents a day in the week
type Weekday string

// Weekday definitions
const (
	Monday    Weekday = "mon"
	Tuesday   Weekday = "tue"
	Wednesday Weekday = "wed"
	Thursday  Weekday = "thu"
	Friday    Weekday = "fri"
	Saturday  Weekday = "sat"
	Sunday    Weekday = "sun"
)

type Restriction string

const (
	SunUp         Restriction = "sunUp"
	SunDown       Restriction = "sunDown"
	NoRestriction Restriction = ""
)

type Function string

const (
	TurnOn  Function = "turnOn"
	TurnOff Function = "turnOff"
	Dim     Function = "dim"
)

type TimeEvent struct {
	Hour        int       // -1 if OnSunset or OnSunrise
	Minute      int       // -1 if OnSunset or OnSunrise
	OnSunrise   bool      // false if Hour, Minute != -1 or OnSunset = true
	OnSunset    bool      // false if Hour, Minute != -1 or OnRise =
	Offset      int       // Minutes, can be negative (only for OnSunset/OnSunrise).
	Days        []Weekday // 1 - 7 elements
	Restriction Restriction
	Function    Function
	DimLevel    int // -1 if no dim
	DeviceID    int // -1 if group
	GroupID     int // -1 if device
}

func parseEvent(eventConfig string) *TimeEvent {
	return nil
}

// checkTrigger checks if it is time to trigger the event. It only checks the minutes, not
// the seconds. If it is time to trig the function will be called.
func (t TimeEvent) checkTrigger(now time.Time, sunset time.Time, sunrise time.Time) {
	fmt.Println("Hello")
}
