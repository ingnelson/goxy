package main

import (
	"flag"
	goxy "github.com/scotow/goxy/server2"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

var (
	localHTTP = flag.String("h", ":8080", "local HTTP address used by Goxy clients (address:port)")
	remoteTCP = flag.String("r", "localhost:22", "remote TCP address (host:port)")
)

func main() {
	flag.Parse()

	localHTTPAddr, err := net.ResolveTCPAddr("tcp", *localHTTP)
	if err != nil {
		log.WithFields(log.Fields{
			"address": *localHTTP,
			"error":   err,
		}).Error("Invalid listening address.")
	}

	remoteTCPAddr, err := net.ResolveTCPAddr("tcp", *remoteTCP)
	if err != nil {
		log.WithFields(log.Fields{
			"address": *remoteTCP,
			"error":   err,
		}).Error("Invalid remote TCP address.")
	}

	listener, err := goxy.NewListener(localHTTPAddr)
	if err != nil {
		log.WithFields(log.Fields{
			"address": *localHTTP,
			"error":   err,
		}).Error("Cannot start HTTP listener.")
	}

	go listener.Start()

	log.WithFields(log.Fields{
		"local":  *localHTTP,
		"remote": *remoteTCP,
	}).Info("Goxy server started.")

	for {
		goxyConn, err := listener.Accept()
		if err != nil {
			log.WithFields(log.Fields{
				"address": *localHTTP,
				"error":   err,
			}).Warning("Cannot accept Goxy connection.")
		}

		log.WithFields(log.Fields{
			"local":  goxyConn.LocalAddr(),
			"remote": goxyConn.RemoteAddr(),
		}).Info("Goxy connection accepted.")

		tcpConn, err := net.DialTCP("tcp", nil, remoteTCPAddr)
		if err != nil {
			log.WithFields(log.Fields{
				"address": *remoteTCP,
				"error":   err,
			}).Error("Cannot open TCP connection.")
		}

		log.WithFields(log.Fields{
			"address": *remoteTCP,
		}).Info("TCP connection created.")

		go io.Copy(tcpConn, goxyConn)
		go io.Copy(goxyConn, tcpConn)
	}
}
