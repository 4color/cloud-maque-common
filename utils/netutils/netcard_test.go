package netutils

import "testing"

func Test_IP(t *testing.T) {

	ip := GetLocalIpV4()
	hostname := GetOsHostname()
	println(ip)
	println(hostname)
}
