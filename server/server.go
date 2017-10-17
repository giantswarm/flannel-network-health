package server

import (
	"fmt"
	"github.com/giantswarm/micrologger"
	"github.com/vishvananda/netlink"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Server struct {
	Logger           micrologger.Logger
	BridgeInterface  string
	FlannelInterface string
	BridgeIP         string
	FlannelIP        string
}

func DefaultConfig() *Server {
	return &Server{}
}

func (s *Server) LoadConfig() bool {
	// load NIC interfaces from ENV
	s.BridgeInterface = os.Getenv("NETWORK_BRIDGE_NAME")
	s.FlannelInterface = os.Getenv("NETWORK_FLANNEL_DEVICE")
	// read flannel file
	fileContent, err := ioutil.ReadFile(os.Getenv("NETWORK_ENV_FILE_PATH"))
	if err != nil {
		s.Logger.Log(fmt.Printf("Error reading flannel file. %v", err.Error()))
		return false
	}
	// get FLANNEL_SUBNET from flannel file via regexp
	r, _ := regexp.Compile("FLANNEL_SUBNET=[0-9]+.[0-9]+.[0-9]+.[0-9]+/[0-9]+")
	flannelLine := r.Find(fileContent)
	// check if regexp returned non-empty line
	if len(flannelLine) < 5 {
		s.Logger.Log(fmt.Print("Unable to find FLANNEL_SUBNET in flannel file"))
		return false
	}

	// parse flannel subnet
	flannelSubnetStr := strings.Split(string(flannelLine), "=")[1]
	flannelIP, _, err := net.ParseCIDR(flannelSubnetStr)
	if err != nil {
		s.Logger.Log(fmt.Printf("Error when parsing flannel subnet. %v", err.Error()))
		return false
	}
	// force ipv4 for later trick
	flannelIP = flannelIP.To4()

	// get flannel ip
	s.FlannelIP = flannelIP.String()
	// get bridge ip, which is just one number bigger than flannel hence the [3]++ trick
	flannelIP[3]++
	s.BridgeIP = flannelIP.String()
	// debug output
	s.Logger.Log(fmt.Printf("Loaded Config: %+v", s))

	return true
}

func (s *Server) CheckBridgeInterface(w http.ResponseWriter, r *http.Request) {
	var healthy bool = true
	// load interface
	bridge, err := netlink.LinkByName(s.BridgeInterface)
	if err != nil {
		healthy = false
		s.Logger.Log(fmt.Printf("Cant find bridge %s. %s", s.BridgeInterface, err))
	}
	// check ip on interface
	ipList, err := netlink.AddrList(bridge, netlink.FAMILY_V4)
	if err != nil || len(ipList) == 0 {
		healthy = false
		s.Logger.Log(fmt.Printf("Missing ip %s in the bridge configuration.", s.BridgeIP))
	}
	if len(ipList) > 0 && ipList[0].IP.String() != s.BridgeIP {
		healthy = false
		s.Logger.Log(fmt.Printf("Wrong ip on interface %s. Expected %s, but found %s.", s.BridgeInterface, s.BridgeIP, ipList[0].IP.String()))
	}

	// if health check failed set response status to 503
	if !healthy {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "FAILED")
	} else {
		fmt.Fprintln(w, "OK")
		s.Logger.Log(fmt.Printf("Healthcheck for bridge %s has been successful. Bridge is present and configured with ip %s.", s.BridgeInterface, s.BridgeIP))
	}
}

func (s *Server) CheckFlannelInterface(w http.ResponseWriter, r *http.Request) {
	var healthy bool = true
	// load interface
	flannel, err := netlink.LinkByName(s.FlannelInterface)
	if err != nil {
		healthy = false
		s.Logger.Log(fmt.Printf("Cant find flannel interface %s. %s", s.FlannelInterface, err))
	}
	// check ip on interface
	ipList, err := netlink.AddrList(flannel, netlink.FAMILY_V4)
	if err != nil || len(ipList) == 0 {
		healthy = false
		s.Logger.Log(fmt.Printf("Missing ip %s in the flannel configuration.", s.FlannelIP))
	}
	if len(ipList) > 0 && ipList[0].IP.String() != s.FlannelIP {
		healthy = false
		s.Logger.Log(fmt.Printf("Wrong ip on interface %s. Expected %s, but found %s.", s.FlannelInterface, s.FlannelIP, ipList[0].IP.String()))
	}

	// if health check failed set response status to 503
	if !healthy {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "FAILED")
	} else {
		fmt.Fprintln(w, "OK")
		s.Logger.Log(fmt.Printf("Healthcheck for flannel interface %s has been successful. Interface is present and configured with ip %s.", s.FlannelInterface, s.FlannelIP))
	}
}
