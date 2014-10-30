package generator

import (
	"bytes"
	"math/rand"
	"os/exec"
	"time"
)

func Generate() string {
	rand.Seed(time.Now().UTC().UnixNano())
	domains := findAvailableDomainName()
	return domains
}

func findAvailableDomainName() string {
	consonants := []rune{'h', 'k', 'h', 'n', 'm', 'y', 'd', 'r'}
	vowels := []rune{'a', 'e', 'i', 'o', 'u'}
	domainName := ""
	available := false

	for !available {
		domainName = constructString(consonants, vowels)
		//fmt.Print("testing " + domainName+".com" + " ... ")
		available = executeHostCommand(domainName + ".com")
		/*if available {
			fmt.Println("available !")
		} else {
			fmt.Println("not available.")
		}*/
	}

	return domainName
}

func constructString(array1 []rune, array2 []rune) string {
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

func executeHostCommand(s string) bool {
	cmd := exec.Command("host", s)
	_, err := cmd.Output()
	if err != nil {
		return true
	} else {
		return false
	}
}
