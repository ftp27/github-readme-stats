package common

import "strings"

type SecondaryErrorMessages struct {
	MaxRetry             string
	NoTokens             string
	UserNotFound         string
	GraphQLError         string
	GithubRestError      string
	WakatimeUserNotFound string
}

const TryAgainLater = "Please try again later"

var SecondaryErrors = SecondaryErrorMessages{
	MaxRetry:             "You can deploy own instance or wait until public will be no longer limited",
	NoTokens:             "Please add an env variable called PAT_1 with your GitHub API token in vercel",
	UserNotFound:         "Make sure the provided username is not an organization",
	GraphQLError:         TryAgainLater,
	GithubRestError:      TryAgainLater,
	WakatimeUserNotFound: "Make sure you have a public WakaTime profile",
}

type CustomError struct {
	Message          string
	Type             string
	SecondaryMessage string
}

func (e *CustomError) Error() string {
	return e.Message
}

const (
	ErrorMaxRetry     = "MAX_RETRY"
	ErrorNoTokens     = "NO_TOKENS"
	ErrorUserNotFound = "USER_NOT_FOUND"
	ErrorGraphQL      = "GRAPHQL_ERROR"
	ErrorGithubRest   = "GITHUB_REST_API_ERROR"
	ErrorWakatime     = "WAKATIME_ERROR"
)

func NewCustomError(message string, errType string) *CustomError {
	secondary := errType
	switch errType {
	case ErrorMaxRetry:
		secondary = SecondaryErrors.MaxRetry
	case ErrorNoTokens:
		secondary = SecondaryErrors.NoTokens
	case ErrorUserNotFound:
		secondary = SecondaryErrors.UserNotFound
	case ErrorGraphQL:
		secondary = SecondaryErrors.GraphQLError
	case ErrorGithubRest:
		secondary = SecondaryErrors.GithubRestError
	case ErrorWakatime:
		secondary = SecondaryErrors.WakatimeUserNotFound
	}
	return &CustomError{Message: message, Type: errType, SecondaryMessage: secondary}
}

type MissingParamError struct {
	MissedParams     []string
	SecondaryMessage string
}

func (e *MissingParamError) Error() string {
	if len(e.MissedParams) == 0 {
		return "Missing params make sure you pass the parameters in URL"
	}
	quoted := []string{}
	for _, param := range e.MissedParams {
		quoted = append(quoted, "\""+param+"\"")
	}
	return "Missing params " + strings.Join(quoted, ", ") + " make sure you pass the parameters in URL"
}

func NewMissingParamError(params []string, secondary string) *MissingParamError {
	return &MissingParamError{MissedParams: params, SecondaryMessage: secondary}
}

func RetrieveSecondaryMessage(err error) string {
	if err == nil {
		return ""
	}
	if custom, ok := err.(*CustomError); ok {
		return custom.SecondaryMessage
	}
	if missing, ok := err.(*MissingParamError); ok {
		return missing.SecondaryMessage
	}
	return ""
}
