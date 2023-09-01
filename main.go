package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	//"github.com/apparentlymart/go-cidr/cidr"
)

var (
	verbose    bool
	parentCidr string
	subnetCidr string
)

func ipToInt(ip net.IP) *big.Int {
	val := &big.Int{}
	val.SetBytes([]byte(ip))
	return val
}

func printIPNet(network *net.IPNet) {
	ip := network.IP
	mask := network.Mask
	maskStr := fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3])

	fmt.Printf("  ip:\t%15s\t\t%08b %08b %08b %08b\n", ip, ip[0], ip[1], ip[2], ip[3])
	fmt.Printf("mask:\t%15s\t\t%08b %08b %08b %08b\n\n", maskStr, mask[0], mask[1], mask[2], mask[3])
}

func cidrsubnetSyntax(parentCidr string, subnetCidr string) (string, error) {
	_, parentNetwork, err := net.ParseCIDR(parentCidr)

	if err != nil {
		return "", fmt.Errorf("parsing parent CDIR failed %w", err)
	}

	_, subNetwork, err := net.ParseCIDR(subnetCidr)

	if err != nil {
		return "", fmt.Errorf("parsing subnet CIDR failed %w", err)
	}

	if verbose {
		fmt.Printf("Network (%s)\n", parentCidr)
		printIPNet(parentNetwork)
		fmt.Printf("Subnetwork (%s)\n", subnetCidr)
		printIPNet(subNetwork)
	}

	parentMaskLen, _ := parentNetwork.Mask.Size()
	subMaskLen, _ := subNetwork.Mask.Size()
	newbits := subMaskLen - parentMaskLen

	parentIP := ipToInt(parentNetwork.IP)
	subIP := ipToInt(subNetwork.IP)

	subIP.Xor(subIP, parentIP)
	bitsToShift := uint(32) - uint(subMaskLen)
	netnum := subIP.Rsh(subIP, bitsToShift)

	cidrsubnet := fmt.Sprintf("cidrsubnet(\"%s\", %d, %d)", parentCidr, newbits, netnum)

	return cidrsubnet, nil
}

func init() {

	if len(os.Args) < 5 {
		flag.PrintDefaults()
		log.Fatal("Missing required flags")
	}

	flag.StringVar(&parentCidr, "parent", "", "parent CIDR")
	flag.StringVar(&subnetCidr, "subnet", "", "subnet CIDR")
	flag.BoolVar(&verbose, "verbose", false, "display colorized output")
	flag.Parse()

}

func main() {
	cidrsubnet, err := cidrsubnetSyntax(parentCidr, subnetCidr)
	if err != nil {
		log.Panic(err)
	}

	fmt.Print(cidrsubnet)
}
