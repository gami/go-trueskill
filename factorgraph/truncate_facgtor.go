package factorgraph

import (
	"math"

	"github.com/gami/go-trueskill/mathmatics"
)

type TruncateFactor struct {
	*FactorBase
	v          *Variable
	vFunc      func(a float64, b float64) float64
	wFunc      func(a float64, b float64) (float64, error)
	drawMargin float64
}

func NewTruncateFactor(
	v *Variable,
	vFunc func(a float64, b float64) float64,
	wFunc func(a float64, b float64) (float64, error),
	drawMargin float64) *TruncateFactor {

	f := &TruncateFactor{
		v:          v,
		vFunc:      vFunc,
		wFunc:      wFunc,
		drawMargin: drawMargin,
	}

	f.FactorBase = NewFactorBase(f, []*Variable{v})

	return f
}

func (f *TruncateFactor) Up() (float64, error) {
	val := f.v
	msg := f.v.messages[f]
	div := val.Divide(msg)
	sqrtPi := math.Sqrt(div.Pi)
	v := f.vFunc(div.Tau/sqrtPi, f.drawMargin*sqrtPi)
	w, err := f.wFunc(div.Tau/sqrtPi, f.drawMargin*sqrtPi)
	if err != nil {
		return 0, err
	}
	denom := 1.0 - w
	pi := div.Pi / denom
	tau := (div.Tau + (sqrtPi * v)) / denom
	return val.updateValue(f, NewVariable(mathmatics.NewGaussian(pi, tau))), nil
}
