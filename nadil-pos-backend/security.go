package main

import (
	"net"
)

func getMacAddr() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {

		return nil, err
	}
	var nethardwares []string
	for _, iface := range interfaces {

		nethardwares = append(nethardwares, iface.HardwareAddr.String())

	}
	return nethardwares, nil
}
