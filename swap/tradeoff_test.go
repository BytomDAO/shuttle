package swap

import (
	"fmt"
	"testing"
)

func TestParseUint64(t *testing.T) {
	s := "12"
	num, err := ParseUint64(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(num)
}
