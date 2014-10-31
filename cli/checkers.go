package main

import (
	"net"
	"sync"
)

type FreeChecker interface {
	Subject() string
	IsFree() bool
}

type HTTPStatusResponse struct {
	Site       string
	StatusCode int
}

func (res HTTPStatusResponse) IsFree() bool {
	return res.StatusCode == 404
}

func (res HTTPStatusResponse) Subject() string {
	return res.Site
}

type DNSStatusResponse struct {
	DomainName string
	IP         *net.IPAddr
}

func (res DNSStatusResponse) IsFree() bool {
	return res.IP == nil
}

func (res DNSStatusResponse) Subject() string {
	return res.DomainName
}

type NameChecker struct {
	Name         string
	chanCheckers chan FreeChecker
	wg           *sync.WaitGroup
	FreeCheckers []FreeChecker
}
