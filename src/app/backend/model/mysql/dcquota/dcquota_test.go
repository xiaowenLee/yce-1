package dcquota

import (
	"fmt"
	"testing"
)

func Test_NewDcQuota(*testing.T) {
	dcQuota := NewDcQuota(1, 1, 1, 1000, 10, 20, 1, 2, 100, 10, 0, 1, "1000", "add dcquota")
	fmt.Printf("%v\n", dcQuota)
}