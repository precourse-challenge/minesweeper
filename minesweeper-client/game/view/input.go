package view

import (
	"bufio"
	"os"
	"strings"
)

func Read() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
