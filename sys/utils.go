package sys

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

const SERVICE_HOST = "SERVICE_HOST"

func GetExecFilename() string {
	_, filename, _, _ := runtime.Caller(1)
	cwdPath, _ := os.Getwd()
	return cwdPath + filename
}

func GetHost() string {
	ret := getIpFromEnv()
	if len(ret) > 0 {
		return ret
	}

	ret = getIpFromSpecialName("eth0", "em1")
	if len(ret) > 0 {
		return ret
	}

	ret = getFirstIp()
	if len(ret) > 0 {
		return ret
	}

	panic(fmt.Errorf("ip not found"))
	return ""
}

func getIpFromEnv() string {
	return os.Getenv(SERVICE_HOST)
}

func getIpFromSpecialName(name ...string) string {
	ifaces, e := net.Interfaces()
	if e != nil {
		panic(e)
	}

	match := func(n string) bool {
		for _, v := range name {
			if v == n {
				return true
			}
		}
		return false
	}

	for _, v := range ifaces {
		if match(v.Name) {
			return _getIpByFace(v)
		}
	}
	return ""
}

func _getIpByFace(iface net.Interface) string {
	if iface.Flags&net.FlagUp == 0 {
		return ""
	}

	if iface.Flags&net.FlagLoopback != 0 {
		return ""
	}

	// ignore docker and warden bridge
	if strings.HasPrefix(iface.Name, "docker") || strings.HasPrefix(iface.Name, "w-") {
		return ""
	}

	addrs, e := iface.Addrs()
	if e != nil {
		return ""
	}

	for _, addr := range addrs {
		var ip net.IP
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
		return ip.String()
	}
	return ""
}

func getFirstIp() string {
	ips := make([]string, 0)
	ifaces, e := net.Interfaces()
	if e != nil {
		panic(e)
	}
	for _, iface := range ifaces {
		ipStr := _getIpByFace(iface)
		if len(ipStr) > 0 {
			ips = append(ips, ipStr)
		}
	}
	if len(ips) > 0 {
		return ips[0]
	}
	return ""
}

func GetFreePort() (int, error) {
	ports, err := GetFreePorts(1)
	if err != nil {
		return 0, err
	}
	if len(ports) == 0 {
		return 0, fmt.Errorf("no useful port")
	}
	return ports[0], nil
}

func GetFreePorts(count int) ([]int, error) {
	var ports []int
	for i := 0; i < count; i++ {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			return nil, err
		}

		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			return nil, err
		}
		defer l.Close()
		ports = append(ports, l.Addr().(*net.TCPAddr).Port)
	}
	return ports, nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func InterruptHandler(errc chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	terminateError := fmt.Errorf("%s", <-c)
	errc <- terminateError
}

func GetFileList(directory string, ext string) []string {
	ext = strings.ToLower(ext)
	arr := make([]string, 0)
	err := filepath.Walk(directory, func(filename string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if len(ext) == 0 || ext == path.Ext(filename) {
			arr = append(arr, filename)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
	return arr
}
