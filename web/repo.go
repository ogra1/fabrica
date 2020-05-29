package web

import (
	"encoding/json"
	"io"
	"net/http"
)

type repoDeleteRequest struct {
	RepoID       string `json:"id"`
	DeleteBuilds bool   `json:"deleteBuilds"`
}

// RepoCreate creates a repository
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
	records, err := srv.BuildSrv.RepoList(false)
	if err != nil {
		formatStandardResponse("list", err.Error(), w)
		return
	}

	formatRecordsResponse(records, w)
}

// RepoDelete remove a repo and, optionally, its builds
func (srv Web) RepoDelete(w http.ResponseWriter, r *http.Request) {
	req := srv.decodeRepoDeleteRequest(w, r)
	if req == nil {
		return
	}

	// Delete the repo
	if err := srv.BuildSrv.RepoDelete(req.RepoID, req.DeleteBuilds); err != nil {
		formatStandardResponse("repo", err.Error(), w)
		return
	}

	formatStandardResponse("", "", w)
}

func (srv Web) decodeRepoDeleteRequest(w http.ResponseWriter, r *http.Request) *repoDeleteRequest {
	// Decode the JSON body
	req := repoDeleteRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	switch {
	// Check we have some data
	case err == io.EOF:
		formatStandardResponse("data", "No request data supplied.", w)
		return nil
		// Check for parsing errors
	case err != nil:
		formatStandardResponse("decode-json", err.Error(), w)
		return nil
	}
	return &req
}
