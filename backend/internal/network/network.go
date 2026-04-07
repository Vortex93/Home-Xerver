package network

import (
	"app/internal/global"
	"app/internal/log"
	"app/internal/utils/arraylist"
	"app/internal/utils/ssh"

	"fmt"
	"net"
	"time"
)

// ========================================
// State
// ========================================
var (
	SURVEILLANCE_IPS = []string{
		"10.0.0.70",
		"10.0.0.71",
		"10.0.0.72",
		"10.0.0.73",
		"10.0.0.74",
		"10.0.0.75",
		"10.0.0.76",
		"10.0.0.78",
		"10.0.0.79",
	}

	ROUTER_IPS = []string{
		"10.0.0.90",
		"10.0.0.91",
		"10.0.0.92",
	}
)
var (
	logger *log.Logger
)

// ========================================
// Bootstrap
// ========================================
func init() {
	var err error
	logger, err = log.Create("network")
	if err != nil {
		panic("Failed to create network logger: " + err.Error())
	}
}

// ========================================
// Functions
// ========================================

// ========================================
// Helpers
// ========================================

// collectCandidateIPs gathers candidate IP addresses from local private IPv4 /24 networks.
func collectCandidateIPs() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	seen := make(map[string]struct{})

	for _, networkInterface := range interfaces {
		if networkInterface.Flags&net.FlagUp == 0 {
			continue
		}

		if networkInterface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addresses, err := networkInterface.Addrs()
		if err != nil {
			continue
		}

		for _, address := range addresses {
			ipNet, ok := address.(*net.IPNet)
			if !ok {
				continue
			}

			ipv4 := ipNet.IP.To4()
			if ipv4 == nil {
				continue
			}

			if !isPrivateIPv4(ipv4) {
				continue
			}

			base := net.IPv4(ipv4[0], ipv4[1], ipv4[2], 0)

			for host := 1; host <= 254; host++ {
				candidate := net.IPv4(base[0], base[1], base[2], byte(host)).String()
				seen[candidate] = struct{}{}
			}
		}
	}

	results := make([]string, 0, len(seen))
	for ip := range seen {
		results = append(results, ip)
	}

	return results
}

// isPrivateIPv4 checks if the given IP address is a private IPv4 address.
func isPrivateIPv4(ip net.IP) bool {
	if len(ip) != 4 {
		return false
	}

	if ip[0] == 10 {
		return true
	}

	if ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31 {
		return true
	}

	if ip[0] == 192 && ip[1] == 168 {
		return true
	}

	return false
}

func checkPort22(ip string) bool {
	timeout := 20 * time.Millisecond
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, 22), timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func checkPort80(ip string) bool {
	timeout := 20 * time.Millisecond
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, 80), timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// ========================================
// Entrpoint
// ========================================
func Start() {
	logger.Info("Started network scanner")

	go func() {
		for {
			// Scan()
			time.Sleep(5 * time.Second)
		}
	}()

	go func() { // Camera 10.0.0.77
        var online bool = false
        var history = arraylist.NewWithLimit[bool](5)
        var lastRestart time.Time

	    for {
			online = checkPort80("10.0.0.77")
			history.Add(online)

			if history.IsFull() && history.All(func(value bool) bool { return value == false }) {
                logger.Info("Camera is offline")

                if time.Since(lastRestart) < 1*time.Minute {
                    logger.Info("Skipping reboot because it was rebooted recently")
                    continue
                }

                lastRestart = time.Now()

                // Restart router to reconnect the camera
                c := ssh.NewClient(ssh.ConnectParams{
                    Host:     "10.0.0.91",
                    Port:     22,
                    Username: global.GetDevice91User(),
                    Password: global.GetDevice91Pass(),
                })

                err := c.Connect()
                if err != nil {
                    logger.Error(fmt.Sprintf("Failed to connect to router for rebooting camera: %v", err))
                }

                _, err = c.Execute(ssh.CommandParams{
                    Command: "reboot",
                })
                if err != nil {
                    logger.Error(fmt.Sprintf("Failed to execute reboot command on router for camera: %v", err))
                }

                logger.Debugf("Executed reboot command on router to restart camera: %v", err)

                // Wait until port 80 is back online
                for {
                    if checkPort80("10.0.0.77") {
                        logger.Info("Camera is back online")
                        break
                    }

                    time.Sleep(1 * time.Second)
                }

                c.Close()
            }

			time.Sleep(1 * time.Second)
		}
	}()

}
