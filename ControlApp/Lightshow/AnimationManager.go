package Lightshow

import (
	"ControlApp/Display"
	"ControlApp/Infrastructure"
	"math/rand"
	"sync"
)

type Animation struct {
	Id                 Display.AnimationId
	Name               string
	Mood               LightingMood
	SecondaryAnimation Display.AnimationId
}

type AnimationManager struct {
	animations  []Animation
	accessLock  *sync.Mutex
	UploadQueue chan Display.AnimationId
}

func LoadAnimations() *AnimationManager {
	uploadQueue := make(chan Display.AnimationId, 2)
	return &AnimationManager{
		animations:  getDefaultAnimations(),
		UploadQueue: uploadQueue,
		accessLock:  &sync.Mutex{},
	}
}

func (manager *AnimationManager) ImportAnimation(animationPath string, name string, mood LightingMood, splitAnimation bool) (Display.AnimationId, error) {
	animationId := rand.Uint32()

	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	if !splitAnimation {
		err := Infrastructure.ExtractFrames(animationId, animationPath)
		if err != nil {
			return 0, err
		}

		animation := Animation{Display.AnimationId(animationId), name, mood, Display.None}
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
	animation := Animation{Display.AnimationId(animationId), name, mood, rightAnimationId}
	manager.animations = append(manager.animations, animation)
	manager.storeDatabase()
	manager.UploadQueue <- animation.Id
	manager.UploadQueue <- rightAnimationId
	return animation.Id, nil
}

func (manager *AnimationManager) GetById(id Display.AnimationId) (bool, Animation) {
	for _, animation := range manager.animations {
		if animation.Id == id {
			return true, animation
		}
	}

	return false, Animation{}
}

func (manager *AnimationManager) storeDatabase() {
	// TODO: store the data in the database
}

func getDefaultAnimations() []Animation {
	return []Animation{
		{Display.AnimationId(446948159), "Nerd Pacman", Regular, Display.None},
		{Display.AnimationId(649833014), "Gottloser Creme", Happy, Display.None},
		{Display.AnimationId(678928891), "DVD Logo", Happy, Display.None},
		{Display.AnimationId(724152790), "Foxi Jumpscare", Regular, Display.None},
		{Display.AnimationId(746302169), "Vaporwave Autobahn", Moody, Display.None},
		{Display.AnimationId(899960868), "Saul 3D", Regular, Display.None},
		{Display.AnimationId(1204539747), "Nyan Cat", Regular, Display.AnimationId(2454484289)},
		{Display.AnimationId(1345034356), "Ash Pat", Happy, Display.None},
		{Display.AnimationId(1884833779), "Gopnik", Party, Display.None},
		{Display.AnimationId(1899868680), "Cat Bounce", Regular, Display.None},
		{Display.AnimationId(1965415769), "Kermit Suizid", Happy, Display.None},
		{Display.AnimationId(2243405019), "Aksel.", Happy, Display.None},
		{Display.AnimationId(2456904767), "Burning Piano man", Moody, Display.None},
		{Display.AnimationId(2500737094), "Doggo dance", Party, Display.None},
		{Display.AnimationId(2574938612), "Pedro", Party, Display.None},
		{Display.AnimationId(2759311642), "Caramelldansen", Party, Display.None},
		{Display.AnimationId(2899126749), "Monke", Moody, Display.None},
		{Display.AnimationId(2939821731), "Ribbons", Moody, Display.None},
		{Display.AnimationId(3343111115), "Spinning Fish", Regular, Display.None},
		{Display.AnimationId(3424648902), "Spinning Neuer", Regular, Display.None},
		{Display.AnimationId(3703776356), "Another doggo dancing", Regular, Display.None},
	}
}
