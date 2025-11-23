package room

import (
	"fmt"
	"log"
	"minesweeper-core/position"
	"minesweeper-infrastructure/dto"
	"minesweeper-infrastructure/network"
	"minesweeper-infrastructure/protocol"
	"minesweeper-server/game/match"
	"sync"
)

type Room struct {
	id          string
	players     [2]*network.Connection
	playerIds   map[*network.Connection]int
	playerCount int
	match       *match.Match
	mutex       sync.Mutex
}

func NewRoom(id string, firstPlayer *network.Connection) *Room {
	return &Room{
		id:          id,
		players:     [2]*network.Connection{firstPlayer, nil},
		playerIds:   map[*network.Connection]int{firstPlayer: 1},
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
	if r.playerIds == nil {
		r.playerIds = make(map[*network.Connection]int)
	}
	r.playerIds[player] = 2
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

func (r *Room) HandleOpen(conn *network.Connection, row, col int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	playerId, err := r.getValidPlayerId(conn)
	if err != nil {
		log.Printf("[Room %s][Open] 유효하지 않은 플레이어 요청 - %v\n", r.id, err)
		return
	}

	cellPosition, err := position.NewCellPosition(row, col)
	if err != nil {
		r.sendError(conn, err.Error())
		log.Printf("[Room %s][Open] 잘못된 좌표 입력 - player %d, row=%d col=%d\n", r.id, playerId, row, col)
		return
	}

	result, err := r.match.Open(playerId, cellPosition)
	if err != nil {
		r.sendError(conn, err.Error())
		log.Printf("[Room %s][Open] Open 실패 - player %d,  %v\n", r.id, playerId, err)
		return
	}

	board1 := dto.ToBoardDto(r.match.GetPlayer1Board())
	board2 := dto.ToBoardDto(r.match.GetPlayer2Board())

	if result.IsGameOver {
		message := protocol.Message{
			Type:    protocol.GameOver,
			Board1:  board1,
			Board2:  board2,
			Winner:  result.Winner,
			Message: fmt.Sprintf("플레이어 %d가 %s 플레이어 %d 승리!", playerId, result.Message, result.Winner),
		}
		log.Printf("[Room %s][Open] 게임 종료 %s\n", r.id, result.Message)
		r.broadcastMessage(message)
	} else {
		message := protocol.Message{
			Type:   protocol.Update,
			Board1: board1,
			Board2: board2,
		}
		log.Printf("[Room %s][Open] 상태 업데이트 브로드캐스트 - player %d\n", r.id, playerId)
		r.broadcastMessage(message)
	}
}

func (r *Room) HandleFlag(conn *network.Connection, row, col int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	playerId, err := r.getValidPlayerId(conn)
	if err != nil {
		log.Printf("[Room %s][Flag] 유효하지 않은 플레이어 요청 - %v\n", r.id, err)
		return
	}

	cellPosition, err := position.NewCellPosition(row, col)
	if err != nil {
		r.sendError(conn, err.Error())
		log.Printf("[Room %s][Flag] 잘못된 좌표 입력 - player %d, row=%d col=%d\n", r.id, playerId, row, col)
		return
	}

	err = r.match.Flag(playerId, cellPosition)
	if err != nil {
		r.sendError(conn, err.Error())
		log.Printf("[Room %s][Flag] Flag 실패 - player %d,  %v\n", r.id, playerId, err)
		return
	}

	board1 := dto.ToBoardDto(r.match.GetPlayer1Board())
	board2 := dto.ToBoardDto(r.match.GetPlayer2Board())

	message := protocol.Message{
		Type:   protocol.Update,
		Board1: board1,
		Board2: board2,
	}
	log.Printf("[Room %s][Flag] 상태 업데이트 브로드캐스트 - player %d\n", r.id, playerId)
	r.broadcastMessage(message)
}

func (r *Room) HandleDisconnect(conn *network.Connection) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	playerId, err := r.getValidPlayerId(conn)
	if err != nil {
		log.Printf("[Room %s][Disconnect] 유효하지 않은 플레이어 요청 - %v\n", r.id, err)
		r.cleanup()
		return
	}

	winner := 3 - playerId

	otherPlayerIndex := winner - 1
	if otherPlayerIndex >= 0 && otherPlayerIndex < len(r.players) && r.players[otherPlayerIndex] != nil {
		message := protocol.Message{
			Type:     protocol.GameOver,
			PlayerId: winner,
			Winner:   winner,
			Board1:   dto.BoardDto{},
			Board2:   dto.BoardDto{},
			Message:  fmt.Sprintf("상대방이 연결을 끊었습니다. 당신이 승리했습니다!"),
		}

		log.Printf("[Room %s][Disconnect] player %d 연결 끊김 - winner=%d\n", r.id, playerId, winner)
		err := r.players[otherPlayerIndex].Send(message)
		if err != nil {
			log.Printf("[Room %s][Disconnect] 상대에게 GameOver 메시지 전송 실패 - %v\n", r.id, err)
		}
	}

	r.cleanup()
}

func (r *Room) sendError(conn *network.Connection, errorMessage string) {
	message := protocol.Message{
		Type:    protocol.Error,
		Message: errorMessage,
	}

	err := conn.Send(message)
	if err != nil {
		log.Printf("[Room %s][sendError] 에러 메시지 전송 실패 - %v\n", r.id, err)
	}
}

func (r *Room) getValidPlayerId(conn *network.Connection) (int, error) {
	playerId := r.getPlayerId(conn)
	if playerId == 0 {
		return 0, fmt.Errorf("유효하지 않은 플레이어의 요청입니다")
	}
	return playerId, nil
}

func (r *Room) getPlayerId(conn *network.Connection) int {
	if r.playerIds == nil {
		return 0
	}
	if id, exists := r.playerIds[conn]; exists {
		return id
	}
	return 0
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
				log.Printf("[Room %s][broadcast] 플레이어 %d에게 메시지 전송 실패 - %v\n", r.id, i+1, err)
			}
		}
	}
}

func (r *Room) cleanup() {
	for i, player := range r.players {
		if player != nil {
			err := player.Close()
			if err != nil {
				log.Printf("[Room %s][cleanup] 플레이어 %d 연결 종료 실패 - %v\n", r.id, i+1, err)
			} else {
				log.Printf("[Room %s][cleanup] 플레이어 %d 연결 종료\n", r.id, i+1)
			}
		}
	}
}
