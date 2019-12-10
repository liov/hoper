package mobile

import (
	"fmt"

	"golang.org/x/mobile/event/key"
)

func Mobile() {
	fmt.Println(key.Event{}.String())
}
