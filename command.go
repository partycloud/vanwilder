package vanwilder

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

type Status struct {
	ID    string
	State string
}

type Command struct {
	Command string
	CmdArgs string `json:"cmd-args"`
	Game    string
	ID      string
	Ports   []docker.Port
	Volume  string
}

func (c *Command) Execute(events chan Status) {
	switch c.Command {
	case "start-game":
		c.StartGame(events)
	case "stop-game":
		c.StopGame(events)
	}
}

func (c *Command) StartGame(events chan Status) {
	portBindings := make(map[docker.Port][]docker.PortBinding, len(c.Ports))

	for _, guest := range c.Ports {
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

	if c.Volume != "" {
		hostConfig.VolumeDriver = "convoy"
		hostConfig.Binds = []string{c.Volume + ":/data:rw"}
	}

	container, err := c.createContainer()
	if err != nil {
		panic(err)
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
		panic(err)
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
}

func (c *Command) StopGame(events chan Status) {
	client, _ := docker.NewClientFromEnv()

	log.WithFields(log.Fields{
		"container": c.ID,
	}).Info("stopping")
	events <- Status{
		ID:    c.ID,
		State: "stopping",
	}

	err := client.StopContainer(c.ID, 60)
	if err != nil {
		panic(err)
	}

	err = client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            c.ID,
		RemoveVolumes: false,
		Force:         true,
	})
	if err != nil {
		panic(err)
	}

	log.WithFields(log.Fields{
		"container": c.ID,
	}).Info("stopped")
	events <- Status{
		ID:    c.ID,
		State: "stopped",
	}
}

func (c *Command) createContainer() (*docker.Container, error) {
	client, _ := docker.NewClientFromEnv()

	createOptions := docker.CreateContainerOptions{
		Name: "game-1",
		Config: &docker.Config{
			AttachStdout: true,
			Cmd:          strings.Split(c.CmdArgs, " "),
			Image:        c.Game,
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
