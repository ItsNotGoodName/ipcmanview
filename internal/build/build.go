package build

var (
	Version string
	Commit  string
	Date    string
	RepoURL string
)

var (
	CommitURL  string
	LicenseURL string
)

func init() {
	CommitURL = RepoURL + "/tree/" + Commit
	LicenseURL = RepoURL + "/blob/master/LICENSE"
}
