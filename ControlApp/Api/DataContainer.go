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

func CreateDataContainer(hardware Infrastructure.HardwareInterface, visuals *Lightshow.VisualManager) *DataContainer {
	result := DataContainer{
		hardware, visuals, LightingInstructionTotal{
			Enable:           false,
			ApplyOnBeat:      false,
			Mode:             0,
			ColorDeviceA:     Color{R: 50, UV: 255},
			ColorDeviceB:     Color{G: 10, A: 200},
			PaletteId:        0,
			DurationMs:       10000,
			PaletteShift:     0,
			Speed:            0,
			TargetBrightness: 0,
			FrequencyHz:      12,
		},
	}

	return &result
}
