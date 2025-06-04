package Lightshow

type switchType uint8

const (
	OnBeat switchType = iota
	InDeadTime
	InCalmMode
)

func (context *AutoModeContext) getNextAnimation(switchType switchType) AnimationsInstruction {

}

func (context *AutoModeContext) getNextLighting(switchType switchType) LightingInstruction {

}
