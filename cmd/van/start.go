package main

import (
	"encoding/json"
	"io"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/partycloud/vanwilder"
	"github.com/shirou/gopsutil/mem"
)

var StartDescription = `
Usage: van start [OPTIONS]

Start listening.

Example:
   $ van start --volume-driver convoy`

type VanStartConfig struct {
	VolumeDriver string `cli:"volume-driver"`
}

var VanStartCommand = cli.Command{
	Name:        "start",
	Usage:       "Starts the party",
	Description: StartDescription,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "volume-driver",
			Value:  "convoy",
			Usage:  "Docker volume driver (default is convoy)",
			EnvVar: "DOCKER_VOLUME_DRIVER",
		},
	},
	Action: func(c *cli.Context) {
		host, _ := os.Hostname()
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		v, _ := mem.VirtualMemory()

		log.WithFields(log.Fields{
			"host":      host,
			"cpus":      runtime.NumCPU(),
			"mem-total": v.Total / 1024 / 1024,
			"mem-free":  v.Free / 1024 / 1024,
		}).Info("environment")

		// socket
		if _, err := os.Stat(sockFile); err == nil {
			log.WithFields(log.Fields{"file": sockFile}).Info("remove-socket")
			if err := os.Remove(sockFile); err != nil {
				panic(err)
			}
		}

		l, err := net.Listen("unix", sockFile)
		if err != nil {
			panic(err)
		}
		defer l.Close()

		sigs := make(chan os.Signal, 1)
		done := make(chan bool, 1)
		events := make(chan vanwilder.Status, 100)

		signal.Notify(sigs, os.Interrupt, os.Kill, syscall.SIGTERM)
		go func() {
			sig := <-sigs
			log.WithField("sig", sig).Info("caught")
			done <- true
		}()

		// this should drain events even when not connected to a socket
		go func() {
			for {
				conn, err := l.Accept()
				if err != nil {
					panic(err)
				}
				commands := make(chan vanwilder.Command)
				go func() {
					for {
						dec := json.NewDecoder(conn)
						for {
							var cmd vanwilder.Command
							if err := dec.Decode(&cmd); err == io.EOF {
								break
							} else if err != nil {
								panic(err)
							}

							commands <- cmd
						}
					}
				}()

				for {
					select {
					case cmd := <-commands:
						go cmd.Execute(events)
					case event := <-events:
						b, _ := json.Marshal(event)
						conn.Write(b)
						conn.Write([]byte("\n"))
					}
				}
			}
		}()

		<-done

	},
}
