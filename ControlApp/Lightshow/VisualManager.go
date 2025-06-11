package Lightshow

import (
	"ControlApp/BoxiBus"
	"ControlApp/Display"
	"ControlApp/Infrastructure"
	"sync"
)

type VisualManager struct {
	autoContext        *AutoModeContext
	animations         *AnimationManager
	palettes           *PaletteManager
	hardwareManager    Infrastructure.HardwareManager
	lightingOverwrite  *LightingInstruction
	animationOverwrite *AnimationsInstruction
	textOverwrite      *TextsInstruction
	accessLock         *sync.Mutex
}

func CreateLightingManager(hardwareManager Infrastructure.HardwareManager) *VisualManager {
	visual := VisualManager{hardwareManager: hardwareManager}
	visual.animations = LoadAnimations()
	visual.palettes = LoadPalettes()
	visual.autoContext = CreateAutoMode(visual, loadConfiguration())
	return &visual
}

type LightingInstruction struct {
	BoxiBus.MessageBlock
	character ModeCharacter
	SlowMode  bool
}

type animationInstruction struct {
	animation Display.AnimationId
	displays  []Display.ServerDisplay
}

type AnimationsInstruction struct {
	animations []animationInstruction
	character  ModeCharacter
	blinkSpeed uint16
}

type textInstruction struct {
	text     string
	displays []Display.ServerDisplay
}
type TextsInstruction []textInstruction

func (manager VisualManager) applyLighting(instruction LightingInstruction) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	if manager.lightingOverwrite != nil {
		return
	}

	manager.hardwareManager.SendLightingInstruction(instruction.MessageBlock)
}

func (manager VisualManager) applyAnimation(instruction AnimationsInstruction) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	if manager.animationOverwrite != nil {
		return
	}

	for _, animation := range instruction.animations {
		manager.hardwareManager.SendAnimationInstruction(animation.animation, animation.displays)
	}
	manager.hardwareManager.SendBrightnessChange(nil, instruction.blinkSpeed)
}

func (manager VisualManager) triggerBeat() {
	manager.hardwareManager.SendBeatToDisplay(false)
}

func (manager VisualManager) getAnimations() *AnimationManager {
	return manager.animations
}

func (manager VisualManager) getPalettes() *PaletteManager {
	return manager.palettes
}

func (manager VisualManager) SetLightingOverwrite(instruction *LightingInstruction) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	manager.lightingOverwrite = instruction
	if instruction != nil {
		manager.hardwareManager.SendLightingInstruction(instruction.MessageBlock)
	}
}

func (manager VisualManager) SetAnimationsOverwrite(instructions *AnimationsInstruction) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	manager.animationOverwrite = instructions
	if instructions != nil {
		for _, animation := range instructions.animations {
			manager.hardwareManager.SendAnimationInstruction(animation.animation, animation.displays)
		}
	}
}

func (manager VisualManager) SetTextsOverwrite(instructions *TextsInstruction) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	manager.textOverwrite = instructions
	if instructions != nil {
		for _, text := range *instructions {
			manager.hardwareManager.SendTextInstruction(text.text, text.displays)
		}
	}
}
