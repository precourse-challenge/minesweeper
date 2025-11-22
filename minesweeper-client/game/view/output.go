package view

import (
	"fmt"
	"minesweeper-core/board"
	"minesweeper-core/cell"
	"minesweeper-core/position"
	"minesweeper-core/util"
	"minesweeper-infrastructure/dto"
	"time"
)

func ShowGameStartMessage() {
	fmt.Println("🎮지뢰찾기 게임을 시작합니다!")
}

func ShowGameModeSelection() {
	fmt.Println("게임 모드를 선택하세요 (single / multi)")
}

func ShowPlayerJoined(playerId int) {
	fmt.Printf("\nPlayer%d (으)로 참가했습니다.\n", playerId)
}

func ShowOpponentWaitMessage() {
	fmt.Println("\n다른 플레이어를 기다리는 중...")
}

func AskGameLevel() {
	fmt.Println("\n난이도를 선택하세요 (easy / normal / hard)")
}

func ShowSelectedGameLevel(level string) {
	fmt.Printf("\n선택된 난이도: %s\n\n", level)
}

func ShowBoard(board *board.Board) {
	showColNumbers(board)

	for row := 0; row < board.GetRowSize(); row++ {
		fmt.Printf("%2d  ", row+1)
		for col := 0; col < board.GetColSize(); col++ {
			cellPosition := util.FatalIfError(position.NewCellPosition(row, col))
			cellSnapshot := board.GetSnapshot(cellPosition)

			status := cellSnapshot.GetStatus()
			count := cellSnapshot.GetAdjacentLandMineCount()
			fmt.Printf("%2s ", signOf(status, count))
		}
		fmt.Println()
	}
	fmt.Println()
}

func ShowMultiBoards(board1Dto, board2Dto dto.BoardDto, playerId int) {
	fmt.Printf("\n       내 게임판 (Player%d)"+
		"               상대방 게임판 (Player%d)\n", playerId, 3-playerId)
	rows := len(board1Dto)
	cols := len(board1Dto[0])

	showMultiColumnNumbers(cols)

	var myBoard, enemyBoard dto.BoardDto
	if playerId == 1 {
		myBoard = board1Dto
		enemyBoard = board2Dto
	} else {
		myBoard = board2Dto
		enemyBoard = board1Dto
	}

	for i := 0; i < rows; i++ {
		fmt.Printf("%2d ", i+1)
		for j := 0; j < cols; j++ {
			fmt.Printf("%2s ", getCellSign(myBoard[i][j]))
		}

		fmt.Printf("  %2d ", i+1)
		for j := 0; j < cols; j++ {
			fmt.Printf("%2s ", getCellSign(enemyBoard[i][j]))
		}
		fmt.Println()
	}
}

func ShowRemainingFlagCount(board *board.Board) {
	remainingFlagCount := board.GetRemainingFlags()
	fmt.Printf("남은 깃발 개수: %d\n", remainingFlagCount)
}

func ShowTotalElapsedTime(elapsedTime time.Duration) {
	minutes := int(elapsedTime.Minutes())
	seconds := int(elapsedTime.Seconds()) % 60

	fmt.Printf("총 소요 시간: %02d:%02d\n\n", minutes, seconds)
}

func AskCommand() {
	fmt.Println("\n명령어를 입력해주세요 (open x y / flag x y / exit)")
}

func ShowCompletionMessage() {
	fmt.Println("모든 지뢰를 찾았습니다! 🎉🎉")
}

func ShowHitMineMessage() {
	fmt.Println("지뢰를 밟았습니다! 💣💣 게임 종료🥺")
}

func ShowWinMessage() {
	fmt.Println("축하합니다! 승리하셨습니다!🎉🎉")
}

func ShowLoseMessage() {
	fmt.Println("패배했습니다. 다음 기회에...")
}

func ShowRestartMessage() {
	fmt.Println("\n게임을 재시작하시겠습니까?")
}

func ShowQuitMessage() {
	fmt.Println("\n프로그램을 종료합니다.")
}

func ShowErrorMessage(err error) {
	fmt.Println("\n[ERROR] " + err.Error() + "\n")
}

func showColNumbers(board *board.Board) {
	colNumbers := generateColNumbers(board.GetColSize())

	fmt.Print("    ")
	for _, n := range colNumbers {
		fmt.Printf("%2d ", n)
	}
	fmt.Println()
}

func showMultiColumnNumbers(cols int) {
	fmt.Print("   ")
	for j := 1; j <= cols; j++ {
		fmt.Printf("%2d ", j)
	}
	fmt.Print("     ")
	for j := 1; j <= cols; j++ {
		fmt.Printf("%2d ", j)
	}
	fmt.Println()
}

func generateColNumbers(colSize int) []int {
	numbers := make([]int, 0, colSize)
	for i := 1; i <= colSize; i++ {
		numbers = append(numbers, i)
	}
	return numbers
}

func signOf(status cell.SnapshotStatus, adjacentLandMineCount int) string {
	switch status {
	case cell.Empty:
		return "■"
	case cell.Flag:
		return "⚑"
	case cell.LandMine:
		return "☼"
	case cell.Number:
		return fmt.Sprintf("%d", adjacentLandMineCount)
	default:
		return "□"
	}
}

func getCellSign(snapshotDto dto.CellSnapshotDto) string {
	return signOf(snapshotDto.Status, snapshotDto.Number)
}
