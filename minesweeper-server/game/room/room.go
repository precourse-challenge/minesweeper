package room

import (
	"fmt"
	"log"
	"minesweeper-infrastructure/dto"
	"minesweeper-infrastructure/network"
	"minesweeper-infrastructure/protocol"
	"minesweeper-server/game/match"
	"sync"
)

type Room struct {
	id          string
	players     [2]*network.Connection
	playerCount int
	match       *match.Match
	mutex       sync.Mutex
}

func NewRoom(id string, firstPlayer *network.Connection) *Room {
	return &Room{
		id:          id,
		players:     [2]*network.Connection{firstPlayer, nil},
		playerCount: 1,
		match:       match.NewMatch(),
	}
}

func (r *Room) AddPlayer(player *network.Connection) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.playerCount >= 2 {
		return false
	}

	r.players[1] = player
	r.playerCount = 2
	return true
}

func (r *Room) IsFull() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.playerCount == 2
}

func (r *Room) StartGame() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.match.InitializeGame()

	board1 := dto.ToBoardDto(r.match.GetPlayer1Board())
	board2 := dto.ToBoardDto(r.match.GetPlayer2Board())

	message := protocol.Message{
		Type:    protocol.Start,
		Board1:  board1,
		Board2:  board2,
		Message: "게임이 시작되었습니다!",
	}
	log.Println(message.Message)

	r.broadcastMessage(message)
}

func (r *Room) HandleDisconnect(playerId int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	winner := 3 - playerId
	message := protocol.Message{
		Type:    protocol.GameOver,
		Winner:  winner,
		Message: fmt.Sprintf("플레이어 %d가 연결을 끊었습니다. 플레이어 %d 승리!", playerId, winner),
	}

	otherPlayerIndex := winner - 1
	if r.players[otherPlayerIndex] != nil {
		err := r.players[otherPlayerIndex].Send(message)
		if err != nil {
			log.Println("상대 플레이어에게 종료 메시지 전송을 실패했습니다.", err)
		}
	}

	r.cleanup()
}

func (r *Room) GetId() string {
	return r.id
}

func (r *Room) broadcastMessage(message protocol.Message) {
	for i, player := range r.players {
		if player != nil {
			messageCopy := message
			messageCopy.PlayerId = i + 1
			err := player.Send(messageCopy)
			if err != nil {
				fmt.Println("메시지 전송을 실패했습니다.", err)
			}
		}
	}
}

func (r *Room) cleanup() {
	for _, player := range r.players {
		if player != nil {
			err := player.Close()
			if err != nil {
				log.Println("플레이어 연결 종료를 실패했습니다.", err)
			}
		}
	}
}
