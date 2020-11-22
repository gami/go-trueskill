package trueskill_test

import (
	"fmt"

	"github.com/gami/go-trueskill"
)

func ExampleTrueSkill_Rate() {

	ts := trueskill.NewTrueSkill()

	r1 := ts.CreateRating()
	r2 := ts.CreateRating()

	calced, err := ts.Rate([][]*trueskill.Rating{{r1}, {r2}})
	if err != nil {
		panic(err)
	}

	for i, rs := range calced {
		for k, r := range rs {
			fmt.Printf("team=%v member=%v score=%f\n", i, k, ts.Expose(r))
		}
	}

	// Output:
	// team=0 member=0 score=7.881405
	// team=1 member=0 score=-0.910259
}
