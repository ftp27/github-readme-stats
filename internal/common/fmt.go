package common

import (
	"fmt"
	"math"
)

func KFormatter(num float64, precision *int) string {
	abs := math.Abs(num)
	sign := math.Copysign(1, num)

	if precision != nil {
		return fmt.Sprintf("%.*fk", *precision, sign*(abs/1000))
	}

	if abs < 1000 {
		if math.Mod(abs, 1) == 0 {
			return fmt.Sprintf("%.0f", sign*abs)
		}
		return fmt.Sprintf("%g", sign*abs)
	}

	return fmt.Sprintf("%.1fk", sign*(abs/1000))
}
