package mathmatics

import "math"

// Gaussian represents a model for the normal distribution.
type Gaussian struct {
	Pi  float64 // Precision, the inverse of the variance.
	Tau float64 // Precision adjusted mean, the precision multiplied by the mean.
}

func NewGaussianFromPrecision(pi float64, tau float64) *Gaussian {
	return &Gaussian{
		Pi:  pi,
		Tau: tau,
	}
}

func NewGaussianFromDistribution(mu float64, sigma float64) *Gaussian {
	pi := math.Pow(sigma, -2)
	tau := pi * mu

	return &Gaussian{
		Pi:  pi,
		Tau: tau,
	}
}

func (g *Gaussian) Mu() float64 {
	if g.Pi == 0 {
		return 0
	}
	return g.Tau / g.Pi
}

func (g *Gaussian) Sigma() float64 {
	if g.Pi == 0 {
		return math.Inf(0)
	}
	return math.Sqrt(1.0 / g.Pi)
}

func (g *Gaussian) Multiply(a *Gaussian) *Gaussian {
	return &Gaussian{
		Pi:  g.Pi + a.Pi,
		Tau: g.Tau + a.Tau,
	}
}

func (g *Gaussian) Divide(a *Gaussian) *Gaussian {
	return &Gaussian{
		Pi:  g.Pi - a.Pi,
		Tau: g.Tau - a.Tau,
	}
}

func (g *Gaussian) Equals(a *Gaussian) bool {
	return g.Pi == a.Pi && g.Tau == a.Tau
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
