package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

// BuildLog fetches a build with its logs
func (srv Web) BuildLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bld, err := srv.BuildSrv.BuildGet(vars["id"])
	if err != nil {
		formatStandardResponse("logs", err.Error(), w)
		return
	}

	formatRecordResponse(bld, w)
}
