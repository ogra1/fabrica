package web

import "net/http"

// RepoCreate initiates a build
func (srv Web) RepoCreate(w http.ResponseWriter, r *http.Request) {
	req := srv.decodeBuildRequest(w, r)
	if req == nil {
		return
	}

	// Create the repo
	repoID, err := srv.BuildSrv.RepoCreate(req.Repo)
	if err != nil {
		formatStandardResponse("repo", err.Error(), w)
		return
	}

	formatStandardResponse("", repoID, w)
}

// RepoList lists the watched repos
func (srv Web) RepoList(w http.ResponseWriter, r *http.Request) {
	records, err := srv.BuildSrv.RepoList()
	if err != nil {
		formatStandardResponse("list", err.Error(), w)
		return
	}

	formatRecordsResponse(records, w)
}
