package factorgraph

type PriorFactor struct {
	FactorBase
}

func (f *PriorFactor) down() float64 {

}

class PriorFactor(Factor):

    def __init__(self, var, val, dynamic=0):
        super(PriorFactor, self).__init__([var])
        self.val = val
        self.dynamic = dynamic

    def down(self):
        sigma = math.sqrt(self.val.sigma ** 2 + self.dynamic ** 2)
        value = Gaussian(self.val.mu, sigma)
        return self.var.update_value(self, value=value)