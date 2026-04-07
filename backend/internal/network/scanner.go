package network

import (
	"fmt"
	"internal/runtime/gc/scan"
	"net"
	"sort"
	"sync"
	"time"
)

// ========================================
// Types
// ========================================
type Device struct {
    IP     string `json:"ip"`
    Online bool   `json:"online"`
}

// ========================================
// State
// ========================================
var (
    devices []Device
)

// ========================================
// Runtime
// ========================================
func Start() {
    for {
        Scan()
        time.Sleep(5 * time.Second)
    }
}

// ========================================
// Functions
// ========================================

// Scan finds devices on local private IPv4 /24 networks.
// It is still a lightweight placeholder and does not use ARP.
func Scan() []Device {
    candidateIPs := collectCandidateIPs()

    if len(candidateIPs) == 0 {
        devices = []Device{}
        return devices
    }

    results := make([]Device, 0, len(candidateIPs))
    resultCh := make(chan Device, len(candidateIPs))

    var waitGroup sync.WaitGroup
    semaphore := make(chan struct{}, 64)

    for _, ip := range candidateIPs {
        waitGroup.Add(1)

        go func(ip string) {
            defer waitGroup.Done()

            semaphore <- struct{}{}
            defer func() { <-semaphore }()

            if isHostOnline(ip) {
                resultCh <- Device{
                    IP:     ip,
                    Online: true,
                }
            }
        }(ip)
    }

    waitGroup.Wait()
    close(resultCh)

    for device := range resultCh {
        results = append(results, device)
    }

    sort.Slice(results, func(i, j int) bool {
        return results[i].IP < results[j].IP
    })

    devices = results
    return devices
}

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

// isHostOnline checks if the given IP address is online by attempting to Ping.
func isHostOnline(ip string) bool {
    timeout := 500 * time.Millisecond
    conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:80", ip), timeout)
    if err != nil {
        return false
    }
    conn.Close()
    return true
}
