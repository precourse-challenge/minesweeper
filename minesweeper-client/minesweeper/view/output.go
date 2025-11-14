package view

import (
	"fmt"
	"minesweeper-client/minesweeper/board"
	"minesweeper-client/minesweeper/cell"
	"minesweeper-client/minesweeper/position"
	"minesweeper-client/minesweeper/util"
)

func ShowGameStartMessage() {
	fmt.Println("🎮지뢰찾기 게임을 시작합니다!")
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

			fmt.Printf("%2s ", signOf(cellSnapshot))
		}
		fmt.Println()
	}
	fmt.Println()
}

func ShowRemainingFlagCount(board *board.Board) {
	remainingFlagCount := board.GetRemainingFlags()
	fmt.Printf("남은 깃발 개수: %d\n\n", remainingFlagCount)
}

func AskCommand() {
	fmt.Println("명령어를 입력해주세요 (open x y / flag x y / quit)")
}

func ShowWinMessage() {
	fmt.Println("모든 지뢰를 찾았습니다! 🎉🎉")
}

func ShowLoseMessage() {
	fmt.Println("지뢰를 밟았습니다! 💣💣 게임 종료🥺")
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

func generateColNumbers(colSize int) []int {
	numbers := make([]int, 0, colSize)
	for i := 1; i <= colSize; i++ {
		numbers = append(numbers, i)
	}
	return numbers
}

func signOf(snapshot cell.Snapshot) string {
	switch snapshot.GetStatus() {
	case cell.Empty:
		return "■"
	case cell.Flag:
		return "⚑"
	case cell.LandMine:
		return "☼"
	case cell.Number:
		return fmt.Sprintf("%d", snapshot.GetAdjacentLandMineCount())
	default:
		return "□"
	}
}
