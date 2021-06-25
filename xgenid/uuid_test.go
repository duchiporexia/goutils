package xgenid

import (
	"fmt"
	"testing"
)

func TestUuid(t *testing.T) {
	fmt.Printf("uuid:%s\n", Uuid())
}
