package generator

import (
	"fmt"
	"github.com/douguohai/gen-id/utils"
	"testing"
)

func TestXPaddingZeroForNumberStart(t *testing.T) {
	fmt.Println(utils.PaddingZeroForNumberStart(6, "010000000"))
}

func TestXGetBirthDay(t *testing.T) {
	idcard := utils.NewIDCard("342423199604212299")
	fmt.Println(idcard.GetBirthDay())
}
