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
	defaultBlinkSpeed = 600
)

func (context *AutoModeContext) getNextAnimation(switchType switchType) AnimationsInstruction {
	baseMood := context.Configuration.Mood

	//When in a calmer section of a beat mode, randomly pick between moody and happy
	if (baseMood == Regular || baseMood == Party) && (switchType == InDeadTime || switchType == InCalmMode) {
		randNbr := rand.Intn(2)
		if randNbr == 0 {
			baseMood = Moody
		} else {
			baseMood = Happy
		}
	}

	animationManager := context.manager.GetAnimations()
	animationManager.accessLock.Lock()
	defer animationManager.accessLock.Unlock()

	//Find valid animations to switch to
	validIds := make([]Display.AnimationId, 0)
	for id, animation := range animationManager.animations {
		if animation.IsNsfw && !context.Configuration.AllowNsfw {
			continue
		}

		if animation.Mood == baseMood || animation.Mood == Regular && baseMood == Party {
			validIds = append(validIds, id)
		}
	}

	if len(validIds) < 1 {
		return AnimationsInstruction{Character: Unknown}
	}

	var dsp1A, dsp1B, dsp2A, dsp2B Display.AnimationId

	mirrorAcrossScreens := rand.Intn(2)
	generateBoxiScreens := func() (Display.AnimationId, Display.AnimationId) {
		randomIndex := rand.Intn(len(validIds))
		firstAnimation := animationManager.animations[validIds[randomIndex]]

		//If picked animation is played across two screens, do that
		if firstAnimation.SecondaryAnimation != Display.None {
			return firstAnimation.SecondaryAnimation, firstAnimation.Id
		}

		if mirrorAcrossScreens != 0 {
			return firstAnimation.Id, firstAnimation.Id
		}

		randomIndex = rand.Intn(len(validIds))
		secondAnimation := animationManager.animations[validIds[randomIndex]]
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

	instructions := make([]AnimationInstruction, 0)
	for animationId, display := range screensPerAnimation {
		instructions = append(instructions, AnimationInstruction{animationId, []Display.ServerDisplay{display}})
	}

	return AnimationsInstruction{instructions, character, blinkSpeed}
}

func (context *AutoModeContext) getNextLighting(switchType switchType) LightingInstruction {
	baseMood := context.Configuration.Mood
	var possibleModes []BoxiBus.LightingModeId
	var possiblePalettes []Palette

	if (baseMood == Regular || baseMood == Party) && (switchType == InDeadTime || switchType == InCalmMode) {
		//When in a calmer section of a beat mode, randomly pick between moody and happy
		randNbr := rand.Intn(2)
		if randNbr == 0 {
			baseMood = Moody
		} else {
			baseMood = Happy
		}

		// Also only allow the transition Beat -> FadeToColor -> PaletteFade
		if switchType == InCalmMode {
			possibleModes = []BoxiBus.LightingModeId{BoxiBus.FadeToColor}
			possiblePalettes = []Palette{
				{0, "UV", []BoxiBus.Color{{255, 0, 0, 0, 0, 255}}, nil},
				{1, "Blue", []BoxiBus.Color{{0, 0, 255, 0, 0, 0}}, nil},
				{2, "Amber", []BoxiBus.Color{{0, 0, 0, 0, 255, 0}}, nil},
			}
		} else {
			possibleModes = []BoxiBus.LightingModeId{BoxiBus.PaletteFade}
		}
	}

	if possibleModes == nil {
		possibleModes = getLightingModesByMood(baseMood)
	}
	randNbr := rand.Intn(len(possibleModes))
	mode := possibleModes[randNbr]

	if possiblePalettes == nil {
		possiblePalettes = context.manager.GetPalettes().GetPalettesForMood(baseMood)
	}
	if possiblePalettes == nil || len(possiblePalettes) == 0 {
		possiblePalettes = getDefaultPalettes()
	}
	randNbr = rand.Intn(len(possiblePalettes))
	palette := possiblePalettes[randNbr].Colors

	randNbr = rand.Intn(context.Configuration.HueShiftChance)
	hueShift := 0

	if randNbr == 0 {
		hueShift = rand.Intn(len(palette))
	}

	// Figure out whether the mode should be applied on the next beat
	applyOnNextBeat := false
	if switchType == OnBeat {
		applyOnNextBeat = true
	}

	// If first beat since a while, have a chance for strobe to flash bang you
	if switchType == FirstBeat && context.Configuration.StrobeChance > 0 {
		randNbr = rand.Intn(context.Configuration.StrobeChance)
		if randNbr == 0 {
			mode = BoxiBus.Strobe
			palette = []BoxiBus.Color{{0, 0, 0, 255, 0, 0}}
			applyOnNextBeat = false
		}
	}

	messages := getLightingMessages(context.Configuration, mode, palette, byte(hueShift), applyOnNextBeat)
	character := getLightingModeCharacter(mode)
	return LightingInstruction{messages, character}
}

func getLightingModesByMood(mood LightingMood) []BoxiBus.LightingModeId {
	switch mood {
	case Happy, Moody:
		return []BoxiBus.LightingModeId{BoxiBus.FadeToColor, BoxiBus.PaletteFade}
	default:
		return []BoxiBus.LightingModeId{BoxiBus.PaletteSwitch, BoxiBus.PaletteBrightnessFlash, BoxiBus.PaletteHueFlash}
	}
}

func getLightingMessages(config AutoModeConfiguration, mode BoxiBus.LightingModeId, palette []BoxiBus.Color, hueShift byte, applyOnBeat bool) BoxiBus.MessageBlock {
	switch mode {
	case BoxiBus.Off:
		return BoxiBus.CreateLightingOff(applyOnBeat)
	case BoxiBus.SetColor:
		return BoxiBus.CreateLightingSetColor(palette[0], palette[int(hueShift)%len(palette)], applyOnBeat)
	case BoxiBus.FadeToColor:
		return BoxiBus.CreateLightingFadeToColor(palette[0], palette[int(hueShift)%len(palette)], config.FadeToColorCycles, applyOnBeat)
	case BoxiBus.PaletteFade:
		result, err := BoxiBus.CreateLightingPaletteFade(palette, config.PaletteFadeCycles, hueShift, applyOnBeat)
		if err == nil {
			return result
		}
	case BoxiBus.PaletteSwitch:
		result, err := BoxiBus.CreateLightingPaletteSwitch(palette, hueShift, applyOnBeat)
		if err == nil {
			return result
		}
	case BoxiBus.PaletteBrightnessFlash:
		result, err := BoxiBus.CreateLightingPaletteBrightnessFlash(palette, config.FlashFadeoutSpeed, config.FlashTargetBrightness, hueShift, applyOnBeat)
		if err == nil {
			return result
		}
	case BoxiBus.PaletteHueFlash:
		result, err := BoxiBus.CreateLightingPaletteHueFlash(palette, config.HueFlashFadeoutSpeed, config.FlashHueShift, applyOnBeat)
		if err == nil {
			return result
		}
	case BoxiBus.Strobe:
		return BoxiBus.CreateLightingStrobe(palette[0], config.StrobeFrequency, config.StrobeRolloff, applyOnBeat)
	}

	return nil
}
