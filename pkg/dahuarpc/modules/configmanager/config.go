package configmanager

import (
	"fmt"
	"slices"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type SwitchMode int

const (
	SwitchModeGeneral SwitchMode = iota
	SwitchModeDay
	SwitchModeNight
	SwitchModeTime
	SwitchModeBrightness
)

func (m SwitchMode) String() string {
	switch m {
	case SwitchModeGeneral:
		return "general"
	case SwitchModeDay:
		return "day"
	case SwitchModeNight:
		return "night"
	case SwitchModeTime:
		return "time"
	case SwitchModeBrightness:
		return "brightness"
	default:
		return "unknown"
	}
}

type VideoInMode struct {
	Config      []int                    `json:"Config"`
	Mode        int                      `json:"Mode"`
	TimeSection [][]dahuarpc.TimeSection `json:"TimeSection"`
}

func (m VideoInMode) SwitchMode() (SwitchMode, error) {
	if m.Mode == 0 && slices.Equal(m.Config, []int{2}) {
		return SwitchModeGeneral, nil
	}
	if m.Mode == 0 && slices.Equal(m.Config, []int{0}) {
		return SwitchModeDay, nil
	}
	if m.Mode == 0 && slices.Equal(m.Config, []int{1}) {
		return SwitchModeNight, nil
	}
	if m.Mode == 1 && slices.Equal(m.Config, []int{0, 1}) {
		return SwitchModeTime, nil
	}
	if m.Mode == 2 && slices.Equal(m.Config, []int{0, 1}) {
		return SwitchModeBrightness, nil
	}
	return 0, fmt.Errorf("unknown SwitchMode")
}

func (m VideoInMode) SetSwitchMode(mode SwitchMode) VideoInMode {
	switch mode {
	case SwitchModeGeneral:
		m.Mode = 0
		m.Config = []int{2}
	case SwitchModeDay:
		m.Mode = 0
		m.Config = []int{0}
	case SwitchModeNight:
		m.Mode = 0
		m.Config = []int{1}
	case SwitchModeTime:
		m.Mode = 1
		m.Config = []int{0, 1}
	case SwitchModeBrightness:
		m.Mode = 2
		m.Config = []int{0, 1}
	}
	return m
}
