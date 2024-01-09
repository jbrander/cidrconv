package cidr

import (
	"fmt"
	"log"
	"math/big"
	"net"
)

func LogIPNet(network *net.IPNet) {
	ip := network.IP
	mask := network.Mask
	maskStr := fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3])

	log.Printf("  ip:\t%15s\t\t%08b %08b %08b %08b", ip, ip[0], ip[1], ip[2], ip[3])
	log.Printf("mask:\t%15s\t\t%08b %08b %08b %08b\n\n", maskStr, mask[0], mask[1], mask[2], mask[3])

}

func Subnet(base *net.IPNet, newBits int, smlNum int) (*net.IPNet, error) {
	num := big.NewInt(int64(smlNum))

	ip := base.IP
	mask := base.Mask
	LogIPNet(base)

	parentLen, addrLen := mask.Size()
	newPrefixLen := parentLen + newBits

	if newPrefixLen > addrLen {
		return nil, fmt.Errorf("insufficient address space to extend prefix of %d by %d", parentLen, newBits)
	}

	maxNetNum := uint64(1<<uint64(newBits)) - 1
	if num.Uint64() > maxNetNum {
		return nil, fmt.Errorf("prefix extension of %d does not accommodate a subnet numbered %d", newBits, num)
	}

	return &net.IPNet{
		IP:   insertNumIntoIP(ip, num, newPrefixLen),
		Mask: net.CIDRMask(newPrefixLen, addrLen),
	}, nil
}
