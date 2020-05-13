package web

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"path"
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

// BuildDownload fetches the built snap
func (srv Web) BuildDownload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bld, err := srv.BuildSrv.BuildGet(vars["id"])
	if err != nil {
		formatStandardResponse("logs", err.Error(), w)
		return
	}

	// Get the filename of the download
	filename := path.Base(bld.Download)

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	download, err := os.Open(bld.Download)
	if err != nil {
		formatStandardResponse("download", err.Error(), w)
		return
	}
	defer download.Close()

	io.Copy(w, download)
}
