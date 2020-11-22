package trueskill

import "github.com/gami/go-trueskill/mathmatics"

// Rating represents a player’s skill as Gaussian distrubution.
type Rating struct {
	mu     float64 // the mean.
	sigma  float64 // the square root of the variance.
	weight float64 // default 1
}

func (r *Rating) gaussian() *mathmatics.Gaussian {
	return mathmatics.NewGaussianFromDistribution(r.mu, r.sigma)
}
