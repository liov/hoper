package constraints

import (
	"golang.org/x/exp/constraints"
	"time"
)

type Integer interface {
}

type Number interface {
}

type Ordered interface {
	constraints.Ordered | time.Time
}
