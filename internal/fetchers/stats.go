package fetchers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ftp27/github-readme-stats/internal/common"
	"github.com/ftp27/github-readme-stats/internal/github"
	"github.com/spf13/viper"
)

const graphQLReposField = `
  repositories(first: 100, ownerAffiliations: OWNER, orderBy: {direction: DESC, field: STARGAZERS}, after: $after) {
    totalCount
    nodes {
      name
      stargazers {
        totalCount
      }
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
`

const graphQLReposQuery = `
  query userInfo($login: String!, $after: String) {
    user(login: $login) {
      ` + graphQLReposField + `
    }
  }
`

const graphQLStatsQuery = `
  query userInfo($login: String!, $after: String, $includeMergedPullRequests: Boolean!, $includeDiscussions: Boolean!, $includeDiscussionsAnswers: Boolean!, $startTime: DateTime = null) {
    user(login: $login) {
      name
      login
      commits: contributionsCollection (from: $startTime) {
        totalCommitContributions,
      }
      reviews: contributionsCollection {
        totalPullRequestReviewContributions
      }
      repositoriesContributedTo(first: 1, contributionTypes: [COMMIT, ISSUE, PULL_REQUEST, REPOSITORY]) {
        totalCount
      }
      pullRequests(first: 1) {
        totalCount
      }
      mergedPullRequests: pullRequests(states: MERGED) @include(if: $includeMergedPullRequests) {
        totalCount
      }
      openIssues: issues(states: OPEN) {
        totalCount
      }
      closedIssues: issues(states: CLOSED) {
        totalCount
      }
      followers {
        totalCount
      }
      repositoryDiscussions @include(if: $includeDiscussions) {
        totalCount
      }
      repositoryDiscussionComments(onlyAnswers: true) @include(if: $includeDiscussionsAnswers) {
        totalCount
      }
      ` + graphQLReposField + `
    }
  }
`

type StatsData struct {
	Name                     string
	TotalStars               int
	TotalCommits             int
	TotalIssues              int
	TotalPRs                 int
	TotalPRsMerged           int
	MergedPRsPercentage      float64
	TotalReviews             int
	TotalDiscussionsStarted  int
	TotalDiscussionsAnswered int
	ContributedTo            int
	Rank                     common.Rank
}

type graphQLRepoNode struct {
	Name       string `json:"name"`
	Stargazers struct {
		TotalCount int `json:"totalCount"`
	} `json:"stargazers"`
}

type graphQLPageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor"`
}

type graphQLUser struct {
	Name    string `json:"name"`
	Login   string `json:"login"`
	Commits struct {
		TotalCommitContributions int `json:"totalCommitContributions"`
	} `json:"commits"`
	Reviews struct {
		TotalPullRequestReviewContributions int `json:"totalPullRequestReviewContributions"`
	} `json:"reviews"`
	RepositoriesContributedTo struct {
		TotalCount int `json:"totalCount"`
	} `json:"repositoriesContributedTo"`
	PullRequests struct {
		TotalCount int `json:"totalCount"`
	} `json:"pullRequests"`
	MergedPullRequests struct {
		TotalCount int `json:"totalCount"`
	} `json:"mergedPullRequests"`
	OpenIssues struct {
		TotalCount int `json:"totalCount"`
	} `json:"openIssues"`
	ClosedIssues struct {
		TotalCount int `json:"totalCount"`
	} `json:"closedIssues"`
	Followers struct {
		TotalCount int `json:"totalCount"`
	} `json:"followers"`
	RepositoryDiscussions struct {
		TotalCount int `json:"totalCount"`
	} `json:"repositoryDiscussions"`
	RepositoryDiscussionComments struct {
		TotalCount int `json:"totalCount"`
	} `json:"repositoryDiscussionComments"`
	Repositories struct {
		TotalCount int               `json:"totalCount"`
		Nodes      []graphQLRepoNode `json:"nodes"`
		PageInfo   graphQLPageInfo   `json:"pageInfo"`
	} `json:"repositories"`
}

type graphQLUserData struct {
	User graphQLUser `json:"user"`
}

