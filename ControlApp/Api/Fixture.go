package Api

import (
	"ControlApp/Infrastructure"
	"ControlApp/Lightshow"
)

type Fixture struct {
	Hardware Infrastructure.HardwareInterface
	Visuals  Lightshow.VisualManager
}
