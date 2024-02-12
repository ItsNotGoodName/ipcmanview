package build

import "time"

var (
	commit  = ""
	date    = ""
	version = "dev"
	repoURL = ""
)

func init() {
	date, _ := time.Parse(time.RFC3339, date)

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
	Date       time.Time
	RepoURL    string
	CommitURL  string
	LicenseURL string
	ReleaseURL string
}
