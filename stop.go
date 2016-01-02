package vanwilder

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

type Stop struct{}

func (s *Stop) Execute(r *Req, events chan interface{}) error {
	client, _ := docker.NewClientFromEnv()

	log.WithFields(log.Fields{
		"container": r.ID,
	}).Info("stopping")
	events <- Status{
		ID:    r.ID,
		State: "stopping",
	}

	err := client.StopContainer(r.ID, 60)
	if err != nil {
		return err
	}

	err = client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            r.ID,
		RemoveVolumes: false,
		Force:         true,
	})
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"container": r.ID,
	}).Info("stopped")
	events <- Status{
		ID:    r.ID,
		State: "stopped",
	}

	return nil
}
