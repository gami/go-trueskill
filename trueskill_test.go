package trueskill_test

import (
	"fmt"

	"github.com/gami/go-trueskill"
)

func ExampleTrueSkill_Rate() {

	ts := trueskill.NewTrueSkill(
		trueskill.MU(250),
		trueskill.Beta(125),
	)

	r1 := ts.CreateRating()
	r2 := ts.CreateRating()

	rs, err := ts.Rate1v1(r1, r2)
	if err != nil {
		panic(err)
	}

	for i, r := range rs {
		fmt.Printf("team=%v mu=%v sigma=%v score=%f\n", i, r.Mu, r.Sigma, (ts.Expose(r) + 2000))
	}

	rs, err = ts.Rate1v1(rs[0], rs[1])
	if err != nil {
		panic(err)
	}

	for i, r := range rs {
		fmt.Printf("team=%v mu=%v sigma=%v score=%f\n", i, r.Mu, r.Sigma, (ts.Expose(r) + 2000))
	}

	rs, err = ts.Rate1v1(rs[0], rs[1])
	if err != nil {
		panic(err)
	}

	for i, r := range rs {
		fmt.Printf("team=%v mu=%v sigma=%v score=%f\n", i, r.Mu, r.Sigma, (ts.Expose(r) + 2000))
	}

	rs, err = ts.Rate1v1(rs[0], rs[1])
	if err != nil {
		panic(err)
	}

	for i, r := range rs {
		fmt.Printf("team=%v mu=%v sigma=%v score=%f\n", i, r.Mu, r.Sigma, (ts.Expose(r) + 2000))
	}

	// Output:
	// team=0 member=0 score=7.881405
	// team=1 member=0 score=-0.910259
}
