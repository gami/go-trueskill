package trueskill

// Rating represents a player’s skill as Gaussian distrubution.
type Rating struct {
	mu    float64 // the mean.
	sigma float64 // the square root of the variance.
}

type RatingGroup map[string]*Rating
