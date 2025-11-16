package level

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_입력받은_난이도를_level로_파싱한다(t *testing.T) {
	// given
	tests := []string{"easy", "normal", "hard"}

	// when & then
	for _, input := range tests {
		gameLevel, _ := From(input)
		assert.NotNil(t, gameLevel)
	}
}

func Test_잘못된_난이도를_입력하면_에러가_발생한다(t *testing.T) {
	// given
	tests := []string{"", "wrong", "medium"}

	// when & then
	for _, input := range tests {
		_, err := From(input)
		assert.EqualError(t, err, "난이도는 (easy / normal / hard) 중 하나를 입력해야 합니다")
	}
}
