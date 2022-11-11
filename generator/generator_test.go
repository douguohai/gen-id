package generator

import (
	"fmt"
	"github.com/douguohai/gen-id/utils"
	"testing"
)

func TestXPaddingZeroForNumberStart(t *testing.T) {
	fmt.Println(utils.PaddingZeroForNumberStart(6, "010000000"))
}
