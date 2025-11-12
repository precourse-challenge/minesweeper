package view

import "fmt"

func ShowGameStartMessage() {
	fmt.Println("🎮지뢰찾기 게임을 시작합니다!")
}
func AskGameLevel() {
	fmt.Println("난이도를 선택하세요 (easy / normal / hard)")
}
func ShowSelectedGameLevel(level string) {
	fmt.Printf("선택된 난이도: %s\n", level)
}
