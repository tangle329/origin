package netutils

import (
	"net"
	"strings"
	"testing"
)

func TestConversion(t *testing.T) {
	ip := net.ParseIP("10.1.2.3")
	if ip == nil {
		t.Fatal("Failed to parse IP")
	}

	u := IPToUint32(ip)
	t.Log(u)
	ip2 := Uint32ToIP(u)
	t.Log(ip2)

	if !ip2.Equal(ip) {
		t.Fatal("Conversion back and forth failed")
	}
}

func TestGenerateGateway(t *testing.T) {
	sna, err := NewSubnetAllocator("10.1.0.0/16", 8, nil)
	if err != nil {
		t.Fatal("Failed to initialize IP allocator: ", err)
	}

	sn, err := sna.GetNetwork()
	if err != nil {
		t.Fatal("Failed to get network: ", err)
	}
	if sn.String() != "10.1.0.0/24" {
		t.Fatalf("Did not get expected subnet (sn=%s)", sn.String())
	}

	gatewayIP := GenerateDefaultGateway(sn)
	if gatewayIP.String() != "10.1.0.1" {
		t.Fatalf("Did not get expected gateway IP Address (gatewayIP=%s)", gatewayIP.String())
	}
}

func TestParseCIDRMask(t *testing.T) {
	tests := []struct {
		cidr       string
		fixedShort string
		fixedLong  string
	}{
		{
			cidr: "192.168.0.0/16",
		},
		{
			cidr: "192.168.1.0/24",
		},
		{
			cidr: "192.168.1.1/32",
		},
		{
			cidr:       "192.168.1.0/16",
			fixedShort: "192.168.0.0/16",
			fixedLong:  "192.168.1.0/32",
		},
		{
			cidr:       "192.168.1.1/24",
			fixedShort: "192.168.1.0/24",
			fixedLong:  "192.168.1.1/32",
		},
	}

	for _, test := range tests {
		_, err := ParseCIDRMask(test.cidr)
		if test.fixedShort == "" && test.fixedLong == "" {
			if err != nil {
				t.Fatalf("unexpected error parsing CIDR mask %q: %v", test.cidr, err)
			}
		} else {
			if err == nil {
				t.Fatalf("unexpected lack of error parsing CIDR mask %q", test.cidr)
			}
			if !strings.Contains(err.Error(), test.fixedShort) {
				t.Fatalf("error does not contain expected string %q: %v", test.fixedShort, err)
			}
			if !strings.Contains(err.Error(), test.fixedLong) {
				t.Fatalf("error does not contain expected string %q: %v", test.fixedLong, err)
			}
		}
	}
}
