package repo

import (
	"github.com/ogra1/fabrica/domain"
	"os"
	"path"
)

// List returns a list of the builds that have been requested
func (bld *BuildService) List() ([]domain.Build, error) {
	return bld.Datastore.BuildList()
}

// BuildGet returns a build with its logs
func (bld *BuildService) BuildGet(id string) (domain.Build, error) {
	return bld.Datastore.BuildGet(id)
}

// BuildDelete deletes a build with its logs and snap
func (bld *BuildService) BuildDelete(id string) error {
	// Get the path to the built file and remove it
	build, err := bld.Datastore.BuildGet(id)
	if err != nil {
		return err
	}

	// Remove the stored files
	if build.Download != "" {
		dir := path.Dir(build.Download)
		os.RemoveAll(dir)
	}

	// Stop and delete the running container
	if build.Container != "" {
		bld.LXDSrv.StopAndDeleteContainer(build.Container)
	}

	// Remove the database records for the build and logs
	return bld.Datastore.BuildDelete(id)
}
