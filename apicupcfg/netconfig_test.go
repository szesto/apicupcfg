package apicupcfg

import (
	"fmt"
	"testing"
)

func TestNetmask(t *testing.T) {
	dotmask := "255.255.252.0"
	bitmask, err := DotsToBits(dotmask)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%x\n", bitmask)

	expnetbits := 0xfffffc00
	if bitmask != uint32(expnetbits) {
		t.Errorf("DotsToBits(%s) = %x, expected %x\n", dotmask, bitmask, uint32(expnetbits))
	}
}

func TestDotmaskToNethost(t *testing.T) {
	dotmask := "255.255.252.0"
	netmask, hostmask, err := DotmaskToNethost(dotmask)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%x %x\n", netmask, hostmask)

	expnetbits := 0xfffffc00
	exphostbits := 0x3ff

	if netmask != uint32(expnetbits) || hostmask != uint32(exphostbits) {
		t.Errorf("DotmaskToNethost(%s) = %x, %x, expected %x, %x", dotmask, netmask, hostmask, expnetbits, exphostbits)
	}
}

func TestSplitAddress(t *testing.T) {
	dotmask := "255.255.255.0"
	dotaddress := "172.20.165.191"

	netbits, hostbits, err := SplitAddress(dotaddress, dotmask)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%x %x\n", netbits, hostbits)

	dotaddress = "172.20.164.1"
	netbits, hostbits, err = SplitAddress(dotaddress, dotmask)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%x %x\n", netbits, hostbits)
}

func TestDecodeAddress(t *testing.T) {
	dotmask := "255.255.252.0"
	dotaddress := "172.20.165.191"
	dotgw := "172.20.164.1"

	DecodeAddress(dotaddress, dotgw, dotmask)
}