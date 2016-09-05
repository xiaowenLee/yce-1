package time

import (
	"fmt"
	"testing"
)

func Test_LocalTime(*testing.T) {
	l := NewLocalTime()
	fmt.Printf("%s\n", l.String())
}
