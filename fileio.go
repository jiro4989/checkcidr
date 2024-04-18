package main

import (
	"bufio"
	"net"
	"os"
	"strings"
)

func ReadCIDR(path string) ([]*net.IPNet, error) {
	cp, err := os.Open(path)
	defer cp.Close()
	if err != nil {
		return nil, err
	}

	cidr := make([]*net.IPNet, 0)
	sc := bufio.NewScanner(cp)
	for sc.Scan() {
		l := strings.TrimSpace(sc.Text())
		_, ip, err := net.ParseCIDR(l)
		if err != nil {
			return nil, err
		}
		cidr = append(cidr, ip)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	return cidr, nil
}
