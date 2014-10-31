package generator

import (
	"bytes"
	"math/rand"
	"os/exec"
	"time"
	//"net"
	//"fmt"
	//"sync"
)

rand.Seed(time.Now().UTC().UnixNano())

consonants := []rune{'h', 'k', 'h', 'n', 'm', 'y', 'd', 'r'}
vowels := []rune{'a', 'e', 'i', 'o', 'u'}

var sites = []string{
	"https://github.com/%s",
	"https://twitter.com/%s",
	"https://facebook.com/%s",
}

var domains chan string
var httpDomains chan string

func Generate() {
	go func(){
		for i := 0; i < 100; i++ {
			domains <- ConstructString(consonants, vowels)
		}
	}()

	go func(){
		for domain := range domains {
			httpDomains <- domain
		}
	}()
}

/*func findAvailableDomainName() string {
	domainName := ""
	available := false

	for !available {
		domainName = ConstructString(consonants, vowels)
		//available = executeHostCommand(domainName + ".com")
	}

	return domainName
}*/

func ConstructString(array1 []rune, array2 []rune) string {
	s := &bytes.Buffer{}
	for i := 0; i < 3; i++ {
		s.WriteRune(pickRandomLetter(array1))
		s.WriteRune(pickRandomLetter(array2))
	}
	return s.String()
}

func pickRandomLetter(array []rune) rune {
	return array[randInt(1, len(array))]
}
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

type Check interface {
	Name() string
	CheckIfFree() string
}

type HTTPStatusResponse struct {
	Site string
	StatusCode int
}

func (res HTTPStatusResponse) CheckIfFree() bool{
	return res.StatusCode == 404
}

func (res HTTPStatusResponse) Name() string{
	return res.Site
}

/*func checkSites(chanRes chan<- Check){
	for domain := range httpDomains {

	}
}*/

func executeHostCommand(s string) bool {
	cmd := exec.Command("host", s)
	_, err := cmd.Output()
	if err != nil {
		return true
	} else {
		return false
	}
}
