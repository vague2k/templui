package ctxkeys

type contextKey string

const (
	URLPathValue = contextKey("url_path_value")
	Version      = contextKey("version")
	GitHubStars  = contextKey("github_stars")
)
