package Lightshow

import (
	"ControlApp/BoxiBus"
	"ControlApp/Display"
	"math/rand"
)

type switchType uint8

const (
	FirstBeat switchType = iota
	OnBeat
	InDeadTime
	InCalmMode
)

const (
	defaultBlinkSpeed = 20
)

func (context *AutoModeContext) getNextAnimation(switchType switchType) AnimationsInstruction {
	baseMood := context.Configuration.Mood

	//When in a calmer section of a beat mode, randomly pick between moody and happy
	if (baseMood == Regular || baseMood == Party) && switchType != OnBeat {
		randNbr := rand.Intn(2)
		if randNbr == 0 {
			baseMood = Moody
		} else {
			baseMood = Happy
		}
	}

	animationManager := context.switcher.getAnimations()
	animationManager.accessLock.Lock()
	defer animationManager.accessLock.Unlock()

	//Find valid animations to switch to
	validIndices := make([]int, 0)
	for index, animation := range animationManager.animations {
		if animation.Mood == baseMood || animation.Mood == Regular && baseMood == Party {
			validIndices = append(validIndices, index)
		}
	}

	if len(validIndices) < 1 {
		return AnimationsInstruction{character: Unknown}
	}

	var dsp1A, dsp1B, dsp2A, dsp2B Display.AnimationId

	mirrorAcrossScreens := rand.Intn(2)
	generateBoxiScreens := func() (Display.AnimationId, Display.AnimationId) {
		randomIndex := rand.Intn(len(validIndices))
		firstAnimation := animationManager.animations[validIndices[randomIndex]]

		//If picked animation is played across two screens, do that
		if firstAnimation.SecondaryAnimation != Display.None {
			var secondAnimation Animation
			foundSecondAnimation := false

			for _, a := range animationManager.animations {
				if a.Id == firstAnimation.SecondaryAnimation {
					secondAnimation = a
					foundSecondAnimation = true
					break
				}
			}

			if foundSecondAnimation {
				return firstAnimation.Id, secondAnimation.Id
			}
		}

		if mirrorAcrossScreens != 0 {
			return firstAnimation.Id, firstAnimation.Id
		}

		randomIndex = rand.Intn(len(validIndices))
		secondAnimation := animationManager.animations[validIndices[randomIndex]]
		return firstAnimation.Id, secondAnimation.Id
	}

	dsp1A, dsp1B = generateBoxiScreens()

	mirrorAcrossBoxis := rand.Intn(2)
	if mirrorAcrossBoxis == 0 {
		dsp2A = dsp1A
		dsp2B = dsp1B
	} else {
		dsp2A, dsp2B = generateBoxiScreens()
	}

	doDaBounce := rand.Intn(7)
	var blinkSpeed uint16
	if doDaBounce == 6 && baseMood == Party {
		blinkSpeed = defaultBlinkSpeed
	}

	character := Calm
	if baseMood == Regular || baseMood == Party {
		character = Rhythmic
	}

	//Find grouped animations
	screensPerAnimation := make(map[Display.AnimationId]Display.ServerDisplay)
	screensPerAnimation[dsp1A] |= Display.Boxi1D1
	screensPerAnimation[dsp1B] |= Display.Boxi1D2
	screensPerAnimation[dsp2A] |= Display.Boxi2D1
	screensPerAnimation[dsp2B] |= Display.Boxi2D2

	instructions := make([]animationInstruction, 0)
	for animationId, display := range screensPerAnimation {
		instructions = append(instructions, animationInstruction{animationId, []Display.ServerDisplay{display}})
	}

	return AnimationsInstruction{instructions, character, blinkSpeed}
}

func (context *AutoModeContext) getNextLighting(switchType switchType) LightingInstruction {
	baseMood := context.Configuration.Mood
	var possibleModes []BoxiBus.LightingModeId

	if (baseMood == Regular || baseMood == Party) && switchType != OnBeat {
		//When in a calmer section of a beat mode, randomly pick between moody and happy
		randNbr := rand.Intn(2)
		if randNbr == 0 {
			baseMood = Moody
		} else {
			baseMood = Happy
		}

		// Also only allow the transition Beat -> FadeToColor -> PaletteFade
		if switchType == InDeadTime {
			possibleModes = []BoxiBus.LightingModeId{BoxiBus.FadeToColor}
		} else {
			possibleModes = []BoxiBus.LightingModeId{BoxiBus.PaletteFade}
		}
	} else {
		//If in any other section of any other mode, just pick mode from any of the default mode list
		possibleModes = getLightingModesByMood(baseMood)
	}

	randNbr := rand.Intn(len(possibleModes))
	mode := possibleModes[randNbr]

	randNbr = rand.Intn(context.Configuration.HueShiftChance)
	hueShift := randNbr == 0

	palette := []BoxiBus.Color{{255, 255, 255, 0, 0, 0}}
	// TODO: Implement palette management/selection

	if switchType == FirstBeat && context.Configuration.StrobeChance > 0 {
		randNbr = rand.Intn(context.Configuration.StrobeChance)
		if randNbr == 0 {
			mode = BoxiBus.Strobe
			palette = []BoxiBus.Color{{0, 0, 0, 255, 0, 0}}
		}
	}

	// TODO: Convert mode, palette and hue shift parameter to mode struct and return
}

func getLightingModesByMood(mood LightingMood) []BoxiBus.LightingModeId {
	switch mood {
	case Happy, Moody:
		return []BoxiBus.LightingModeId{BoxiBus.FadeToColor, BoxiBus.PaletteFade}
	case Regular, Party:
		return []BoxiBus.LightingModeId{BoxiBus.PaletteSwitch, BoxiBus.PaletteBrightnessFlash, BoxiBus.PaletteHueFlash}
	}
}
