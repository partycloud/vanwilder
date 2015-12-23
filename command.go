package vanwilder

import (
	"fmt"
	"os"

	"github.com/fsouza/go-dockerclient"
)

type Command struct {
	Command string
	Game    string
	CmdArgs string `json:"cmd-args"`
	Volume  string
}

func (c *Command) Execute() {
	switch c.Command {
	case "start-game":
		c.StartGame()
	}
}

func (c *Command) StartGame() {
	client, _ := docker.NewClientFromEnv()

	container, err := client.CreateContainer(docker.CreateContainerOptions{
		Name: "game-1",
		Config: &docker.Config{
			AttachStdout: true,
			Cmd:          []string{c.CmdArgs},
			Image:        c.Game,
		},
	})
	if err != nil {
		panic(err)
	}

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

	err = client.StartContainer(container.ID, hostConfig)
	if err != nil {
		panic(err)
	}

	fmt.Println(container)

	err = client.Logs(docker.LogsOptions{
		Follow:       true,
		Container:    container.ID,
		OutputStream: os.Stdout,
		Stdout:       true,
	})
	if err != nil {
		panic(err)
	}
}
