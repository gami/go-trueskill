package factorgraph

type Factor interface {
	Up() float64
	Down() float64
	Var() *Variable
	ToString() string
}

type FactorBase struct {
}

func (f *FactorBase) Up(*FactorBase) float64 {
	return 0
}

func (f *FactorBase) Down(*FactorBase) float64 {
	return 0
}

// def __init__(self, variables):
// self.vars = variables
// for var in variables:
// 	var[self] = Gaussian()

// @property
// def var(self):
// assert len(self.vars) == 1
// return self.vars[0]
