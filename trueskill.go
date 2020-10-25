package trueskill

import (
	"errors"
)

// Trueskill represents envirionment of rating
type TrueSkill struct {
	mu              float64 // the initial mean of ratings.
	sigma           float64 // the initial standard deviation of ratings. The recommended value is a third of mu.
	beta            float64 // the distance which guarantees about 76% chance of winning. The recommended value is a half of sigma.
	tau             float64 // the dynamic factor which restrains a fixation of rating. The recommended value is sigma per cent.
	drawProbability float64 // the draw probability between two teams. It can be a float or function which returns a float by the given two rating (team performance) arguments and the beta value. If it is a float, the game has fixed draw probability. Otherwise, the draw probability will be decided dynamically per each match.
}

type option func(*TrueSkill)

func NewTrueSkill(options ...option) *TrueSkill {
	s := &TrueSkill{
		mu:              25.0,
		sigma:           8.333333333333334,
		beta:            4.166666666666667,
		tau:             0.08333333333333334,
		drawProbability: 0.1,
	}

	for _, opt := range options {
		opt(s)
	}

	return s
}

func MU(v float64) option {
	return func(s *TrueSkill) {
		s.mu = v
	}
}

func Sigma(v float64) option {
	return func(s *TrueSkill) {
		s.sigma = v
	}
}

func Beta(v float64) option {
	return func(s *TrueSkill) {
		s.beta = v
	}
}

func Tau(v float64) option {
	return func(s *TrueSkill) {
		s.tau = v
	}
}

func DrawProbability(v float64) option {
	return func(s *TrueSkill) {
		s.drawProbability = v
	}
}

func (s *TrueSkill) CreateRating() *Rating {
	return &Rating{
		mu:    s.mu,
		sigma: s.sigma,
	}
}

// Rate recalculates ratings by the ranking table:
func (s *TrueSkill) Rate(ratingGroups []RatingGroup) ([]RatingGroup, error) {
	if err := s.validateRatingGroup(ratingGroups); err != nil {
		return nil, err
	}

	// # build factor graph
	// args = (sorted_rating_groups, sorted_ranks, sorted_weights)
	// builders = self.factor_graph_builders(*args)
	// args = builders + (min_delta,)
	// layers = self.run_schedule(*args)

	// # make result
	// rating_layer, team_sizes = layers[0], _team_sizes(sorted_rating_groups)
	// transformed_groups = []
	// for start, end in zip([0] + team_sizes[:-1], team_sizes):
	//     group = []
	//     for f in rating_layer[start:end]:
	//         group.append(Rating(float(f.var.mu), float(f.var.sigma)))
	//     transformed_groups.append(tuple(group))
	// by_hint = lambda x: x[0]
	// unsorting = sorted(zip((x for x, __ in sorting), transformed_groups),
	//                    key=by_hint)
	// if keys is None:
	//     return [g for x, g in unsorting]
	// # restore the structure with input dictionary keys
	// return [dict(zip(keys[x], g)) for x, g in unsorting]

	return ratingGroups, nil
}

func (s *TrueSkill) validateRatingGroup(ratingGroups []RatingGroup) error {
	if len(ratingGroups) < 2 {
		return errors.New("need multiple rating groups")
	}

	for _, rs := range ratingGroups {
		if len(rs) < 2 {
			return errors.New("each group must contain multiple ratings")
		}
	}

	return nil
}
