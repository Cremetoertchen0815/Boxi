package Api

import (
	"ControlApp/Infrastructure"
	"ControlApp/Lightshow"
)

type DataContainer struct {
	Hardware                Infrastructure.HardwareInterface
	Visuals                 *Lightshow.VisualManager
	OverrideLightingCurrent LightingInstructionTotal
}
