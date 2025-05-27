package Display

import (
	"encoding/binary"
)

func (manager *ServerManager) UploadAnimation(animationId uint32, frames []string, displays ServerDisplay) {
	servers := make(map[byte]*Server)

	//Find all display servers
	for display := range displays {
		server, ok := manager.connections[byte(display)/2]

		//If specified display has no registered server, skip display
		if !ok {
			continue
		}

		servers[server.identifier] = server
	}

	//Check, which display servers still need the animations
	var serversToBeUpdates []byte
	animationIdBytes := binary.BigEndian.AppendUint32([]byte{}, animationId)
	frameCount := uint16(len(frames))
	for _, server := range servers {
		exists, err := server.sendInstructionWithCallback(DoesAnimationExist, frameCount, animationIdBytes)
		if err == nil && !exists {
			serversToBeUpdates = append(serversToBeUpdates, server.identifier)
		}
	}

	//Send frames
}

func (manager *ServerManager) PlayAnimation(animationId uint32, displays ServerDisplay) {
	payload := binary.BigEndian.AppendUint32([]byte{}, animationId)

	for serverId, server := range manager.connections {
		//If the current server isn't linked to any display that should show, continue.
		if displays&(allLocalDisplays<<(serverId*2)) != 0 {
			continue
		}

		localDisplays := (displays >> (serverId * 2)) & allLocalDisplays
		server.sendInstructionWithoutCallback(PlayAnimation, uint16(localDisplays), payload)
	}
}

func (manager *ServerManager) DisplayText(textToDisplay string, displays ServerDisplay) {
	payload := []byte(textToDisplay)

	for serverId, server := range manager.connections {
		//If the current server isn't linked to any display that should show, continue.
		if displays&(allLocalDisplays<<(serverId*2)) != 0 {
			continue
		}

		localDisplays := (displays >> (serverId * 2)) & allLocalDisplays
		server.sendInstructionWithoutCallback(PlayAnimation, uint16(localDisplays), payload)
	}
}
