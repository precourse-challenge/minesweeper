package matchmaking

import (
	"fmt"
	"minesweeper-infrastructure/network"
	"minesweeper-server/game/room"
	"sync"
)

type Matchmaker struct {
	rooms      map[string]*room.Room
	roomMutex  sync.Mutex
	nextRoomId int
}

func NewMatchmaker() *Matchmaker {
	return &Matchmaker{
		rooms:      make(map[string]*room.Room),
		nextRoomId: 1,
	}
}

func (m *Matchmaker) FindOrCreateRoom(conn *network.Connection) (*room.Room, int) {
	m.roomMutex.Lock()
	defer m.roomMutex.Unlock()

	for _, room := range m.rooms {
		if !room.IsFull() {
			if room.AddPlayer(conn) {
				return room, 2
			}
		}
	}
	roomId := m.generateRoomId()
	m.nextRoomId++

	gameRoom := room.NewRoom(roomId, conn)
	m.rooms[roomId] = gameRoom

	return gameRoom, 1
}

func (m *Matchmaker) RemoveRoom(roomId string) {
	m.roomMutex.Lock()
	defer m.roomMutex.Unlock()
	delete(m.rooms, roomId)
}

func (m *Matchmaker) generateRoomId() string {
	return fmt.Sprintf("room-%d", m.nextRoomId)
}
