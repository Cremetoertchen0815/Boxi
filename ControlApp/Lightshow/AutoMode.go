package Lightshow

import (
	"ControlApp/BoxiBus"
	"ControlApp/Infrastructure"
	"github.com/stianeikeland/go-rpio/v4"
	"math/rand"
	"time"
)

type VisualSwitch interface {
	applyLighting(instruction Infrastructure.LightingInstruction)
	applyAnimation(instruction Infrastructure.AnimationInstruction)
	triggerBeat()
}

type AutoModeContext struct {
	Configuration      AutoModeConfiguration
	switcher           VisualSwitch
	beatInputPin       rpio.Pin
	lastBeat           *time.Time
	lightingDeadTime   *time.Duration
	animationDeadTime  *time.Duration
	lightingBeatsLeft  int
	animationBeatsLeft int
}

type AutoModeConfiguration struct {
	Mood                 LightingMood
	MinTimeBetweenBeats  time.Duration
	NormalBeatModeTiming TimingConstraint //The timing constraints of all regular beat-based modes
	StrobeModeTiming     TimingConstraint //The timing constraints of the strobe mode
	AnimationTiming      TimingConstraint //The timing constraints of an energetic animation
	SlowModeDeadTime     time.Duration    //The duration since the last mode change in slow mode when forcibly switching to another mode.
}

type TimingConstraint struct {
	MinNumberOfBeats int           //The least number of beats before switching to the next mode.
	MaxNumberOfBeats int           //The most number of beats before switching to the next mode.
	NoBeatDeadTime   time.Duration //The duration since the last beat when forcibly switching to a calm mode.
}

type LightingMood uint8

const (
	Chill LightingMood = iota
	Moody
	Regular
	Party
)

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

		pendingBeat := context.beatInputPin.Read() == rpio.High
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
			context.switcher.triggerBeat()

			// Count down the display beat timer and play new animation if limit was reached
			context.animationBeatsLeft--
			if context.animationBeatsLeft <= 0 {
				for _, animation := range context.getNextAnimation(OnBeat) {
					context.switcher.applyAnimation(animation)
				}

				context.animationBeatsLeft = getNextBeatConstraint(context.Configuration.AnimationTiming)
				context.animationDeadTime = &context.Configuration.AnimationTiming.NoBeatDeadTime
			}

			// Count down the animation beat timer and play new animation if limit was reached
			context.lightingBeatsLeft--
			if context.lightingBeatsLeft <= 0 {
				lighting := context.getNextLighting(OnBeat)
				context.switcher.applyLighting(lighting)

				timingConstraint := context.Configuration.NormalBeatModeTiming
				if lighting.Mode == BoxiBus.Strobe {
					timingConstraint = context.Configuration.StrobeModeTiming
				}

				context.animationBeatsLeft = getNextBeatConstraint(timingConstraint)
				context.lightingDeadTime = &timingConstraint.NoBeatDeadTime
			}

			return
		}

		//Check beat dead time for animation
		if context.animationDeadTime != nil && lastBeat.Add(*context.animationDeadTime).Before(time.Now()) {
			context.animationDeadTime = nil

			for _, animation := range context.getNextAnimation(InDeadTime) {
				context.switcher.applyAnimation(animation)
			}
		}
	}
}

func getNextBeatConstraint(constraint TimingConstraint) int {
	return rand.Intn(constraint.MaxNumberOfBeats-constraint.MinNumberOfBeats) + constraint.MinNumberOfBeats
}
