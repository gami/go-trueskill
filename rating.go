package trueskill

import "github.com/gami/go-trueskill/mathmatics"

// Rating represents a player’s skill as Gaussian distrubution.
type Rating struct {
	Mu     float64 // the mean.
	Sigma  float64 // the square root of the variance.
	Weight float64 // default 1
}

type ratingOpt func(*Rating)

func NewRating(mu float64, sigma float64, weight float64) *Rating {
	return &Rating{
		Mu:     mu,
		Sigma:  sigma,
		Weight: weight,
	}
}

func (r *Rating) gaussian() *mathmatics.Gaussian {
	return mathmatics.NewGaussianFromDistribution(r.Mu, r.Sigma)
}
