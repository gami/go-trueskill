package trueskill

import (
	"errors"
	"fmt"
	"math"

	"github.com/chobie/go-gaussian"
	"github.com/gami/go-trueskill/factorgraph"
	"github.com/gami/go-trueskill/mathmatics"
)

const (
	// MinDelta is a basis to check reliability of the result.
	MinDelta = 0.001
)

// TrueSkill represents envirionment of rating
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
		mu:     s.mu,
		sigma:  s.sigma,
		weight: 1,
	}
}

// Rate recalculates ratings by the ranking table:
func (s *TrueSkill) Rate(ratingGroups [][]*Rating) ([][]*Rating, error) {
	if err := s.validateRatingGroup(ratingGroups); err != nil {
		return nil, err
	}

	flattenRatings := make([]*Rating, 0)
	sortedRanks := make([]int, 0)
	rank := 0
	for _, rg := range ratingGroups {
		for _, r := range rg {
			flattenRatings = append(flattenRatings, r)
			sortedRanks = append(sortedRanks, rank)
		}
		rank++
	}

	ratingVars := make([]*factorgraph.Variable, 0, len(flattenRatings))
	perfVars := make([]*factorgraph.Variable, 0, len(flattenRatings))
	flattenWeights := make([]float64, 0, len(flattenRatings))
	for _, r := range flattenRatings {
		ratingVars = append(ratingVars, factorgraph.NewVariable(mathmatics.NewGaussianFromDistribution(0, 0)))
		perfVars = append(perfVars, factorgraph.NewVariable(mathmatics.NewGaussianFromDistribution(0, 0)))
		flattenWeights = append(flattenWeights, r.weight)
	}

	teamPerfVars := make([]*factorgraph.Variable, 0, len(ratingGroups))
	for i := 0; i < len(teamPerfVars); i++ {
		teamPerfVars = append(teamPerfVars, factorgraph.NewVariable(mathmatics.NewGaussianFromDistribution(0, 0)))
	}

	teamDiffVars := make([]*factorgraph.Variable, 0, len(ratingGroups)-1)
	for i := 0; i < len(ratingGroups)-1; i++ {
		teamDiffVars = append(teamDiffVars, factorgraph.NewVariable(mathmatics.NewGaussianFromDistribution(0, 0)))
	}

	teamSizes := teamSizes(ratingGroups)

	layers, err := s.runSchedule(
		ratingVars,
		flattenRatings,
		perfVars,
		teamPerfVars,
		teamSizes,
		flattenWeights,
		teamDiffVars,
		sortedRanks,
		ratingGroups,
	)
	if err != nil {
		return nil, err
	}

	transformedGroups := make([][]*Rating, 0, len(teamSizes))

	trimmed := []int{0}
	trimmed = append(trimmed, teamSizes[0:len(teamSizes)-1]...)

	for i := 0; i < len(teamSizes); i++ {
		group := make([]*Rating, 0)
		glayers := layers[trimmed[i]:teamSizes[i]]
		for _, layer := range glayers {
			r := &Rating{
				mu:    layer.Var().Mu(),
				sigma: layer.Var().Sigma(),
			}
			group = append(group, r)
		}
		transformedGroups = append(transformedGroups, group)
	}

	return transformedGroups, nil
}

func (s *TrueSkill) validateRatingGroup(ratingGroups [][]*Rating) error {
	if len(ratingGroups) < 2 {
		return errors.New("need multiple rating groups")
	}

	for _, rs := range ratingGroups {
		if len(rs) < 1 {
			return errors.New("each group must contain multiple ratings")
		}
	}

	return nil
}

