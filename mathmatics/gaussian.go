package mathmatics

import "math"

// Gaussian represents a model for the normal distribution.
type Gaussian struct {
	tau float64
	pi  float64
}

func NewGaussian(pi float64, tau float64) *Gaussian {
	return &Gaussian{
		pi:  pi,
		tau: tau,
	}
}

func NewGaussianFromDistribution(mu float64, sigma float64) *Gaussian {
	pi := math.Pow(sigma, -2)
	tau := pi * mu

	return &Gaussian{
		pi:  pi,
		tau: tau,
	}
}

func (g *Gaussian) Mu() float64 {
	if g.pi == 0 {
		return 0
	}
	return g.tau / g.pi
}

func (g *Gaussian) Sigma() float64 {
	if g.pi == 0 {
		return math.Inf(0)
	}
	return math.Sqrt(1 / g.pi)
}

func (g *Gaussian) Multiply(a *Gaussian) *Gaussian {
	return &Gaussian{
		pi:  g.pi + a.pi,
		tau: g.tau + a.tau,
	}
}

func (g *Gaussian) Divide(a *Gaussian) *Gaussian {
	return &Gaussian{
		pi:  g.pi - a.pi,
		tau: g.tau - a.tau,
	}
}

func (g *Gaussian) Equals(a *Gaussian) bool {
	return g.pi == a.pi && g.tau == a.tau
}

func (g *Gaussian) LessThan(a *Gaussian) bool {
	return g.Mu() < a.Mu()
}

func (g *Gaussian) LessThanEqual(a *Gaussian) bool {
	return g.Mu() <= a.Mu()
}

func (g *Gaussian) GreaterThan(a *Gaussian) bool {
	return g.Mu() < a.Mu()
}

func (g *Gaussian) GreaterEqual(a *Gaussian) bool {
	return g.Mu() <= a.Mu()
}
