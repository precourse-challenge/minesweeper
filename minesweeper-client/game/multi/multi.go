package multi

import (
	"fmt"
	"log"
	"minesweeper-client/game/user"
	"minesweeper-client/game/view"
	"minesweeper-core/position"
	"os"
	"strconv"
	"strings"
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

	<-m.session.joinDone

	for {
		action, cellPosition, err := m.readCommand()
		if err != nil {
			view.ShowErrorMessage(err)
			continue
		}

		err = m.handleActionOnCell(action, cellPosition)
		if err != nil {
			view.ShowErrorMessage(err)
		}
	}
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

	inputAction := commands[0]
	action := user.ActionFrom(inputAction)

	if action == user.Exit {
		view.ShowQuitMessage()
		os.Exit(0)
	}

	if len(commands) != 3 {
		return user.UnknownAction, nil, fmt.Errorf("명령어 형식이 올바르지 않습니다")
	}

	inputRow := commands[1]
	inputCol := commands[2]

	row, err := m.getSelectedRowIndex(inputRow)
	if err != nil {
		return user.UnknownAction, nil, err
	}

	col, err := m.getSelectedColIndex(inputCol)
	if err != nil {
		return user.UnknownAction, nil, err
	}

	cellPosition, err := position.NewCellPosition(row, col)
	if err != nil {
		return user.UnknownAction, nil, err
	}

	return action, cellPosition, nil
}

func (m *MultiMode) getSelectedColIndex(inputCol string) (int, error) {
	col, err := strconv.Atoi(inputCol)
	if err != nil {
		return 0, fmt.Errorf("좌표는 숫자여야 합니다")
	}
	return col - 1, nil
}

func (m *MultiMode) getSelectedRowIndex(inputRow string) (int, error) {
	row, err := strconv.Atoi(inputRow)
	if err != nil {
		return 0, fmt.Errorf("좌표는 숫자여야 합니다")
	}
	return row - 1, nil
}

func (m *MultiMode) handleActionOnCell(action user.Action, cellPosition *position.CellPosition) error {
	if action == user.Open {
		err := m.session.Open(cellPosition.RowIndex(), cellPosition.ColIndex())
		if err != nil {
			return err
		}
		return nil
	}
	if action == user.Flag {
		err := m.session.Flag(cellPosition.RowIndex(), cellPosition.ColIndex())
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("잘못된 명령어를 입력했습니다")
}

func (m *MultiMode) closeConnection(gameClient *Session) {
	err := gameClient.Close()
	if err != nil {
		err := fmt.Errorf("클라이언트 연결 종료를 실패했습니다 %w", err)
		view.ShowErrorMessage(err)
	}
}