func FetchStats(username string, includeAllCommits bool, excludeRepo []string, includeMergedPRs bool, includeDiscussions bool, includeDiscussionAnswers bool, commitsYear *int) (StatsData, error) {
	if username == "" {
		return StatsData{}, common.NewMissingParamError([]string{"username"}, "")
	}

	startTime := ""
	if commitsYear != nil {
		startTime = fmtDateYear(*commitsYear)
	}

	statsResponse, err := statsFetcher(username, includeMergedPRs, includeDiscussions, includeDiscussionAnswers, startTime)
	if err != nil {
		return StatsData{}, err
	}

	if len(statsResponse.Errors) > 0 {
		errType := statsResponse.Errors[0].Type
		if errType == "NOT_FOUND" {
			return StatsData{}, common.NewCustomError(statsResponse.Errors[0].Message, common.ErrorUserNotFound)
		}
		if statsResponse.Errors[0].Message != "" {
			return StatsData{}, common.NewCustomError(statsResponse.Errors[0].Message, common.ErrorGraphQL)
		}
		return StatsData{}, common.NewCustomError("Something went wrong while trying to retrieve the stats data using the GraphQL API.", common.ErrorGraphQL)
	}

	var data graphQLUserData
	if err := json.Unmarshal(statsResponse.Data, &data); err != nil {
		return StatsData{}, err
	}

	user := data.User
	stats := StatsData{
		Name: user.Name,
	}
	if stats.Name == "" {
		stats.Name = user.Login
	}

	if includeAllCommits {
		totalCommits, err := totalCommitsFetcher(username)
		if err != nil {
			return StatsData{}, err
		}
		stats.TotalCommits = totalCommits
	} else {
		stats.TotalCommits = user.Commits.TotalCommitContributions
	}

	stats.TotalPRs = user.PullRequests.TotalCount
	if includeMergedPRs {
		stats.TotalPRsMerged = user.MergedPullRequests.TotalCount
		if stats.TotalPRs > 0 {
			stats.MergedPRsPercentage = (float64(stats.TotalPRsMerged) / float64(stats.TotalPRs)) * 100
		}
	}

	stats.TotalReviews = user.Reviews.TotalPullRequestReviewContributions
	stats.TotalIssues = user.OpenIssues.TotalCount + user.ClosedIssues.TotalCount
	if includeDiscussions {
		stats.TotalDiscussionsStarted = user.RepositoryDiscussions.TotalCount
	}
	if includeDiscussionAnswers {
		stats.TotalDiscussionsAnswered = user.RepositoryDiscussionComments.TotalCount
	}
	stats.ContributedTo = user.RepositoriesContributedTo.TotalCount

	excluded := append([]string{}, excludeRepo...)
	excluded = append(excluded, common.ExcludeRepositories()...)
	excludedSet := map[string]bool{}
	for _, repo := range excluded {
		excludedSet[repo] = true
	}

	totalStars := 0
	for _, node := range user.Repositories.Nodes {
		if excludedSet[node.Name] {
			continue
		}
		totalStars += node.Stargazers.TotalCount
	}
	stats.TotalStars = totalStars

	stats.Rank = common.CalculateRank(
		includeAllCommits,
		stats.TotalCommits,
		stats.TotalPRs,
		stats.TotalIssues,
		stats.TotalReviews,
		stats.TotalStars,
		user.Followers.TotalCount,
	)

	return stats, nil
}

