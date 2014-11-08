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

var namesGenerator chan string

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	t := NewThrottler(10)
	namesGenerator = make(chan string)
	go func() {
		for {
			t.Acquire()
			namesGenerator <- generator.ConstructString(consonants, vowels)
		}
		close(namesGenerator)
	}()

	// Check availability of .com, twitter, facebook and github
	mainWg := &sync.WaitGroup{}
	for name := range namesGenerator {
		nameChecker := &NameChecker{
			Name:         name,
			wg:           &sync.WaitGroup{},
			chanCheckers: make(chan FreeChecker),
		}
		go nameChecker.check(mainWg, t)
	}
	mainWg.Wait()
}

func (nameChecker *NameChecker) check(mainWg *sync.WaitGroup, t *Throttler) {
	defer mainWg.Done()

	mainWg.Add(1)
	nameChecker.wg.Add(2)
	go nameChecker.CheckDNS()
	go nameChecker.CheckHTTP()

	go func() {
		nameChecker.wg.Wait()
		close(nameChecker.chanCheckers)
	}()

	for freeChecker := range nameChecker.chanCheckers {
		nameChecker.FreeCheckers = append(nameChecker.FreeCheckers, freeChecker)
	}

	fmt.Println(nameChecker.Name+" ? ", nameChecker.IsFree())
	t.Release()
}

func (nameChecker *NameChecker) CheckDNS() {
	defer nameChecker.wg.Done()

	// Add DNS FreeChecker to chanCheckers
	ip, err := net.ResolveIPAddr("ip", nameChecker.Name+".com")
	if err != nil {
		nameChecker.chanCheckers <- DNSStatusResponse{nameChecker.Name, nil}
	} else {
		nameChecker.chanCheckers <- DNSStatusResponse{nameChecker.Name, ip}
	}
}

func (nameChecker *NameChecker) CheckHTTP() {
	defer nameChecker.wg.Done()

	httpWg := &sync.WaitGroup{}

	// Add HTTP FreeCheckers to chanCheckers
	for _, site := range sites {
		httpWg.Add(1)
		go func(site string) {
			defer httpWg.Done()
			url := fmt.Sprintf(site, nameChecker.Name)
			res, err := http.Get(url)
			if err != nil {
				fmt.Printf("error %s %s %s\n", site, nameChecker.Name, err)
				return
			}
			nameChecker.chanCheckers <- HTTPStatusResponse{Site: url, StatusCode: res.StatusCode}
			res.Body.Close()
		}(site)
	}

	httpWg.Wait()
}

func (nameChecker *NameChecker) IsFree() bool {
	for _, freeChecker := range nameChecker.FreeCheckers {
		if !freeChecker.IsFree() {
			return false
		}
	}
	return true
}
