package Lightshow

import (
	"math/rand"
	"time"
)

type Manager interface {
	applyLighting(instruction LightingInstruction)
	applyAnimation(instruction AnimationsInstruction)
	triggerBeat()
	getBeatState() bool
	GetAnimations() *AnimationManager
	GetPalettes() *PaletteManager
}

type AutoModeContext struct {
	Configuration         AutoModeConfiguration
	manager               Manager
	lastBeat              *time.Time
	lightingSwitchToCalm  *time.Time
	animationSwitchToCalm *time.Time
	lightingDeadTime      *time.Duration
	animationDeadTime     *time.Duration
	lightingBeatsLeft     int
	animationBeatsLeft    int
	wasInCalmMode         bool
}

const loopDelayMs = 5

func CreateAutoMode(switcher Manager, configuration AutoModeConfiguration) *AutoModeContext {
	lightingSwitchTime := time.Now().Add(configuration.LightingCalmModeBoring)
	animationSwitchTime := time.Now().Add(configuration.AnimationCalmModeBoring)
	result := &AutoModeContext{
		Configuration:         configuration,
		manager:               switcher,
		lightingSwitchToCalm:  &lightingSwitchTime,
		animationSwitchToCalm: &animationSwitchTime,
	}

	go result.calculateAutoMode()
	return result
}

func (context *AutoModeContext) calculateAutoMode() {

	for {
		time.Sleep(loopDelayMs * time.Millisecond)

		//Only count beat if we're not in an exclusively calm mood
		pendingBeat := context.manager.getBeatState() && !context.Configuration.Mood.IsCalm()
		now := time.Now()

		var lastBeat time.Time
		if context.lastBeat == nil {
			lastBeat = now
		} else {
			lastBeat = *context.lastBeat
		}

		isBeat := pendingBeat && (context.lastBeat == nil || now.After(context.lastBeat.Add(context.Configuration.MinTimeBetweenBeats)))
		if isBeat {
			context.lastBeat = &now
			context.manager.triggerBeat()
			context.lightingSwitchToCalm = nil
			context.animationSwitchToCalm = nil

			// Count down the display beat timer and play new animation if limit was reached
			context.animationBeatsLeft--
			if context.animationBeatsLeft <= 0 {
				animation := context.getNextAnimation(OnBeat)
				context.manager.applyAnimation(animation)

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

				context.manager.applyLighting(lighting)

				timingConstraint, ok := context.Configuration.LightingModeTiming[lighting.character]
				if ok {
					context.lightingBeatsLeft = getNextBeatConstraint(timingConstraint)
					context.lightingDeadTime = &timingConstraint.NoBeatDeadTime
				}
			}

			continue
		}

		//Check beat dead time for animation
		if context.animationDeadTime != nil && lastBeat.Add(*context.animationDeadTime).Before(time.Now()) {
			context.animationDeadTime = nil
			context.animationBeatsLeft = 0
			animation := context.getNextAnimation(InDeadTime)
			context.manager.applyAnimation(animation)

			if animation.character == Calm {
				timeWhenBoring := time.Now().Add(context.Configuration.AnimationCalmModeBoring)
				context.animationSwitchToCalm = &timeWhenBoring
			}
		}

		//Check beat dead time for lighting
		if context.lightingDeadTime != nil && lastBeat.Add(*context.lightingDeadTime).Before(time.Now()) {
			context.lightingDeadTime = nil
			context.lightingBeatsLeft = 0
			lighting := context.getNextLighting(InDeadTime)
			context.manager.applyLighting(lighting)
			context.wasInCalmMode = lighting.character == Calm

			if context.wasInCalmMode {
				timeWhenBoring := time.Now().Add(context.Configuration.LightingCalmModeBoring)
				context.lightingSwitchToCalm = &timeWhenBoring
			}
		}

		//Check if the calm animation is boring
		if context.animationSwitchToCalm != nil && context.animationSwitchToCalm.Before(time.Now()) {
			context.animationSwitchToCalm = nil
			animation := context.getNextAnimation(InCalmMode)
			context.manager.applyAnimation(animation)

			if animation.character == Calm {
				timeWhenBoring := time.Now().Add(context.Configuration.AnimationCalmModeBoring)
				context.animationSwitchToCalm = &timeWhenBoring
			}
		}

		//Check if the calm lighting is boring
		if context.lightingSwitchToCalm != nil && context.lightingSwitchToCalm.Before(time.Now()) {
			context.lightingSwitchToCalm = nil
			lighting := context.getNextLighting(InCalmMode)
			context.manager.applyLighting(lighting)
			context.wasInCalmMode = true

			if lighting.character == Calm {
				timeWhenBoring := time.Now().Add(context.Configuration.LightingCalmModeBoring)
				context.lightingSwitchToCalm = &timeWhenBoring
			}
		}

		if isBeat {
			context.wasInCalmMode = false
		}
	}
}

func getNextBeatConstraint(constraint TimingConstraint) int {
	return rand.Intn(constraint.MaxNumberOfBeats-constraint.MinNumberOfBeats) + constraint.MinNumberOfBeats
}
