package mode

import (
	"fmt"
	"minesweeper-client/game/view"
	"minesweeper-core/board"
	"minesweeper-core/level"
	"minesweeper-core/position"
	"minesweeper-core/user"
	"os"
	"strconv"
	"strings"
	"time"
)

type SingleMode struct {
	board *board.Board
}

func NewSingleMode(level level.GameLevel) *SingleMode {
	return &SingleMode{board: board.NewBoard(level)}
}

func (m *SingleMode) Start() {
	m.board.InitializeGame()

	startTime := time.Now()

	for m.board.IsInProgress() {
		view.ShowBoard(m.board)
		view.ShowRemainingFlagCount(m.board)

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

	view.ShowBoard(m.board)

	view.ShowTotalElapsedTime(time.Since(startTime))

	if m.board.IsWinStatus() {
		view.ShowWinMessage()
	}
	if m.board.IsLoseStatus() {
		view.ShowLoseMessage()
	}
}

func (m *SingleMode) readCommand() (user.Action, *position.CellPosition, error) {
	view.AskCommand()
	inputCommand := view.Read()

	action, cellPosition, err := m.parseCommand(inputCommand)
	if err != nil {
		return user.Unknown, nil, err
	}
	return action, cellPosition, nil
}

func (m *SingleMode) handleActionOnCell(action user.Action, cellPosition *position.CellPosition) error {
	if action == user.Open {
		err := m.board.Open(cellPosition)
		if err != nil {
			return err
		}
		return nil
	}
	if action == user.Flag {
		err := m.board.Flag(cellPosition)
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("잘못된 명령어를 입력했습니다")
}

func (m *SingleMode) parseCommand(inputCommand string) (user.Action, *position.CellPosition, error) {
	commands := strings.Fields(inputCommand)

	inputAction := commands[0]
	action := user.From(inputAction)

	if action == user.Exit {
		view.ShowQuitMessage()
		os.Exit(0)
	}

	if len(commands) != 3 {
		return user.Unknown, nil, fmt.Errorf("명령어 형식이 올바르지 않습니다")
	}

	inputRow := commands[1]
	inputCol := commands[2]

	row, err := m.getSelectedRowIndex(inputRow)
	if err != nil {
		return user.Unknown, nil, err
	}

	col, err := m.getSelectedColIndex(inputCol)
	if err != nil {
		return user.Unknown, nil, err
	}

	cellPosition, err := position.NewCellPosition(row, col)
	if err != nil {
		return user.Unknown, nil, err
	}

	if m.board.IsOutOfBounds(cellPosition) {
		return user.Unknown, nil, fmt.Errorf("올바르지 않은 좌표값입니다")
	}

	return action, cellPosition, nil
}

func (m *SingleMode) getSelectedColIndex(inputCol string) (int, error) {
	col, err := strconv.Atoi(inputCol)
	if err != nil {
		return 0, fmt.Errorf("좌표는 숫자여야 합니다")
	}
	return col - 1, nil
}

func (m *SingleMode) getSelectedRowIndex(inputRow string) (int, error) {
	row, err := strconv.Atoi(inputRow)
	if err != nil {
		return 0, fmt.Errorf("좌표는 숫자여야 합니다")
	}
	return row - 1, nil
}
