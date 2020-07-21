package geom

import (
	"fmt"
	"testing"
)

func BenchmarkDistance(b *testing.B) {
	testcases := []*Hex{
		&Hex{x: 0, y: 18},
		&Hex{x: 1, y: 1},
		&Hex{x: 2400, y: 180},
		&Hex{x: -108, y: 1027},
		&Hex{x: 24321, y: 273},
		&Hex{x: 273, y: 24321},
		&Hex{x: -730, y: -5200},
		&Hex{x: -18, y: 103},
		&Hex{x: 12, y: -1200},
		&Hex{x: 8000, y: 9000},
	}
	var total float64
	for i := 0; i < b.N; i++ {
		for first := 0; first < len(testcases); first++ {
			for second := 0; second < len(testcases); second++ {
				if first == second {
					continue
				}
				total += Distance(testcases[first], testcases[second])
			}
		}
	}
	fmt.Println(total)
}

func BenchmarkDistanceSquared(b *testing.B) {
	testcases := []*Hex{
		&Hex{x: 0, y: 18},
		&Hex{x: 1, y: 1},
		&Hex{x: 2400, y: 180},
		&Hex{x: -108, y: 1027},
		&Hex{x: 24321, y: 273},
		&Hex{x: 273, y: 24321},
		&Hex{x: -730, y: -5200},
		&Hex{x: -18, y: 103},
		&Hex{x: 12, y: -1200},
		&Hex{x: 8000, y: 9000},
	}
	var total float64
	for i := 0; i < b.N; i++ {
		for first := 0; first < len(testcases); first++ {
			for second := 0; second < len(testcases); second++ {
				if first == second {
					continue
				}
				total += DistanceSquared(testcases[first], testcases[second])
			}
		}
	}
	fmt.Println(total)
}
