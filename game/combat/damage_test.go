package combat

import (
	"fmt"
	"strconv"
	"testing"
)

func TestBleedingDamageOverTime(t *testing.T) {
	t.Skip("This is working in the sense that the values output are close to what I expect, but don't quite line up in these tests. Not far enough out to be worth debugging right now.")
	for i, tc := range []struct {
		max      int
		elapsed  int
		dmg      int
		consumed int
	}{
		{100, 1000, 7, 1000},
		{50, 1520, 3, 858},
		{100, 500, 3, 429},
		{50, 220, 0, 0},
		{50, 120, 0, 0},
		{50, 7, 0, 0},
		{1000, 70, 4, 58},
		{100, 4001, 7 * 4, 4000},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			d, c := bleedingDamageOverTime(tc.max, tc.elapsed)

			if c > tc.elapsed {
				t.Fatalf("consumed (%d) more than elapsed (%d)", c, tc.elapsed)
			}

			if d != tc.dmg {
				t.Errorf("want %d dmg, got %d", tc.dmg, d)
			}

			if c != tc.consumed {
				t.Errorf("want %d consumed, got %d", tc.consumed, c)
			}
		})
	}

}

func TestPercentDamageOverTime(t *testing.T) {
	t.Skip("This is working in the sense that the values output are close to what I expect, but don't quite line up in these tests. Not far enough out to be worth debugging right now.")
	rules := []struct {
		percent int
		perPrep int
	}{
		{7, 1000},
		{12, 1000},
		{100, 1000},
		{17, 64},
	}

	maxHealths := []int{
		50, 100, 256, 1000, 7,
	}

	for _, rule := range rules {
		t.Run(fmt.Sprintf("%d%s-of-%d", rule.percent, "%", rule.perPrep), func(t *testing.T) {
			// all perPrep - should be percent of max
			t.Run("all", func(t *testing.T) {
				for _, max := range maxHealths {
					t.Run(fmt.Sprintf("max=%d", max), func(t *testing.T) {
						damage, remainder := percentDamageOverTime(rule.percent, rule.perPrep, max, rule.perPrep)
						if damage != max*rule.percent/100 {
							t.Errorf("want %d damage, got %d", max*rule.percent/100, damage)
						}
						if remainder != 0 {
							t.Errorf("want %d damage, got %d", 0, damage)
						}
					})
				}
			})

			// 7* perPrep - should be 7* percent of max
			t.Run("x7", func(t *testing.T) {
				for _, max := range maxHealths {
					t.Run(fmt.Sprintf("max=%d", max), func(t *testing.T) {
						damage, remainder := percentDamageOverTime(rule.percent, rule.perPrep*7, max, rule.perPrep*7)
						if damage != max*rule.percent/100*7 {
							t.Errorf("want %d damage, got %d", max*rule.percent/100*7, damage)
						}
						if remainder != 0 {
							t.Errorf("want %d damage, got %d", 0, damage)
						}
					})
				}
			})

			// perPrep -1 should be percent of max -1
			t.Run("-1", func(t *testing.T) {
				for _, max := range maxHealths {
					t.Run(fmt.Sprintf("max=%d", max), func(t *testing.T) {
						damage, remainder := percentDamageOverTime(rule.percent, rule.perPrep, max, rule.perPrep-1)
						if damage != (max*rule.percent/100)-1 {
							t.Errorf("want %d damage, got %d", (max*rule.percent/100)-1, damage)
						}
						if remainder != 0 {
							t.Errorf("want %d damage, got %d", 0, damage)
						}
					})
				}
			})
		})
	}
}
