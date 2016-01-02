package vanwilder

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

type Start struct {
}

func (s *Start) Execute(r *Req, events chan interface{}) error {
	portBindings := make(map[docker.Port][]docker.PortBinding, len(r.Ports))

	for _, guest := range r.Ports {
		parts := strings.Split(string(guest), "/")
		portBindings[guest] = []docker.PortBinding{
			docker.PortBinding{
				HostPort: parts[0],
			},
		}
	}

	hostConfig := &docker.HostConfig{
		PortBindings: portBindings,
	}

	if r.Volume != "" {
		hostConfig.VolumeDriver = "convoy"
		hostConfig.Binds = []string{r.Volume + ":/data:rw"}
	}

	container, err := s.createContainer(r.Game, r.CmdArgs)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"ID": container.ID,
	}).Info("starting")
	events <- Status{
		ID:    container.ID,
		State: "starting",
	}

	client, _ := docker.NewClientFromEnv()
	err = client.StartContainer(container.ID, hostConfig)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"ID": container.ID,
	}).Info("started")
	events <- Status{
		ID:    container.ID,
		State: "started",
	}

	// fmt.Println(container)

	// err = client.Logs(docker.LogsOptions{
	// 	Follow:       true,
	// 	Container:    container.ID,
	// 	OutputStream: os.Stdout,
	// 	Stdout:       true,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	return nil
}

func (s *Start) createContainer(game, cmdArgs string) (*docker.Container, error) {
	client, _ := docker.NewClientFromEnv()

	createOptions := docker.CreateContainerOptions{
		Name: "game-1",
		Config: &docker.Config{
			AttachStdout: true,
			Cmd:          strings.Split(cmdArgs, " "),
			Image:        game,
		},
	}

	container, err := client.CreateContainer(createOptions)
	if err == docker.ErrContainerAlreadyExists {
		log.Info("removing existing container")
		err = client.RemoveContainer(docker.RemoveContainerOptions{
			ID: createOptions.Name,
		})
		if err != nil {
			return nil, err
		}

		container, err = client.CreateContainer(createOptions)
		if err != nil {
			return nil, err
		}

	} else if err != nil {
		return nil, err
	}

	return container, nil
}
