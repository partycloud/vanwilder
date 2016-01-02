package vanwilder

import (
	"errors"

	"github.com/fsouza/go-dockerclient"
)

type Status struct {
	ID    string
	State string
}

type Cmd interface {
	Execute(r *Req, events chan interface{}) error
}

type Req struct {
	Command string
	CmdArgs string `json:"cmd-args"`
	Game    string
	ID      string
	Ports   []docker.Port
	Volume  string
}

var (
	ErrCommandNotFound = errors.New("command not found")
)

func (r *Req) Process(events chan interface{}) {
	if err := r.process(events); err != nil {
		// TODO: error handling
		panic(err)
	}
}

func (r *Req) process(events chan interface{}) error {
	cmd, err := findCommand(r.Command)
	if err != nil {
		return err
	}

	return cmd.Execute(r, events)
}

func findCommand(name string) (Cmd, error) {
	switch name {
	case "start-game":
		return &Start{}, nil
	case "stop-game":
		return &Stop{}, nil
	case "list-volumes":
		return &Volumes{}, nil
	}

	return nil, ErrCommandNotFound
}
