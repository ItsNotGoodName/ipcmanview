package build

type Build struct {
	BuiltBy    string
	Commit     string
	Date       string
	Version    string
	RepoURL    string
	ReleaseURL string
}

func (b Build) CommitURL() string {
	return b.RepoURL + "/tree/" + b.Commit
}

func (b Build) LicenseURL() string {
	return b.RepoURL + "/blob/master/LICENSE"
}

var Current Build