// runSchedule sends messages within every nodes of the factor graph until the result is reliable.
func (s *TrueSkill) runSchedule(
	ratingVars []*factorgraph.Variable,
	flattenRatings []*Rating,
	perfVars []*factorgraph.Variable,
	teamPerfVars []*factorgraph.Variable,
	teamSizes []int,
	flattenWeights []float64,
	teamDiffVars []*factorgraph.Variable,
	sortedRanks []int,
	sortedRatingGroups [][]*Rating,
) ([]*factorgraph.PriorFactor, error) {
	ratingLayer := s.buildRatingLayer(ratingVars, flattenRatings)
	perfLayer := s.buildPerfLayer(ratingVars, perfVars)
	teamPerfLayer := s.buildTeamPerfLayer(
		teamPerfVars,
		perfVars,
		teamSizes,
		flattenWeights,
	)

	for _, f := range ratingLayer {
		f.Down()
	}
	for _, f := range perfLayer {
		f.Down()
	}
	for _, f := range teamPerfLayer {
		f.Down()
	}

	// Arrow #1, #2, #3
	teamDiffLayer := s.buildTeamDiffLayer(teamPerfVars, teamDiffVars)
	truncLayer := s.buildTruncLayer(teamDiffVars, sortedRanks, sortedRatingGroups)
	teamDiffLen := len(teamDiffLayer)

	for index := 0; index <= 10; index++ {
		delta := 0.0
		if teamDiffLen == 1 {
			// Only two teams
			teamDiffLayer[0].Down()
			delta = truncLayer[0].Up()
		} else {
			// Multiple teams
			for z := 0; z < teamDiffLen-1; z++ {
				teamDiffLayer[z].Down()
				delta = math.Max(delta, truncLayer[z].Up())
				teamDiffLayer[z].SetPointer(1)
				teamDiffLayer[z].Up()
			}

			for z := teamDiffLen - 1; z > 0; z-- {
				teamDiffLayer[z].Down()
				delta = math.Max(delta, truncLayer[z].Up())
				teamDiffLayer[z].SetPointer(0)
				teamDiffLayer[z].Up()
			}
		}

		// Repeat until too small update
		if delta <= MinDelta {
			break
		}
	}

	// Up both ends
	teamDiffLayer[0].SetPointer(0)
	teamDiffLayer[0].Up()
	teamDiffLayer[teamDiffLen-1].SetPointer(1)
	teamDiffLayer[teamDiffLen-1].Up()

	// Up the remainder of the black arrows
	for _, f := range teamPerfLayer {
		f.SetPointer(0)
		for x := 0; x < len(f.Vars)-1; x++ {
			f.Up()
		}
	}

	for _, f := range perfLayer {
		f.Up()
	}

	return ratingLayer, nil
}

func (s *TrueSkill) buildRatingLayer(ratingVars []*factorgraph.Variable, flattenRatings []*Rating) []*factorgraph.PriorFactor {
	layers := make([]*factorgraph.PriorFactor, 0, len(ratingVars))

	for i, v := range ratingVars {
		f := factorgraph.NewPriorFactor(v, flattenRatings[i].gaussian(), s.tau)
		layers = append(layers, f)
	}

	return layers
}

func (s *TrueSkill) buildPerfLayer(ratingVars []*factorgraph.Variable, perfVars []*factorgraph.Variable) []factorgraph.Factor {
	layer := make([]factorgraph.Factor, 0, len(ratingVars))

	b := math.Pow(s.beta, 2)

	for i, v := range ratingVars {
		f := factorgraph.NewLikelihoodFactor(v, perfVars[i], b)
		layer = append(layer, f)
	}

	return layer
}

func (s *TrueSkill) buildTeamPerfLayer(
	teamPerfVars []*factorgraph.Variable,
	perfVars []*factorgraph.Variable,
	teamSizes []int,
	flattenWeights []float64,
) []*factorgraph.SumFactor {

	team := 0

	layer := make([]*factorgraph.SumFactor, 0, len(teamPerfVars))

	for _, v := range teamPerfVars {
		s := 0
		if team > 0 {
			s = teamSizes[team-1]
		}

		e := teamSizes[team]

		team++

		f := factorgraph.NewSumFactor(v, perfVars[s:e], flattenWeights[s:e])
		layer = append(layer, f)
	}

	return layer
}

