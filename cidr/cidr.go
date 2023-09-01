//code borrowed from github.com/apparentlymart/go-cidr

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

func ipToInt(ip net.IP) (*big.Int, int) {
	val := &big.Int{}
	val.SetBytes([]byte(ip))
	if len(ip) == net.IPv4len {
		return val, 32
	} else if len(ip) == net.IPv6len {
		return val, 128
	} else {
		panic(fmt.Errorf("unsupported address length %d", len(ip)))
	}
}

func insertNumIntoIP(ip net.IP, bigNum *big.Int, prefixLen int) net.IP {
	ipInt, totalBits := ipToInt(ip)
	bigNum.Lsh(bigNum, uint(totalBits-prefixLen))
	ipInt.Or(ipInt, bigNum)
	return intToIP(ipInt, totalBits)
}

func intToIP(ipInt *big.Int, bits int) net.IP {
	ipBytes := ipInt.Bytes()
	ret := make([]byte, bits/8)
	// Pack our IP bytes into the end of the return array,
	// since big.Int.Bytes() removes front zero padding.
	for i := 1; i <= len(ipBytes); i++ {
		ret[len(ret)-i] = ipBytes[len(ipBytes)-i]
	}
	return net.IP(ret)
}
