package factorgraph

import "github.com/gami/go-trueskill/mathmatics"

type Factor interface {
	Up() float64
	Down() float64
	Var() *Variable
}

type FactorBase struct {
	Vars []*Variable
}

func NewFactorBase(factor Factor, vars []*Variable) *FactorBase {
	f := &FactorBase{
		Vars: vars,
	}

	for _, v := range vars {
		v.messages[factor] = &mathmatics.Gaussian{}
	}

	return f
}

func (f *FactorBase) Up() float64 {
	return 0
}

func (f *FactorBase) Down() float64 {
	return 0
}

func (f *FactorBase) Var() *Variable {
	if len(f.Vars) == 0 {
		return nil
	}

	return f.Vars[0]
}
