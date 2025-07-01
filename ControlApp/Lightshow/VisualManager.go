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
	animationOverwrite            *AnimationsInstruction
	brightnessValue               uint8
	lightingCurrentAutoSelection  LightingInstruction
	animationCurrentAutoSelection AnimationsInstruction
	textValues                    TextsInstruction
	accessLock                    *sync.Mutex
}

func CreateVisualManager(hardwareManager Infrastructure.HardwareInterface) *VisualManager {
	visual := VisualManager{hardwareManager: hardwareManager, accessLock: &sync.Mutex{}, brightnessValue: 0xFF}
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

type AnimationInstruction struct {
	Animation Display.AnimationId
	Displays  []Display.ServerDisplay
}

type AnimationsInstruction struct {
	Animations []AnimationInstruction
	Character  ModeCharacter
	BlinkSpeed uint16
}

type TextInstruction struct {
	Text     string
	Displays []Display.ServerDisplay
}
type TextsInstruction []TextInstruction

func (manager *VisualManager) applyLighting(instruction LightingInstruction) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	manager.lightingCurrentAutoSelection = instruction
	if manager.lightingIsOverwritten {
		return
	}

	manager.hardwareManager.SendLightingInstruction(instruction.MessageBlock)
}

func (manager *VisualManager) applyAnimation(instruction AnimationsInstruction) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	manager.animationCurrentAutoSelection = instruction

	// Collect overwritten displays
	var displaysToExclude []int
	if manager.animationOverwrite != nil {
		for _, animation := range manager.animationOverwrite.Animations {
			for _, display := range animation.Displays {
				for i := 1; i < 5; i++ {
					if int(display)&i != 0 {
						displaysToExclude = append(displaysToExclude, i)
					}
				}
			}
		}
	}

	for _, animation := range instruction.Animations {

		displays := make(map[int]bool)

		// Add displays to view with
		for _, display := range animation.Displays {
			for i := 1; i < 5; i++ {
				if int(display)&i != 0 {
					displays[i] = true
				}
			}
		}

		//Remove overridden display
		for _, display := range displaysToExclude {
			displays[display] = false
		}

		var displaySlice []Display.ServerDisplay
		for id, isIt := range displays {
			if !isIt {
				continue
			}

			displaySlice = append(displaySlice, Display.ServerDisplay(id))
		}

		manager.hardwareManager.SendAnimationInstruction(animation.Animation, displaySlice)

	}

	if manager.animationOverwrite == nil {
		manager.hardwareManager.SendBrightnessChange(nil, instruction.BlinkSpeed)
	}
}

func (manager *VisualManager) triggerBeat() {
	manager.hardwareManager.SendBeatToDisplay(false)
}

func (manager *VisualManager) getBeatState() bool {
	return manager.hardwareManager.GetBeatState()
}

func (manager *VisualManager) GetAnimations() *AnimationManager {
	return manager.animations
}

func (manager *VisualManager) GetPalettes() *PaletteManager {
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

	manager.animationOverwrite = instructions
	if instructions == nil {
		for _, animation := range manager.animationCurrentAutoSelection.Animations {
			manager.hardwareManager.SendAnimationInstruction(animation.Animation, animation.Displays)
		}
		manager.hardwareManager.SendBrightnessChange(nil, manager.animationCurrentAutoSelection.BlinkSpeed)
	} else {
		for _, animation := range instructions.Animations {
			manager.hardwareManager.SendAnimationInstruction(animation.Animation, animation.Displays)
		}
		manager.hardwareManager.SendBrightnessChange(nil, instructions.BlinkSpeed)
	}
}

func (manager *VisualManager) SetTexts(instructions TextsInstruction) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	valueSent := make(map[Display.ServerDisplay]bool)
	manager.textValues = instructions
	for _, text := range instructions {
		manager.hardwareManager.SendTextInstruction(text.Text, text.Displays)

		for _, display := range text.Displays {
			valueSent[display] = true
		}
	}

	var displaysToClear []Display.ServerDisplay
	for _, display := range manager.hardwareManager.GetConnectedDisplays() {
		success, result := valueSent[display]
		if !success || !result {
			displaysToClear = append(displaysToClear, display)
		}
	}

	manager.hardwareManager.SendTextInstruction(" ", displaysToClear)
}

func (manager *VisualManager) SetBrightness(value float64) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	bright := value
	manager.hardwareManager.SendBrightnessChange(&bright, 0)
}

func (manager *VisualManager) getAllAnimations() []Display.AnimationId {
	var ids []Display.AnimationId

	for _, data := range manager.animations.animations {
		ids = append(ids, data.Id)
	}

	return ids
}

func (manager *VisualManager) ImportAnimation(path string, name string, mood LightingMood, splitVideo bool, isNsfw bool) (Display.AnimationId, error) {
	return manager.animations.ImportAnimation(path, name, mood, splitVideo, isNsfw)
}

func (manager *VisualManager) GetConfiguration() *AutoModeConfiguration {
	return &manager.autoContext.Configuration
}

func (manager *VisualManager) MarkLightshowAsDirty() {
	manager.autoContext.isDirty = true
}

func (manager *VisualManager) watchForAnimationUploads() {
	for {
		animationId := <-manager.animations.UploadQueue
		manager.hardwareManager.UploadAnimation(animationId)
	}
}
