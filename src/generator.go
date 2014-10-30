package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	cnt := 0
	domains := []string{}
	for cnt < 10 {
		domains = append(domains, findAvailableDomainName())
		cnt++
	}
	fmt.Print("Here's 10 available domains ! ")
	for i := 0; i < len(domains); i++ {
		fmt.Print(domains[i] + " ")
	}
}

func findAvailableDomainName() string {
	available := false
	result := ""
	for !available {
		consonne := []string{"h", "k", "h", "n", "m", "y", "d", "r"}
		voyelle := []string{"a", "e", "i", "o", "u"}
		s := constructString(consonne, voyelle)
		fmt.Print("testing " + s + "... ")
		available = executeHostCommand(s + ".com")
		if available {
			fmt.Println("available !")
		} else {
			fmt.Println("not available.")
		}
		result = s
	}
	return result
}

func constructString(array1 []string, array2 []string) string {
	s := []string{}
	for i := 0; i < 3; i++ {
		c1 := pickRandomLetter(array1)
		c2 := pickRandomLetter(array2)
		s = append(s, c1)
		s = append(s, c2)
	}
	str := strings.Join(s, "")
	return str
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func pickRandomLetter(array []string) string {
	return array[randInt(1, len(array))]
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
