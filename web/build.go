package web

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type buildRequest struct {
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
	KeyID  string `json:"keyId"`
}

// Build initiates a build
func (srv Web) Build(w http.ResponseWriter, r *http.Request) {
	req := srv.decodeBuildRequest(w, r)
	if req == nil {
		return
	}

	// Start the build for the repo with this ID
	buildID, err := srv.BuildSrv.Build(req.Repo)
	if err != nil {
		formatStandardResponse("build", err.Error(), w)
		return
	}

	formatStandardResponse("", buildID, w)
}

// BuildList lists the build requests
func (srv Web) BuildList(w http.ResponseWriter, r *http.Request) {
	builds, err := srv.BuildSrv.List()
	if err != nil {
		formatStandardResponse("list", err.Error(), w)
		return
	}

	formatRecordsResponse(builds, w)
}

// BuildDelete deletes a build with its logs
func (srv Web) BuildDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := srv.BuildSrv.BuildDelete(vars["id"])
	if err != nil {
		formatStandardResponse("logs", err.Error(), w)
		return
	}

	formatStandardResponse("", "", w)
}

func (srv Web) decodeBuildRequest(w http.ResponseWriter, r *http.Request) *buildRequest {
	// Decode the JSON body
	req := buildRequest{}
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
