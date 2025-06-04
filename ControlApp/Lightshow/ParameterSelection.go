package Lightshow

import "ControlApp/Infrastructure"

type switchType uint8

const (
	OnBeat switchType = iota
	InDeadTime
	InSlowMode
)

func (context *AutoModeContext) getNextAnimation(switchType switchType) []Infrastructure.AnimationInstruction {

}

func (context *AutoModeContext) getNextLighting(switchType switchType) Infrastructure.LightingInstruction {

}