func statsFetcher(username string, includeMergedPRs bool, includeDiscussions bool, includeDiscussionAnswers bool, startTime string) (github.GraphQLResponse, error) {
	var statsResponse github.GraphQLResponse
	hasNextPage := true
	endCursor := ""

	for hasNextPage {
		variables := map[string]any{
			"login":                     username,
			"first":                     100,
			"after":                     nil,
			"includeMergedPullRequests": includeMergedPRs,
			"includeDiscussions":        includeDiscussions,
			"includeDiscussionsAnswers": includeDiscussionAnswers,
			"startTime":                 nil,
		}
		if endCursor != "" {
			variables["after"] = endCursor
		}
		if startTime != "" {
			variables["startTime"] = startTime
		}

		query := graphQLStatsQuery
		if endCursor != "" {
			query = graphQLReposQuery
		}

		resp, err := github.RetryGraphQL(func(token string) (github.GraphQLResponse, *github.RestError, error) {
			return github.DoGraphQL(query, variables, token)
		}, github.RetryConfig{IsRateLimited: isRateLimitedGraphQL})
		if err != nil {
			errMessage := strings.ToLower(err.Error())
			switch {
			case strings.Contains(errMessage, "no github api tokens found"):
				return github.GraphQLResponse{}, common.NewCustomError("No GitHub API tokens found", common.ErrorNoTokens)
			case strings.Contains(errMessage, "all github api tokens failed"):
				return github.GraphQLResponse{}, common.NewCustomError("Downtime due to GitHub API rate limiting", common.ErrorMaxRetry)
			default:
				return github.GraphQLResponse{}, err
			}
		}

		statsResponse = mergeGraphQL(statsResponse, resp)

		var data graphQLUserData
		if err := json.Unmarshal(statsResponse.Data, &data); err != nil {
			return github.GraphQLResponse{}, err
		}

		repoNodes := data.User.Repositories.Nodes
		if len(repoNodes) == 0 {
			hasNextPage = false
			break
		}

		nodesWithStars := 0
		for _, node := range repoNodes {
			if node.Stargazers.TotalCount != 0 {
				nodesWithStars++
			}
		}

		hasNextPage = strings.ToLower(viper.GetString("FETCH_MULTI_PAGE_STARS")) == "true" && len(repoNodes) == nodesWithStars && data.User.Repositories.PageInfo.HasNextPage
		endCursor = data.User.Repositories.PageInfo.EndCursor
	}

	return statsResponse, nil
}

func mergeGraphQL(base github.GraphQLResponse, next github.GraphQLResponse) github.GraphQLResponse {
	if len(base.Data) == 0 {
		return next
	}
	if len(next.Data) == 0 {
		return base
	}

	var baseData graphQLUserData
	var nextData graphQLUserData
	if err := json.Unmarshal(base.Data, &baseData); err != nil {
		return next
	}
	if err := json.Unmarshal(next.Data, &nextData); err != nil {
		return base
	}

	baseData.User.Repositories.Nodes = append(baseData.User.Repositories.Nodes, nextData.User.Repositories.Nodes...)
	merged, err := json.Marshal(baseData)
	if err != nil {
		return base
	}

	base.Data = merged
	return base
}

func isRateLimitedGraphQL(resp github.GraphQLResponse) bool {
	if len(resp.Errors) == 0 {
		return false
	}
	err := resp.Errors[0]
	if err.Type == "RATE_LIMITED" {
		return true
	}
	return strings.Contains(strings.ToLower(err.Message), "rate limit")
}

func totalCommitsFetcher(username string) (int, error) {
	resp, err := retryRest(username)
	if err != nil {
		return 0, err
	}
	if resp.TotalCount == 0 {
		return 0, common.NewCustomError("Could not fetch total commits.", common.ErrorGithubRest)
	}
	return resp.TotalCount, nil
}

func retryRest(username string) (github.RestResponse, error) {
	keys := github.GetPATKeys()
	if len(keys) == 0 {
		return github.RestResponse{}, common.NewCustomError("No GitHub API tokens found", common.ErrorNoTokens)
	}

	for _, key := range keys {
		token := viper.GetString(key)
		resp, restErr, _, err := github.SearchCommits(username, token)
		if err != nil {
			return github.RestResponse{}, err
		}

		if restErr != nil {
			msg := strings.ToLower(restErr.Message)
			if msg == "bad credentials" || msg == "sorry. your account was suspended." {
				continue
			}
			return github.RestResponse{}, errors.New(restErr.Message)
		}

		return resp, nil
	}

	return github.RestResponse{}, common.NewCustomError("Downtime due to GitHub API rate limiting", common.ErrorMaxRetry)
}

func fmtDateYear(year int) string {
	return fmt.Sprintf("%d-01-01T00:00:00Z", year)
}
