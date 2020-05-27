package web

import (
	"net/http"
)

type imageAlias struct {
	Alias     string `json:"alias"`
	Available bool   `json:"available"`
}

var aliases = []string{"fabrica-bionic", "fabrica-xenial"}

// ImageAliases checks if the image aliases are available
func (srv Web) ImageAliases(w http.ResponseWriter, r *http.Request) {
	responseList := []imageAlias{}

	for _, a := range aliases {
		// Check if the alias is available
		err := srv.LXDSrv.GetImageAlias(a)

		responseList = append(responseList, imageAlias{
			Alias:     a,
			Available: err == nil,
		})
	}

	formatRecordsResponse(responseList, w)
}
