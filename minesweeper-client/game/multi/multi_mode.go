package multi

import (
	"fmt"
	"log"
	"minesweeper-client/game/user"
	"minesweeper-client/game/view"
	"minesweeper-core/position"
	"strconv"
	"strings"
)

type MultiMode struct {
	session              *Session
	sessionEventChannels *SessionEventChannels
	inputReady           chan struct{}
	done                 chan struct{}
}

func NewMultiMode() *MultiMode {
	eventChannels := NewSessionEventChannels()

	session, err := NewSession("127.0.0.1:8080", eventChannels)
	if err != nil {
		log.Fatal("서버 연결 실패:", err)
	}

	return &MultiMode{
		session:              session,
		sessionEventChannels: eventChannels,
		inputReady:           make(chan struct{}, 5),
		done:                 make(chan struct{}),
	}
}

func (m *MultiMode) Start() {
	err := m.session.JoinGame()
	if err != nil {
		log.Fatal("게임 참가에 실패했습니다.", err)
		return
	}

	go m.session.StartReceiving()
	go m.handleSessionEvents()

	m.runInputLoop()
}

func (m *MultiMode) handleSessionEvents() {
	defer close(m.done)

	for {
		select {

		case e := <-m.sessionEventChannels.JoinedChan:
			view.ShowPlayerJoined(e.PlayerId)
			view.ShowOpponentWaitMessage()

		case e := <-m.sessionEventChannels.StartChan:
			fmt.Println(e.Message)
			view.ShowMultiBoards(e.Board1, e.Board2, e.PlayerId)
			m.signalInputReady()

		case e := <-m.sessionEventChannels.UpdateChan:
			view.ShowMultiBoards(e.Board1, e.Board2, e.PlayerId)
			m.signalInputReady()

		case e := <-m.sessionEventChannels.ErrorChan:
			fmt.Println(e.Message)
			view.ShowErrorMessage(e.Err)
			m.signalInputReady()

		case e := <-m.sessionEventChannels.GameOverChan:
			fmt.Println(e.Message)
			m.displayGameOver(e)
			m.closeConnection(m.session)
			return
		}
	}
}

func (m *MultiMode) displayGameOver(e GameOverEvent) {
	view.ShowMultiBoards(e.Board1, e.Board2, e.PlayerId)

	if e.Winner == e.PlayerId {
		view.ShowWinMessage()
	} else {
		view.ShowLoseMessage()
	}
}

func (m *MultiMode) runInputLoop() {
	for {
		select {
		case <-m.inputReady:
			exit := m.processUserInput()
			if exit {
				return
			}
		case <-m.done:
			return
		}
	}
}

func (m *MultiMode) processUserInput() bool {
	view.AskCommand()
	action, cellPosition, err := m.readCommand()
	if err != nil {
		view.ShowErrorMessage(err)
		m.signalInputReady()
		return false
	}

	if action == user.Exit {
		err := m.session.Exit()
		if err != nil {
			view.ShowErrorMessage(err)
		}
		view.ShowQuitMessage()
		return true
	}

	err = m.handleActionOnCell(action, cellPosition)
	if err != nil {
		view.ShowErrorMessage(err)
		m.signalInputReady()
	}

	return false
}

func (m *MultiMode) readCommand() (user.Action, *position.CellPosition, error) {
	inputCommand := view.Read()

	action, cellPosition, err := m.parseCommand(inputCommand)
	if err != nil {
		return user.UnknownAction, nil, err
	}
	return action, cellPosition, nil
}

func (m *MultiMode) parseCommand(inputCommand string) (user.Action, *position.CellPosition, error) {
	commands := strings.Fields(inputCommand)

	if len(commands) == 0 {
		return user.UnknownAction, nil, fmt.Errorf("명령어를 입력해주세요")
	}

	inputAction := commands[0]
	action := user.ActionFrom(inputAction)

	if action == user.Exit {
		return user.Exit, nil, nil
	}

	if len(commands) != 3 {
		return user.UnknownAction, nil, fmt.Errorf("명령어 형식이 올바르지 않습니다")
	}

	inputRow := commands[1]
	inputCol := commands[2]

	row, err := m.parseCoordinate(inputRow)
	if err != nil {
		return user.UnknownAction, nil, err
	}

	col, err := m.parseCoordinate(inputCol)
	if err != nil {
		return user.UnknownAction, nil, err
	}

	cellPosition, err := position.NewCellPosition(row, col)
	if err != nil {
		return user.UnknownAction, nil, err
	}

	return action, cellPosition, nil
}

func (m *MultiMode) parseCoordinate(input string) (int, error) {
	coordinate, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("좌표는 숫자여야 합니다")
	}
	return coordinate - 1, nil
}

func (m *MultiMode) handleActionOnCell(action user.Action, pos *position.CellPosition) error {
	switch action {
	case user.Open:
		return m.session.Open(pos.RowIndex(), pos.ColIndex())
	case user.Flag:
		return m.session.Flag(pos.RowIndex(), pos.ColIndex())
	default:
		return fmt.Errorf("잘못된 명령어입니다")
	}
}

func (m *MultiMode) signalInputReady() {
	select {
	case m.inputReady <- struct{}{}:
	default:
	}
}

func (m *MultiMode) closeConnection(gameClient *Session) {
	err := gameClient.Close()
	if err != nil {
		err := fmt.Errorf("연결 종료를 실패했습니다 %w", err)
		view.ShowErrorMessage(err)
	}
}
