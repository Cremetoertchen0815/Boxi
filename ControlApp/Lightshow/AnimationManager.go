package Lightshow

import (
	"ControlApp/Display"
	"ControlApp/Infrastructure"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
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
	animations  map[Display.AnimationId]Animation
	accessLock  *sync.Mutex
	UploadQueue chan Display.AnimationId
}

const animationsConfigPath = "Configuration/animations.json"

func LoadAnimations() *AnimationManager {
	configFile, err := os.Open(animationsConfigPath)

	var config map[Display.AnimationId]Animation
	if err != nil {
		log.Fatalf("Config file for animations could not be accessed! %s", err)
	}

	defer func(configFile *os.File) {
		_ = configFile.Close()
	}(configFile)

	jsonParser := json.NewDecoder(configFile)

	err = jsonParser.Decode(&config)
	if err != nil {
		log.Fatalf("Invalid JSON format of animations config file! %s", err)
	}

	uploadQueue := make(chan Display.AnimationId, 2)
	return &AnimationManager{
		animations:  config,
		UploadQueue: uploadQueue,
		accessLock:  &sync.Mutex{},
	}
}

func (manager *AnimationManager) storeConfiguration() {
	configFile, err := os.OpenFile(animationsConfigPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)

	if err != nil {
		log.Fatalf("Config file for animations could not be opened for writing! %s", err)
	}

	defer func(configFile *os.File) {
		_ = configFile.Close()
	}(configFile)

	jsonParser := json.NewEncoder(configFile)
	err = jsonParser.Encode(manager.animations)
	if err != nil {
		log.Fatalf("Configuration for animations could be JSON encoded! %s", err)
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
		manager.animations[animation.Id] = animation
		manager.storeConfiguration()
		manager.UploadQueue <- animation.Id
		storeThumbnail(uint32(animation.Id))
		return animation.Id, nil
	}

	secondaryAnimationId := rand.Uint32()

	err := Infrastructure.ExtractDoubleFrames(animationId, secondaryAnimationId, animationPath)
	if err != nil {
		return 0, err
	}

	rightAnimationId := Display.AnimationId(secondaryAnimationId)
	animation := Animation{Display.AnimationId(animationId), name, mood, nsfw, rightAnimationId}
	manager.animations[animation.Id] = animation
	manager.storeConfiguration()
	manager.UploadQueue <- animation.Id
	manager.UploadQueue <- rightAnimationId
	storeThumbnail(uint32(animation.Id))
	return animation.Id, nil
}

func storeThumbnail(id uint32) {
	sourcePath := fmt.Sprintf("animations/%d/0001.png", id)
	destinationPath := fmt.Sprintf("Frontend/static/thumbs/%d.png", id)

	//Make sure source does exist
	if _, err := os.Stat(sourcePath); err != nil {
		return
	}

	//Make sure destination doesn't exist
	if _, err := os.Stat(destinationPath); err == nil {
		return
	}

	_ = os.Link(sourcePath, destinationPath)
}

func (manager *AnimationManager) GetById(id Display.AnimationId) (bool, Animation) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	item, success := manager.animations[id]
	return success, item
}

func (manager *AnimationManager) GetAll() []Animation {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	var animations []Animation
	for _, animation := range manager.animations {
		animations = append(animations, animation)
	}
	return animations
}

func (manager *AnimationManager) RemoveAnimation(animationId Display.AnimationId) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	delete(manager.animations, animationId)
	manager.storeConfiguration()
}
