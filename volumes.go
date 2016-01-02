package vanwilder

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

type Volumes struct{}

func (v *Volumes) Execute(r *Req, events chan interface{}) error {
	client, _ := docker.NewClientFromEnv()

	vols, err := client.ListVolumes(docker.ListVolumesOptions{})
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"volumes": len(vols),
	}).Info("list-volumes")
	events <- vols

	return nil
}
