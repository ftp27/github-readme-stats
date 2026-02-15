package common

import "math"

func exponentialCDF(x float64) float64 {
	return 1 - math.Pow(2, -x)
}

func logNormalCDF(x float64) float64 {
	return x / (1 + x)
}

type Rank struct {
	Level      string
	Percentile float64
}

func CalculateRank(allCommits bool, commits int, prs int, issues int, reviews int, stars int, followers int) Rank {
	commitsMedian := 250.0
	if allCommits {
		commitsMedian = 1000
	}

	prsMedian := 50.0
	issuesMedian := 25.0
	reviewsMedian := 2.0
	starsMedian := 50.0
	followersMedian := 10.0

	commitsWeight := 2.0
	prsWeight := 3.0
	issuesWeight := 1.0
	reviewsWeight := 1.0
	starsWeight := 4.0
	followersWeight := 1.0

	totalWeight := commitsWeight + prsWeight + issuesWeight + reviewsWeight + starsWeight + followersWeight

	thresholds := []float64{1, 12.5, 25, 37.5, 50, 62.5, 75, 87.5, 100}
	levels := []string{"S", "A+", "A", "A-", "B+", "B", "B-", "C+", "C"}

	rank := 1 - (commitsWeight*exponentialCDF(float64(commits)/commitsMedian)+
		prsWeight*exponentialCDF(float64(prs)/prsMedian)+
		issuesWeight*exponentialCDF(float64(issues)/issuesMedian)+
		reviewsWeight*exponentialCDF(float64(reviews)/reviewsMedian)+
		starsWeight*logNormalCDF(float64(stars)/starsMedian)+
		followersWeight*logNormalCDF(float64(followers)/followersMedian))/totalWeight

	level := levels[len(levels)-1]
	for idx, threshold := range thresholds {
		if rank*100 <= threshold {
			level = levels[idx]
			break
		}
	}

	return Rank{Level: level, Percentile: rank * 100}
}
