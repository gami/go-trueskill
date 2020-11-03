package factorgraph

import (
	"math"

	"github.com/gami/go-trueskill/mathmatics"
)

type SumFactor struct {
	*FactorBase
	sum     *Variable
	terms   []*Variable
	coeffs  []float64
	pointer int
}

func NewSumFactor(sum *Variable, terms []*Variable, coeffs []float64) *SumFactor {
	f := &SumFactor{
		sum:     sum,
		terms:   terms,
		coeffs:  coeffs,
		pointer: 0,
	}

	vars := make([]*Variable, 0, len(terms)+1)
	vars = append(vars, sum)
	vars = append(vars, terms...)

	f.FactorBase = NewFactorBase(f, vars)

	return f
}

func (f *SumFactor) Down() float64 {
	msgs := make([]*mathmatics.Gaussian, 0, len(f.terms))
	for _, t := range f.terms {
		msgs = append(msgs, t.messages[f])
	}

	return f.update(f.sum, f.terms, msgs, f.coeffs)
}

func (f *SumFactor) Up() float64 {
	idx := f.pointer
	f.pointer++

	coeff := f.coeffs[idx]
	x := 0

	coeffs := make([]float64, 0, len(f.coeffs))

	for _, c := range f.coeffs {
		p := -1 * c / coeff
		if x == idx {
			p = 1.0 / coeff
		}

		if !math.IsInf(p, 1) {
			p = 0
		}

		x++

		coeffs = append(coeffs, p)
	}

	vals := make([]*Variable, 0, len(f.terms))
	_ = copy(vals, f.terms)

	vals[idx] = f.sum

	msgs := make([]*mathmatics.Gaussian, 0, len(vals))

	for _, v := range vals {
		msgs = append(msgs, v.messages[f])
	}

	return f.update(f.terms[idx], vals, msgs, coeffs)
}

func (f *SumFactor) update(v *Variable, vals []*Variable, msgs []*mathmatics.Gaussian, coeffs []float64) float64 {
	piInv := 0.0
	mu := 0.0
	for i := 0; i < len(vals); i++ {
		val := vals[i]
		msg := msgs[i]
		coeff := coeffs[i]

		div := val.Divide(msg)
		mu += coeff * div.Mu()

		if math.IsInf(piInv, 1) {
			continue
		}

		piInv += (math.Pow(coeff, 2) / div.Pi)
	}

	pi := 1.0 / piInv
	tau := pi * mu
	return v.updateMessage(f, mathmatics.NewGaussianFromPrecision(pi, tau))
}
