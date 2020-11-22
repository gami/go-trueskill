package factorgraph

import "github.com/gami/go-trueskill/mathmatics"

type LikelihoodFactor struct {
	*FactorBase
	mean     *Variable
	value    *Variable
	variance float64
}

func NewLikelihoodFactor(mean *Variable, value *Variable, variance float64) *LikelihoodFactor {
	f := &LikelihoodFactor{
		mean:     mean,
		value:    value,
		variance: variance,
	}

	f.FactorBase = NewFactorBase(f, []*Variable{mean, value})
	return f
}

func (f *LikelihoodFactor) calcA(v *mathmatics.Gaussian) float64 {
	return 1.0 / ((f.variance * v.Pi) + 1.0)
}

func (f *LikelihoodFactor) Down() float64 {
	msg := f.mean.Divide(f.mean.messages[f])
	a := f.calcA(msg)
	return f.value.updateMessage(f, mathmatics.NewGaussian(a*msg.Pi, a*msg.Tau))
}

func (f *LikelihoodFactor) Up() (float64, error) {
	msg := f.value.Divide(f.value.messages[f])
	a := f.calcA(msg)
	return f.mean.updateMessage(f, mathmatics.NewGaussian(a*msg.Pi, a*msg.Tau)), nil
}
