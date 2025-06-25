package Lightshow

import (
	"ControlApp/BoxiBus"
	"ControlApp/Display"
	"ControlApp/Infrastructure"
	"sync"
)

type VisualManager struct {
	autoContext                   *AutoModeContext
	animations                    *AnimationManager
	palettes                      *PaletteManager
	hardwareManager               Infrastructure.HardwareInterface
	lightingIsOverwritten         bool
	animationIsOverwritten        bool
	lightingCurrentAutoSelection  LightingInstruction
	animationCurrentAutoSelection AnimationsInstruction
	textOverwrite                 *TextsInstruction
	accessLock                    *sync.Mutex
}

func CreateVisualManager(hardwareManager Infrastructure.HardwareInterface) *VisualManager {
	visual := VisualManager{hardwareManager: hardwareManager, accessLock: &sync.Mutex{}}
	visual.animations = LoadAnimations()
	visual.palettes = LoadPalettes()
	visual.autoContext = CreateAutoMode(&visual, loadConfiguration())

	// Sync animations when they get uploaded
	go visual.watchForAnimationUploads()

	return &visual
}

type LightingInstruction struct {
	BoxiBus.MessageBlock
	character ModeCharacter
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

func (manager *VisualManager) applyLighting(instruction LightingInstruction) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	manager.lightingCurrentAutoSelection = instruction
	if manager.animationIsOverwritten {
		return
	}

	manager.hardwareManager.SendLightingInstruction(instruction.MessageBlock)
}

func (manager *VisualManager) applyAnimation(instruction AnimationsInstruction) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	manager.animationCurrentAutoSelection = instruction
	if manager.animationIsOverwritten {
		return
	}

	for _, animation := range instruction.animations {
		manager.hardwareManager.SendAnimationInstruction(animation.animation, animation.displays)
	}

	manager.hardwareManager.SendBrightnessChange(nil, instruction.blinkSpeed)
}

func (manager *VisualManager) triggerBeat() {
	manager.hardwareManager.SendBeatToDisplay(false)
}

func (manager *VisualManager) getBeatState() bool {
	return manager.hardwareManager.GetBeatState()
}

func (manager *VisualManager) getAnimations() *AnimationManager {
	return manager.animations
}

func (manager *VisualManager) getPalettes() *PaletteManager {
	return manager.palettes
}

func (manager *VisualManager) SetLightingOverwrite(instruction *LightingInstruction) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	manager.lightingIsOverwritten = instruction != nil
	if instruction == nil {
		manager.hardwareManager.SendLightingInstruction(manager.lightingCurrentAutoSelection.MessageBlock)
	} else {
		manager.hardwareManager.SendLightingInstruction(instruction.MessageBlock)
	}
}

func (manager *VisualManager) SetAnimationsOverwrite(instructions *AnimationsInstruction) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	manager.animationIsOverwritten = instructions != nil
	if instructions == nil {
		for _, animation := range manager.animationCurrentAutoSelection.animations {
			manager.hardwareManager.SendAnimationInstruction(animation.animation, animation.displays)
		}
	} else {
		for _, animation := range instructions.animations {
			manager.hardwareManager.SendAnimationInstruction(animation.animation, animation.displays)
		}
	}
}

func (manager *VisualManager) SetTextsOverwrite(instructions *TextsInstruction) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	manager.textOverwrite = instructions
	if instructions != nil {
		for _, text := range *instructions {
			manager.hardwareManager.SendTextInstruction(text.text, text.displays)
		}
	}
}

func (manager *VisualManager) getAllAnimations() []Display.AnimationId {
	var ids []Display.AnimationId

	for _, data := range manager.animations.animations {
		ids = append(ids, data.Id)
	}

	return ids
}

func (manager *VisualManager) ImportAnimation(path string, name string, mood LightingMood, splitVideo bool) (Display.AnimationId, error) {
	return manager.animations.ImportAnimation(path, name, mood, splitVideo)
}

func (manager *VisualManager) watchForAnimationUploads() {
	for {
		animationId := <-manager.animations.UploadQueue
		manager.hardwareManager.UploadAnimation(animationId)
	}
}
