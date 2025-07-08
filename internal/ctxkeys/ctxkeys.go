package ctxkeys

type contextKey string

const (
	URLPathValue = contextKey("url_path_value")
	GitHubStars  = contextKey("github_stars")
)
