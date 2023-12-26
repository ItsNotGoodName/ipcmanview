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
		Date:       date,
		Version:    version,
		RepoURL:    repoURL,
		CommitURL:  repoURL + "/tree/" + commit,
		LicenseURL: repoURL + "/blob/master/LICENSE",
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
}
