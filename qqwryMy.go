package main

import (
	"flag"
	"log"
	"net"
	"strings"

	"gopkg.in/olahol/melody.v1"
)

func init_qqwry() {
	datFile := flag.String("qqwry", "./qqwry.dat", "纯真 IP 库的地址")
	IPData.FilePath = *datFile
	res := IPData.InitIPData()

	if v, ok := res.(error); ok {
		log.Panic(v)
	}
}

func getIpAddr(s *melody.Session) (string, string) {
	ip, _, _ := net.SplitHostPort(strings.TrimSpace(s.Request.RemoteAddr))

	forwarded := s.Request.Header.Get("X-Forwarded-For")
	real := s.Request.Header.Get("X-Real-IP")

	if real != "" {
		ip = real
	} else if forwarded != "" {
		ip = forwarded
	}
	qqWry := NewQQwry()
	resultQQwry := qqWry.Find(ip)
	ipAddr := resultQQwry.Country
	return ip, ipAddr
}
