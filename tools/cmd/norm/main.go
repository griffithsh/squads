// command norm is an interface to the NormFloat64 method in the standard
// library rand package that is useful when designing bell curves.
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"
)

var (
	stdDev float64
	mean   float64
	bucket float64
)

func flags() {
	stdDevFlag := flag.Float64("std-dev", 1.4, "standard deviation to provide rand.NormFloat64")
	meanFlag := flag.Float64("mean", 8.0, "mean to provide rand.NormFloat64")
	bucketFlag := flag.Float64("bucket", 0.5, "how large the result buckets should be")
	flag.Parse()

	stdDev = *stdDevFlag
	mean = *meanFlag
	bucket = *bucketFlag
}

func fmtString(s string) string {
	s = fmt.Sprintf("%7v", s)
	s = s[0:7]
	return s
}

func main() {
	flags()

	result := map[float64]int{}
	rand.Seed(time.Now().Unix())

	for i := 0; i < 150; i++ {
		roll := rand.NormFloat64()*stdDev + mean

		// Coerce roll to the nearest bucket.
		coerce := math.Round(roll/bucket) * bucket
		result[coerce]++
	}

	// Sort results.
	sortedKeys := make([]float64, 0, len(result))
	for k := range result {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Slice(sortedKeys, func(i, j int) bool {
		return sortedKeys[i] < sortedKeys[j]
	})

	// Format results for printing.
	first := []string{}
	second := []string{}
	for _, key := range sortedKeys {
		first = append(first, fmtString(fmt.Sprintf("%f", key)))
		second = append(second, fmtString(fmt.Sprintf("%d", result[key])))
	}

	fmt.Printf("%s\n%s\n", strings.Join(first, "\t"), strings.Join(second, "\t"))
}
