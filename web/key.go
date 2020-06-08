package web

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ogra1/fabrica/domain"
	"io"
	"net/http"
)

// KeyCreate store a new ssh key
func (srv Web) KeyCreate(w http.ResponseWriter, r *http.Request) {
	req := srv.decodeKeyRequest(w, r)
	if req == nil {
		return
	}

	keyID, err := srv.KeySrv.Create(req.Name, req.Username, req.Data, req.Password)
	if err != nil {
		formatStandardResponse("key", err.Error(), w)
		return
	}

	formatStandardResponse("", keyID, w)
}

// KeyList lists the ssh keys
func (srv Web) KeyList(w http.ResponseWriter, r *http.Request) {
	records, err := srv.KeySrv.List()
	if err != nil {
		formatStandardResponse("list", err.Error(), w)
		return
	}

	formatRecordsResponse(records, w)
}

// KeyDelete removes an unused key
func (srv Web) KeyDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Delete the repo
	if err := srv.KeySrv.Delete(vars["id"]); err != nil {
		formatStandardResponse("key", err.Error(), w)
		return
	}

	formatStandardResponse("", "", w)
}

func (srv Web) decodeKeyRequest(w http.ResponseWriter, r *http.Request) *domain.Key {
	// Decode the JSON body
	req := domain.Key{}
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
