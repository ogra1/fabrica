package web

import (
	"github.com/ogra1/fabrica/domain"
	"net/http"
)

var aliases = []string{"fabrica-bionic", "fabrica-xenial"}

// ImageAliases checks if the image aliases are available
func (srv Web) ImageAliases(w http.ResponseWriter, r *http.Request) {
	responseList := []domain.SettingAvailable{}

	for _, a := range aliases {
		// Check if the alias is available
		err := srv.LXDSrv.GetImageAlias(a)

		responseList = append(responseList, domain.SettingAvailable{
			Name:      a,
			Available: err == nil,
		})
	}

	formatRecordsResponse(responseList, w)
}

// CheckConnections checks the snap interfaces are connected
func (srv Web) CheckConnections(w http.ResponseWriter, r *http.Request) {
	results := srv.LXDSrv.CheckConnections()

	formatRecordsResponse(results, w)
}
