package web

import (
	"github.com/ogra1/fabrica/domain"
	"net/http"
)

// SystemResources monitors the system resources
func (srv Web) SystemResources(w http.ResponseWriter, r *http.Request) {
	cpu, err := srv.SystemSrv.CPU()
	if err != nil {
		formatStandardResponse("system", err.Error(), w)
		return
	}
	mem, err := srv.SystemSrv.Memory()
	if err != nil {
		formatStandardResponse("system", err.Error(), w)
		return
	}

	rec := domain.SystemResources{
		CPU:    cpu,
		Memory: mem,
	}
	formatRecordResponse(rec, w)
}
