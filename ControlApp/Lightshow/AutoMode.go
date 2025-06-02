package Lightshow

import (
	"ControlApp/BoxiBus"
	"ControlApp/Display"
	"github.com/stianeikeland/go-rpio/v4"
	"time"
)

type VisualSwitch interface {
	applyLighting(block BoxiBus.MessageBlock)
	applyAnimation(animation Display.AnimationId, displays []Display.ServerDisplay)
}

type AutoModeContext struct {
	Configuration     AutoModeConfiguration
	switcher          VisualSwitch
	beatInputPin      rpio.Pin
	lastBeat          *time.Time
	slowMode          bool
	numberOfBeatsLeft int
	currentDeadTime   *time.Time
}

type AutoModeConfiguration struct {
	Mood                     LightingMood
	MinTimeBetweenBeats      time.Duration
	NormalBeatModeTiming     TimingConstraint //The timing constraints of all regular beat-based modes
	StrobeModeTiming         TimingConstraint //The timing constraints of the strobe mode
	EnergeticAnimationTiming TimingConstraint //The timing constraints of an energetic animation
}

type TimingConstraint struct {
	MinNumberOfBeats int           //The least number of beats before switching to the next mode.
	MaxNumberOfBeats int           //The most number of beats before switching to the next mode.
	DeadTime         time.Duration //The duration since the last beat when forcibly switching to a calm mode.
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

	result := &AutoModeContext{configuration, switcher, pin, nil, true}

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
		}
	}
}
