package factorgraph

import (
	"math"

	"github.com/gami/go-trueskill/mathmatics"
)

type PriorFactor struct {
	*FactorBase
	variable *Variable
	value    *mathmatics.Gaussian
	dynamic  float64
}

func NewPriorFactor(v *Variable, val *mathmatics.Gaussian, dynamic float64) *PriorFactor {
	f := &PriorFactor{
		variable: v,
		value:    val,
		dynamic:  dynamic,
	}

	f.FactorBase = NewFactorBase(f, []*Variable{v})

	return f
}

func (f *PriorFactor) Down() float64 {
	sigma := math.Sqrt(math.Pow(f.value.Sigma(), 2) + math.Pow(f.dynamic, 2))
	val := mathmatics.NewGaussianFromDistribution(f.value.Mu(), sigma)

	return f.variable.updateValue(f, NewVariable(val))
}
