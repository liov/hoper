package osi

import "os"

// Hostname returns the host name reported by the kernel.
func Hostname() string {
	hostname, _ := os.Hostname()
	return hostname
}
