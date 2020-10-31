package factorgraph

import (
	"math"

	"github.com/gami/go-trueskill/mathmatics"
)

type PriorFactor struct {
	variable *Variable
	value    *mathmatics.Gaussian
	dynamic  float64
}

func (f *PriorFactor) down() float64 {
	sigma := math.Sqrt(math.Pow(f.value.Sigma(), 2) + math.Pow(f.dynamic, 2))
	val := mathmatics.NewGaussian(f.value.Mu(), sigma)

	return f.variable.updateMessage(val)
}
