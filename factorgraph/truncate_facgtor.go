package factorgraph

import (
	"math"

	"github.com/gami/go-trueskill/mathmatics"
)

type TruncateFactor struct {
	*FactorBase
	v          *Variable
	vFunc      func(a float64, b float64) float64
	wFunc      func(a float64, b float64) float64
	drawMargin float64
}

func NewTruncateFactor(
	v *Variable,
	vFunc func(a float64, b float64) float64,
	wFunc func(a float64, b float64) float64,
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

func (f *TruncateFactor) Up() float64 {
	val := f.v
	msg := f.v.messages[f]
	div := val.Divide(msg)
	sqrtPi := math.Sqrt(div.Pi)
	v := f.vFunc(div.Tau/sqrtPi, f.drawMargin*sqrtPi)
	w := f.wFunc(div.Tau/sqrtPi, f.drawMargin*sqrtPi)
	denom := 1.0 - w
	pi := div.Pi / denom
	tau := (div.Tau + (sqrtPi * v)) / denom
	return val.updateValue(f, NewVariable(mathmatics.NewGaussianFromPrecision(pi, tau)))
}
