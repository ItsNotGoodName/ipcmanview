package build

var (
	commit  = ""
	date    = ""
	version = "dev"
	repoURL = ""
)

func init() {
	Current = Build{
		Commit:     commit,
		Version:    version,
		Date:       date,
		RepoURL:    repoURL,
		CommitURL:  repoURL + "/tree/" + commit,
		LicenseURL: repoURL + "/blob/master/LICENSE",
		ReleaseURL: repoURL + "/releases/tag/" + version,
	}
	if repoURL == "" {
		Current.CommitURL = "#"
		Current.LicenseURL = "#"
		Current.ReleaseURL = "#"
	}
}

var Current Build

type Build struct {
	Commit     string
	Version    string
	Date       string
	RepoURL    string
	CommitURL  string
	LicenseURL string
	ReleaseURL string
}
