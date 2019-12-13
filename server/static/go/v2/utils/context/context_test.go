package context2

import (
	"fmt"
	"testing"
)

func TestContext(t *testing.T) {
	ctx := valueContext{}
	fmt.Println(ctx)
}
