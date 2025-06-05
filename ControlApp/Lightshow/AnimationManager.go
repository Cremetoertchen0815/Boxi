package Lightshow

import (
	"ControlApp/Display"
	"sync"
)

type Animation struct {
	Id                 Display.AnimationId
	Mood               LightingMood
	SecondaryAnimation Display.AnimationId
}

type AnimationManager struct {
	animations []Animation
	accessLock *sync.Mutex
}

func LoadAnimations() AnimationManager {
	return AnimationManager{}
}
