package service

import (
	"errors"
	"fmt"
	"net"
	"github.com/shirou/gopsutil/disk"
	"os"
)

func ExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	var ip net.IP
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
		}
	}
	if ip != nil {
		return ip.String(), nil
	} else {
		return "", errors.New("confirm you have no connection to the network")
	}
}

func FileServerUsage ()(usage DiskUsageObject, err error) {
	filePath := os.Getenv("DOGFOOTER_DATA")
	fmt.Fprintf(os.Stderr, "filePath: %v\n", filePath)
	result, err2 := disk.Usage(filePath)
	if err2 != nil {
		err = fmt.Errorf("FileServerUsage err: %v", err2)
		return
	}
	fmt.Fprintf(os.Stderr, "result:%v\n", result)

	usage.Available = result.Free
	usage.Total = result.Total
	usage.Used = result.Used
	usage.UsedPercent = result.UsedPercent

	return
}