func (s *TrueSkill) buildTeamDiffLayer(teamPerfVars []*factorgraph.Variable, teamDiffVars []*factorgraph.Variable) []*factorgraph.SumFactor {
	layer := make([]*factorgraph.SumFactor, 0, len(teamDiffVars))

	team := 0

	for _, v := range teamDiffVars {

		sl := teamPerfVars[team : team+2]
		team++

		f := factorgraph.NewSumFactor(v, sl, []float64{1, -1})
		layer = append(layer, f)
	}

	return layer
}

func (s *TrueSkill) buildTruncLayer(
	teamDiffVars []*factorgraph.Variable,
	sortedRanks []int,
	sortedRatingGroups [][]*Rating,
) []factorgraph.Factor {

	x := 0

	layer := make([]factorgraph.Factor, 0, len(teamDiffVars))

	for _, v := range teamDiffVars {
		size := 0

		for _, r := range sortedRatingGroups[x : x+2] {
			size += len(r)
		}

		drawMargin := s.calcDrawMargin(size)

		vFunc := func(a float64, b float64) float64 { return s.vWin(a, b) }
		wFunc := func(a float64, b float64) float64 { return s.wWin(a, b) }
		if sortedRanks[x] == sortedRanks[x+1] {
			vFunc = func(a float64, b float64) float64 { return s.vDraw(a, b) }
			wFunc = func(a float64, b float64) float64 { return s.wDraw(a, b) }
		}

		x++
		f := factorgraph.NewTruncateFactor(v, vFunc, wFunc, drawMargin)
		layer = append(layer, f)
	}

	return layer
}

func (s *TrueSkill) calcDrawMargin(
	size int,
) float64 {

	g := gaussian.NewGaussian(0.0, 1.0)
	return g.Ppf((s.drawProbability+1.0)/2.0) * math.Sqrt(float64(size)) * s.beta
}

// The non-draw version of "V" function.
// "V" calculates a variation of a mean.
func (s *TrueSkill) vWin(diff float64, drawMargin float64) float64 {
	x := diff - drawMargin
	g := gaussian.NewGaussian(0.0, 1.0)
	denom := g.Cdf(x)
	if denom != 0 && !math.IsNaN(denom) {
		return g.Pdf(x) / denom
	}

	return -1 * x
}

// The draw version of "v" function.
func (s *TrueSkill) vDraw(diff float64, drawMargin float64) float64 {
	absDiff := math.Abs(diff)
	a := drawMargin - absDiff
	b := -1*drawMargin - absDiff

	g := gaussian.NewGaussian(0.0, 1.0)

	denom := g.Cdf(a) - g.Cdf(b)
	numer := g.Pdf(b) - g.Pdf(a)

	if denom != 0 && !math.IsNaN(denom) {
		if diff < 0 {
			return (numer / denom) * -1
		}
		return (numer / denom)
	}
	if diff < 0 {
		return a * -1
	}
	return a
}

// The non-draw version of "W" function.
// "W" calculates a variation of a standard deviation.
func (s *TrueSkill) wWin(diff float64, drawMargin float64) float64 {
	x := diff - drawMargin
	v := s.vWin(diff, drawMargin)
	w := v * (v + x)
	if w > 0 && w < 1 {
		return w
	}

	panic(fmt.Sprintf("wWin floating point error w=%v", w))
}

// The draw version of "w" function.
func (s *TrueSkill) wDraw(diff float64, drawMargin float64) float64 {
	absDiff := math.Abs(diff)
	a := drawMargin - absDiff
	b := -1*drawMargin - absDiff

	g := gaussian.NewGaussian(0.0, 1.0)

	denom := g.Cdf(a) - g.Cdf(b)

	if denom == 0 || math.IsNaN(denom) {
		panic(fmt.Sprintf("wWin floating point error denom=%v", denom))
	}

	v := s.vDraw(absDiff, drawMargin)

	return math.Pow(v, 2) + (a*g.Pdf(a)-b*g.Pdf(b))/denom
}

// Makes a size map of each teams.
func teamSizes(ratingGroups [][]*Rating) []int {
	teamSizes := make([]int, 0, len(ratingGroups))
	size := 0
	for _, r := range ratingGroups {
		teamSizes = append(teamSizes, size+len(r))
		size += len(r)
	}

	return teamSizes
}
