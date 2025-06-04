package Lightshow

import (
	"ControlApp/Infrastructure"
	"sync"
)

type VisualManager struct {
	autoContext        *AutoModeContext
	hardwareManager    Infrastructure.HardwareManager
	lightingOverwrite  *Infrastructure.LightingInstruction
	animationOverwrite []Infrastructure.AnimationInstruction
	textOverwrite      []Infrastructure.TextInstruction
	accessLock         *sync.Mutex
}

func CreateLightingManager(hardwareManager Infrastructure.HardwareManager) *VisualManager {
	visual := VisualManager{hardwareManager: hardwareManager}
	visual.autoContext = CreateAutoMode(visual, loadConfiguration())
	return &visual
}

func loadConfiguration() AutoModeConfiguration {
	return AutoModeConfiguration{}
}

func (manager VisualManager) applyLighting(instruction Infrastructure.LightingInstruction) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	if manager.lightingOverwrite != nil {
		return
	}

	manager.hardwareManager.SendLightingInstruction(instruction)
}

func (manager VisualManager) applyAnimation(instruction Infrastructure.AnimationInstruction) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	if manager.animationOverwrite != nil {
		return
	}

	manager.hardwareManager.SendAnimationInstruction(instruction)
}

func (manager VisualManager) triggerBeat() {
	manager.hardwareManager.SendBeatToDisplay(false)
}

func (manager VisualManager) SetLightingOverwrite(instruction *Infrastructure.LightingInstruction) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	manager.lightingOverwrite = instruction
	if instruction != nil {
		manager.hardwareManager.SendLightingInstruction(*instruction)
	}
}

func (manager VisualManager) SetAnimationsOverwrite(instructions []Infrastructure.AnimationInstruction) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	manager.animationOverwrite = instructions
	if instructions != nil {
		for _, animation := range instructions {
			manager.hardwareManager.SendAnimationInstruction(animation)
		}
	}
}

func (manager VisualManager) SetTextsOverwrite(instructions []Infrastructure.TextInstruction) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	manager.textOverwrite = instructions
	if instructions != nil {
		for _, text := range instructions {
			manager.hardwareManager.SendTextInstruction(text)
		}
	}
}
