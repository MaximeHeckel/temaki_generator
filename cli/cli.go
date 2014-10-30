package main

import (
	"fmt"
	"net"
	"net/http"
	"runtime"
	"sync"

	"github.com/MaximeHeckel/temaki_generator/generator"
)

var (
	consonants = []rune{'h', 'k', 'h', 'n', 'm', 'y', 'd', 'r'}
	vowels     = []rune{'a', 'e', 'i', 'o', 'u'}
)

var sites = []string{
	"https://github.com/%s",
	"https://twitter.com/%s",
	"https://facebook.com/%s",
}

var domains chan string
var httpDomains chan string
var dnsDomains chan string

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	httpDomains = make(chan string, 1)
	dnsDomains = make(chan string, 1)
	domains = make(chan string)

	go func() {
		for domain := range domains {
			httpDomains <- domain
			dnsDomains <- domain
		}
		close(httpDomains)
		close(dnsDomains)
	}()

	go func() {
		for i := 0; i < 100; i++ {
			domains <- generator.ConstructString(consonants, vowels) + ".com"
		}

		close(domains)
	}()
}

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

func main() {
	wg := &sync.WaitGroup{}

	chanRes := make(chan FreeChecker)

	wg.Add(2)
	go checkHTTPSites(wg, chanRes)
	go checkDNSDomains(wg, chanRes)

	go func() {
		wg.Wait()
		close(chanRes)
	}()

	for checker := range chanRes {
		if checker.IsFree() {
			fmt.Println(checker.Subject(), "is free")
		} else {
			fmt.Println(checker.Subject(), "is occupied")
		}
	}
}

func checkHTTPSites(wg *sync.WaitGroup, chanRes chan<- FreeChecker) {
	defer wg.Done()

	for domain := range httpDomains {
		for _, site := range sites {
			wg.Add(1)
			go func(site string) {
				defer wg.Done()
				url := fmt.Sprintf(site, domain)
				res, err := http.Get(url)
				if err != nil {
					fmt.Printf("error %s %s %s\n", site, domain, err)
					return
				}

				chanRes <- HTTPStatusResponse{Site: url, StatusCode: res.StatusCode}
				res.Body.Close()
			}(site)
		}
	}
}

func checkDNSDomains(wg *sync.WaitGroup, chanRes chan<- FreeChecker) {
	defer wg.Done()

	for domain := range dnsDomains {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			ip, err := net.ResolveIPAddr("ip", domain)
			if err != nil {
				chanRes <- DNSStatusResponse{domain, nil}
			} else {
				chanRes <- DNSStatusResponse{domain, ip}
			}
		}(domain)
	}
}
