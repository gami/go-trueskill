package factorgraph

import (
	"math"

	"github.com/gami/go-trueskill/mathmatics"
)

type Variable struct {
	*mathmatics.Gaussian
	messages map[Factor]*mathmatics.Gaussian
}

func NewVariable(g *mathmatics.Gaussian) *Variable {
	return &Variable{
		Gaussian: g,
		messages: make(map[Factor]*mathmatics.Gaussian),
	}
}

func (v *Variable) set(other *Variable) float64 {
	delta := v.delta(other)
	v.Pi = other.Pi
	v.Tau = other.Tau
	return delta
}

func (v *Variable) delta(other *Variable) float64 {
	piDelta := math.Abs(v.Pi - other.Pi)

	if piDelta == math.Inf(1) {
		return 0
	}

	return math.Max(math.Abs(v.Tau-other.Tau), math.Sqrt(piDelta))
}

func (v *Variable) updateMessage(factor Factor, msg *mathmatics.Gaussian) float64 {
	oldMessage := v.messages[factor]
	v.messages[factor] = msg
	return v.set(NewVariable(v.Divide(oldMessage).Multiply(msg)))
}

func (v *Variable) updateValue(factor Factor, val *Variable) float64 {
	oldMessage := v.messages[factor]
	v.messages[factor] = val.Multiply(oldMessage).Divide(v.Gaussian)

	return v.set(val)
}
