package apicupcfg

import (
	"errors"
	"fmt"
	"log"
	"math/bits"
	"strings"
)

func DotsToBits(dotmask string) (uint32, error) {

	// dotmask is exected as 4 decimal numbers a.b.c.d
	sbytes := strings.Split(dotmask, ".")
	if len(sbytes) != 4 {
		return 0, errors.New(fmt.Sprintf("input mask exected to have 4 dot separated bytes, found: %d, mask: %s",
			len(sbytes), dotmask))
	}

	var bitmask uint32

	for idx, sbyte := range sbytes {
		// convert string byte to int
		ibyte := 0
		_, err := fmt.Sscanf(strings.TrimSpace(sbyte), "%d", &ibyte)
		if err != nil {
			return 0, err
		}

		if ibyte < 0 || ibyte > 255 {
			return 0, errors.New(fmt.Sprintf("invalid mask %s, byte %d out of range 0-255", dotmask, idx))
		}

		bitmask |= uint32(ibyte) << ((4 - idx - 1) * 8)
	}

	return bitmask, nil
}

func DotmaskToNethost(dotmask string) (netbits uint32, hostbits uint32, err error) {
	netbits, err = DotsToBits(dotmask)
	if err != nil {
		return 0,0, err
	}

	return netbits, ^netbits, nil
}

func SplitAddress(dotaddress string, dotmask string) (netbits uint32, hostbits uint32, err error) {

	netmask, hostmask, err := DotmaskToNethost(dotmask)
	if err != nil {
		return 0,0, err
	}

	addrbits, err := DotsToBits(dotaddress)
	if err != nil {
		return 0,0, err
	}

	netbits = netmask & addrbits
	hostbits = hostmask & addrbits

	return netbits, hostbits, nil
}

func Byte1(ipbits uint32) uint32 {
	return (ipbits& 0xff000000) >> 24
}

func Byte2(ipbits uint32) uint32 {
	return (ipbits& 0x00ff0000) >> 16
}

func Byte3(ipbits uint32) uint32 {
	return (ipbits& 0x0000ff00) >> 8
}

func Byte4(ipbits uint32) uint32 {
	return ipbits& 0x000000ff
}

func BitsToDots(ipbits uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", Byte1(ipbits), Byte2(ipbits), Byte3(ipbits), Byte4(ipbits))
}

func DecodeAddress(dotip, dotgw, dotmask string) {

	netmask, hostmask, err := DotmaskToNethost(dotmask)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	netbitscount := bits.OnesCount32(netmask)
	hostbitscount := bits.OnesCount32(hostmask)

	// validate no off-bits in netmask (@todo)

	minip := uint32(0)
	bcastip := uint32((2 << (hostbitscount-1)) - 1)
	maxip := uint32(bcastip -1)

	//maxmask := (2 << 31) - 1 - bcastip

	fmt.Printf("net-mask: %s(%x, %d bits), host-mask: %s(%x, %d bits)\n",
		BitsToDots(netmask), netmask, netbitscount,
		BitsToDots(hostmask), hostmask, hostbitscount)

	fmt.Printf("min-ip: %s(%x), max-ip: %s(%x), bcast-ip: %s(%x)\n",
		BitsToDots(minip), minip,
		BitsToDots(maxip), maxip,
		BitsToDots(bcastip), bcastip)

	// use gw and mask to find subnet
	subnetbits, gwbits, err := SplitAddress(dotgw, dotmask)
	fmt.Printf("gw: %s, subnet: %s(%x), gw-ip: %s(%x)\n",
		dotgw,
		BitsToDots(subnetbits), subnetbits,
		BitsToDots(gwbits), gwbits)

	// ip under netmask must match the subnet
	ipbits, err := DotsToBits(dotip)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	subnetdiff := ipbits & netmask - subnetbits
	ipdiff := int(maxip - (ipbits & hostmask))

	subnetstatus := func() string {if subnetdiff == 0 {return "ok"} else {return "bad"}}()
	ipstatus := func() string {if ipdiff >= 0 {return "ok"} else {return "bad"}}()

	fmt.Printf("ip: %s, ip-subnet: %s(%x, subnet-diff: %x %s), ip-host: %s(%x, ip-diff: %x %s)\n",
		dotip,
		BitsToDots(ipbits & netmask), ipbits & netmask, subnetdiff, subnetstatus,
		BitsToDots(ipbits & hostmask), ipbits & hostmask, ipdiff, ipstatus)
}
