package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/partycloud/vanwilder"
	"github.com/shirou/gopsutil/mem"
)

var sockFile = "/var/run/van.sock"

func main() {
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
	signal.Notify(sigs, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Printf("Caught signal %s: shutting down.\n", sig)
		done <- true
	}()

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				panic(err)
			}
			dec := json.NewDecoder(conn)
			for {
				var cmd vanwilder.Command
				if err := dec.Decode(&cmd); err == io.EOF {
					break
				} else if err != nil {
					panic(err)
				}

				go cmd.Execute()
			}
			conn.Close()
		}
	}()

	<-done

	// // websocket
	// conn, err := dialWebSocket("ws://api.partycloud.dev:3000/vw?token1234")
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Println("connection established")
	//
	// for {
	// 	_, r, err := conn.NextReader()
	//
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	evtString, err := ioutil.ReadAll(r)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	fmt.Println(string(evtString))
	// }
}

func dialWebSocket(url string) (*websocket.Conn, error) {
	var dialer *websocket.Dialer
	var dialHeader http.Header

	conn, _, err := dialer.Dial(url, dialHeader)

	return conn, err
}
