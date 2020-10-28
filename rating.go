package trueskill

// Rating represents a player’s skill as Gaussian distrubution.
type Rating struct {
	mu    float64 // the mean.
	sigma float64 // the square root of the variance.
}

// Expose returns the value of the rating exposure.  It starts from 0 and
// converges to the mean.
func (r *Rating) Expose() float64 {
	k := r.mu / r.sigma
	return r.mu - k*r.sigma
}

type RatingGroup map[string]*Rating
