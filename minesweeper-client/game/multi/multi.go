package multi

import (
	"fmt"
	"log"
	"minesweeper-client/game/view"
)

type MultiMode struct {
	session *Session
}

func NewMultiMode() *MultiMode {
	client, err := NewSession("127.0.0.1:8080")
	if err != nil {
		log.Fatal("서버 연결에 실패했습니다.", err)
	}
	return &MultiMode{session: client}
}

func (m *MultiMode) Start() {
	defer m.closeConnection(m.session)

	err := m.session.JoinGame()
	if err != nil {
		log.Fatal("게임 참가에 실패했습니다.", err)
		return
	}

	view.ShowOpponentWaitMessage()

	go m.session.StartReceiving()
}

func (m *MultiMode) closeConnection(gameClient *Session) {
	err := gameClient.Close()
	if err != nil {
		err := fmt.Errorf("클라이언트 연결 종료를 실패했습니다 %w", err)
		view.ShowErrorMessage(err)
	}
}
