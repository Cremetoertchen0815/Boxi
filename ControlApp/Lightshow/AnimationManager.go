package Lightshow

import (
	"ControlApp/Display"
	"ControlApp/Infrastructure"
	"math/rand"
	"sync"
)

type Animation struct {
	Id                 Display.AnimationId
	Mood               LightingMood
	SecondaryAnimation *Display.AnimationId
}

type AnimationManager struct {
	animations  []Animation
	accessLock  *sync.Mutex
	UploadQueue chan Display.AnimationId
}

func LoadAnimations() *AnimationManager {
	uploadQueue := make(chan Display.AnimationId, 2)
	return &AnimationManager{UploadQueue: uploadQueue}
}

func (manager *AnimationManager) ImportAnimation(animationPath string, mood LightingMood, splitAnimation bool) (Display.AnimationId, error) {
	animationId := rand.Uint32()

	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	if !splitAnimation {
		err := Infrastructure.ExtractFrames(animationId, animationPath)
		if err != nil {
			return 0, err
		}

		animation := Animation{Display.AnimationId(animationId), mood, nil}
		manager.animations = append(manager.animations, animation)
		manager.storeDatabase()
		manager.UploadQueue <- animation.Id
		return animation.Id, nil
	}

	secondaryAnimationId := rand.Uint32()

	err := Infrastructure.ExtractDoubleFrames(animationId, secondaryAnimationId, animationPath)
	if err != nil {
		return 0, err
	}

	rightAnimationId := Display.AnimationId(secondaryAnimationId)
	animation := Animation{Display.AnimationId(animationId), mood, &rightAnimationId}
	manager.animations = append(manager.animations, animation)
	manager.storeDatabase()
	manager.UploadQueue <- animation.Id
	manager.UploadQueue <- rightAnimationId
	return animation.Id, nil
}

func (manager *AnimationManager) storeDatabase() {
	// TODO: store the data in the database
}
