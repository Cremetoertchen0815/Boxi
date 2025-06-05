package Lightshow

import (
	"ControlApp/Display"
	"math/rand"
	"slices"
)

type switchType uint8

const (
	OnBeat switchType = iota
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

	var dsp1A, dsp1B, dsp2A, dsp2B Animation

	mirrorAcrossScreens := rand.Intn(2)
	generateBoxiScreens := func() (Animation, Animation) {
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
				return firstAnimation, secondAnimation
			}
		}

		if mirrorAcrossScreens != 0 {
			return firstAnimation, firstAnimation
		}

		randomIndex = rand.Intn(len(validIndices))
		secondAnimation := animationManager.animations[validIndices[randomIndex]]
		return firstAnimation, secondAnimation
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
	blinkSpeed := 0
	if doDaBounce == 6 && baseMood == Party {
		blinkSpeed = defaultBlinkSpeed
	}

	character := Calm
	if baseMood == Regular || baseMood == Party {
		character = Rhythmic
	}
}

func (context *AutoModeContext) getNextLighting(switchType switchType) LightingInstruction {

}
