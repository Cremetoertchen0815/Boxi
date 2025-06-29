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
	IsNsfw             bool
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

func (manager *AnimationManager) ImportAnimation(animationPath string, name string, mood LightingMood, splitAnimation bool, nsfw bool) (Display.AnimationId, error) {
	animationId := rand.Uint32()

	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	if !splitAnimation {
		err := Infrastructure.ExtractFrames(animationId, animationPath)
		if err != nil {
			return 0, err
		}

		animation := Animation{Display.AnimationId(animationId), name, mood, nsfw, Display.None}
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
	animation := Animation{Display.AnimationId(animationId), name, mood, nsfw, rightAnimationId}
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
		{Display.AnimationId(446948159), "Nerd Pacman", Regular, false, Display.None},
		{Display.AnimationId(649833014), "Gottloser Creme", Happy, false, Display.None},
		{Display.AnimationId(678928891), "DVD Logo", Happy, false, Display.None},
		{Display.AnimationId(724152790), "Foxi Jumpscare", Regular, false, Display.None},
		{Display.AnimationId(746302169), "Vaporwave Autobahn", Moody, false, Display.None},
		{Display.AnimationId(899960868), "Saul 3D", Regular, false, Display.None},
		{Display.AnimationId(1204539747), "Nyan Cat", Regular, false, Display.AnimationId(2454484289)},
		{Display.AnimationId(1345034356), "Ash Pat", Happy, false, Display.None},
		{Display.AnimationId(1884833779), "Gopnik", Party, false, Display.None},
		{Display.AnimationId(1899868680), "Cat Bounce", Regular, false, Display.None},
		{Display.AnimationId(1965415769), "Kermit Suizid", Happy, false, Display.None},
		{Display.AnimationId(2243405019), "Aksel.", Happy, false, Display.None},
		{Display.AnimationId(2456904767), "Burning Piano man", Moody, false, Display.None},
		{Display.AnimationId(2500737094), "Doggo dance", Party, false, Display.None},
		{Display.AnimationId(2574938612), "Pedro", Party, false, Display.None},
		{Display.AnimationId(2759311642), "Caramelldansen", Party, false, Display.None},
		{Display.AnimationId(2899126749), "Monke", Moody, false, Display.None},
		{Display.AnimationId(2939821731), "Ribbons", Moody, false, Display.None},
		{Display.AnimationId(3343111115), "Spinning Fish", Regular, false, Display.None},
		{Display.AnimationId(3424648902), "Spinning Neuer", Regular, false, Display.None},
		{Display.AnimationId(3703776356), "Another doggo dancing", Regular, false, Display.None},
	}
}
