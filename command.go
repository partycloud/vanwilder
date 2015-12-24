package vanwilder

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

type Command struct {
	Command string
	CmdArgs string `json:"cmd-args"`
	Id      string
	Game    string
	Volume  string
}

func (c *Command) Execute() {
	switch c.Command {
	case "start-game":
		c.StartGame()
	case "stop-game":
		c.StopGame()
	}
}

func (c *Command) StartGame() {
	var port docker.Port
	port = "25565/tcp"
	portBinding := []docker.PortBinding{docker.PortBinding{HostPort: "25565"}}
	portBindings := make(map[docker.Port][]docker.PortBinding)
	portBindings[port] = portBinding
	hostConfig := &docker.HostConfig{
		Binds:        []string{c.Volume + ":/data:rw"},
		PortBindings: portBindings,
		VolumeDriver: "convoy",
	}

	container, err := c.createContainer()
	if err != nil {
		panic(err)
	}

	client, _ := docker.NewClientFromEnv()
	err = client.StartContainer(container.ID, hostConfig)
	if err != nil {
		panic(err)
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

func (c *Command) StopGame() {
	client, _ := docker.NewClientFromEnv()

	log.WithFields(log.Fields{
		"ID": c.Id,
	}).Info("stopping")
	err := client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            c.Id,
		RemoveVolumes: false,
		Force:         true,
	})
	if err != nil {
		panic(err)
	}
}

func (c *Command) createContainer() (*docker.Container, error) {
	client, _ := docker.NewClientFromEnv()

	createOptions := docker.CreateContainerOptions{
		Name: "game-1",
		Config: &docker.Config{
			AttachStdout: true,
			Cmd:          []string{c.CmdArgs},
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

// func removeContainerByName(name string) err {
//   client, _ := docker.NewClientFromEnv()
//
//   client.ListContainers(docker.ListContainersOptions{
//
//   })
//
//   client.RemoveContainer(docker.RemoveContainerOptions{
//     ID: container.ID,
//   })
// }
