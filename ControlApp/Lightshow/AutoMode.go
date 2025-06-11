package Lightshow

import (
	"github.com/stianeikeland/go-rpio/v4"
	"math/rand"
	"time"
)

type VisualSwitch interface {
	applyLighting(instruction LightingInstruction)
	applyAnimation(instruction AnimationsInstruction)
	triggerBeat()
	getAnimations() *AnimationManager
	getPalettes() *PaletteManager
}

type AutoModeContext struct {
	Configuration         AutoModeConfiguration
	switcher              VisualSwitch
	beatInputPin          rpio.Pin
	lastBeat              *time.Time
	lightingSwitchToCalm  *time.Time
	animationSwitchToCalm *time.Time
	lightingDeadTime      *time.Duration
	animationDeadTime     *time.Duration
	lightingBeatsLeft     int
	animationBeatsLeft    int
	wasInCalmMode         bool
}

const soundInputPin = 16
const loopDelayMs = 5

func CreateAutoMode(switcher VisualSwitch, configuration AutoModeConfiguration) *AutoModeContext {
	pin := rpio.Pin(soundInputPin)
	pin.Input()

	result := &AutoModeContext{Configuration: configuration, switcher: switcher, beatInputPin: pin}

	go result.calculateAutoMode()
	return result
}

func (context *AutoModeContext) calculateAutoMode() {

	for {
		time.Sleep(loopDelayMs * time.Millisecond)

		//Only count beat if we're not in an exclusively calm mood
		pendingBeat := context.beatInputPin.Read() == rpio.High && context.Configuration.Mood.IsCalm()
		now := time.Now()
		context.lightingSwitchToCalm = nil
		context.animationSwitchToCalm = nil

		var lastBeat time.Time
		if context.lastBeat == nil {
			lastBeat = now
		} else {
			lastBeat = *context.lastBeat
		}

		isBeat := pendingBeat && (context.lastBeat == nil || now.After(context.lastBeat.Add(context.Configuration.MinTimeBetweenBeats)))
		if isBeat {
			context.lastBeat = &now
			context.switcher.triggerBeat()

			// Count down the display beat timer and play new animation if limit was reached
			context.animationBeatsLeft--
			if context.animationBeatsLeft <= 0 {
				animation := context.getNextAnimation(OnBeat)
				context.switcher.applyAnimation(animation)

				timingConstraint, ok := context.Configuration.AnimationModeTiming[animation.character]
				if ok {
					context.animationBeatsLeft = getNextBeatConstraint(timingConstraint)
					context.animationDeadTime = &timingConstraint.NoBeatDeadTime
				}
			}

			// Count down the animation beat timer and play new animation if limit was reached
			context.lightingBeatsLeft--
			if context.lightingBeatsLeft <= 0 {
				var lighting LightingInstruction
				if context.wasInCalmMode {
					lighting = context.getNextLighting(FirstBeat)
				} else {
					lighting = context.getNextLighting(OnBeat)
				}

				context.switcher.applyLighting(lighting)

				timingConstraint, ok := context.Configuration.LightingModeTiming[lighting.character]
				if ok {
					context.animationBeatsLeft = getNextBeatConstraint(timingConstraint)
					context.lightingDeadTime = &timingConstraint.NoBeatDeadTime
				}
			}

			return
		}

		//Check beat dead time for animation
		if context.animationDeadTime != nil && lastBeat.Add(*context.animationDeadTime).Before(time.Now()) {
			context.animationDeadTime = nil
			animation := context.getNextAnimation(InDeadTime)
			context.switcher.applyAnimation(animation)

			if animation.character == Calm {
				timeWhenBoring := time.Now().Add(context.Configuration.AnimationCalmModeBoring)
				context.animationSwitchToCalm = &timeWhenBoring
			}
		}

		//Check beat dead time for lighting
		if context.lightingDeadTime != nil && lastBeat.Add(*context.lightingDeadTime).Before(time.Now()) {
			context.lightingDeadTime = nil
			lighting := context.getNextLighting(InDeadTime)
			context.switcher.applyLighting(lighting)
			context.wasInCalmMode = true

			if lighting.character == Calm {
				timeWhenBoring := time.Now().Add(context.Configuration.LightingCalmModeBoring)
				context.lightingSwitchToCalm = &timeWhenBoring
			}
		}

		//Check if the calm animation is boring
		if context.animationSwitchToCalm != nil && context.animationSwitchToCalm.Before(time.Now()) {
			context.animationSwitchToCalm = nil
			animation := context.getNextAnimation(InCalmMode)
			context.switcher.applyAnimation(animation)

			if animation.character == Calm {
				timeWhenBoring := time.Now().Add(context.Configuration.AnimationCalmModeBoring)
				context.animationSwitchToCalm = &timeWhenBoring
			}
		}

		//Check if the calm lighting is boring
		if context.lightingSwitchToCalm != nil && context.animationSwitchToCalm.Before(time.Now()) {
			context.lightingSwitchToCalm = nil
			lighting := context.getNextLighting(InCalmMode)
			context.switcher.applyLighting(lighting)
			context.wasInCalmMode = true

			if lighting.character == Calm {
				timeWhenBoring := time.Now().Add(context.Configuration.LightingCalmModeBoring)
				context.lightingSwitchToCalm = &timeWhenBoring
			}
		}
	}
}

func getNextBeatConstraint(constraint TimingConstraint) int {
	return rand.Intn(constraint.MaxNumberOfBeats-constraint.MinNumberOfBeats) + constraint.MinNumberOfBeats
}
