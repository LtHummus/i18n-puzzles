package main

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/unicode/norm"

	"github.com/lthummus/i18n-puzzles/input"
)

// this should probably be protected by a mutex :)
var cache = map[string]string{}

func buildDatabase(in string) map[string][]byte {
	lines := strings.Split(in, "\n")
	ret := map[string][]byte{}
	for _, curr := range lines {
		parts := strings.Split(curr, " ")
		ret[parts[0]] = []byte(parts[1])
	}

	return ret
}

func getAllNorms(x string) []string {
	if utf8.RuneCountInString(x) == len(x) {
		// is all ascii? no variants possible
		return []string{x}
	}

	// i just realized this is potentially broken if we are given a completely non-normalized character (such as ậ =
	// U+0061 + U+0302 + U+0323). We got lucky here, but we should probably also return the form we got in the
	// result as well if it is not normalized
	nfc := norm.NFC.String(x)
	nfd := norm.NFD.String(x)

	// for some characters (like Ω) are the same in both NFC and NFD forms
	if nfc == nfd {
		return []string{nfc}
	} else {
		return []string{nfc, nfd}
	}
}

// generateAllNormalizations generates all possible unicode normalizations of a string. In this case, because we're checking
// passwords, we can't detect a correct password by normalizing what we have because all we have is a hash we can not reverse.
// so generate every potential combination of normalized rune in the string. This will produce 2^n strings as output,
// where n is the number of potentially denormalized runes
func generateAllNormalizations(x string) []string {
	var chars []string
	for _, curr := range x {
		chars = append(chars, string(curr))
	}

	// this generates all possible normalization variants....for example the string "brûlée" will output
	// [][]string{[]string{"b"}, []string{"r"}, []string{"û", "û"}, []string{"l"}, []string{"é", "é"}, []string{"e"}}
	// where each element of the slice is a slice of possible unicode normalizations for that character (note that above,
	// û and û appear the same, but they have different normalization forms
	var allVariants [][]string
	for _, curr := range chars {
		variants := getAllNorms(curr)
		allVariants = append(allVariants, variants)
	}

	return generateCombinations(allVariants, 0, "")
}

func generateCombinations(variants [][]string, idx int, curr string) []string {
	if idx == len(variants) {
		return []string{curr}
	}

	var res []string
	for _, variant := range variants[idx] {
		combos := generateCombinations(variants, idx+1, curr+variant)
		res = append(res, combos...)
	}
	return res
}

func validLogin(db map[string][]byte, entry string) bool {
	parts := strings.Split(entry, " ")

	username := parts[0]
	n := norm.NFC.String(parts[1])

	cached := cache[username]
	if cached == n {
		return true
	} else if cached != "" {
		return false
	}

	hash := db[username]
	if hash == nil {
		return false
	}

	passwordPotentials := generateAllNormalizations(n)

	for _, curr := range passwordPotentials {
		if err := bcrypt.CompareHashAndPassword(hash, []byte(curr)); err == nil {
			cache[username] = n
			return true
		}
	}

	return false
}

func workerRoutine(db map[string][]byte, jobs <-chan string, results chan<- bool) {
	for curr := range jobs {
		results <- validLogin(db, curr)
	}
}

func main() {
	in, err := input.GetInputUTF8(context.Background(), 10, input.RealInput)
	if err != nil {
		panic(err)
	}
	
	parts := strings.Split(in, "\n\n")

	db := buildDatabase(parts[0])

	attempts := strings.Split(parts[1], "\n")
	fmt.Printf("Read %d attempts\n", len(attempts))

	jobs := make(chan string, len(attempts))
	results := make(chan bool, len(attempts))

	numWorkers := runtime.NumCPU()
	for i := 0; i < numWorkers; i++ {
		go workerRoutine(db, jobs, results)
	}

	fmt.Printf("Spawned %d threads\n", numWorkers)

	start := time.Now()
	for _, curr := range attempts {
		jobs <- curr
	}
	close(jobs)

	valid := 0
	for range len(attempts) {
		if <-results {
			valid++
		}
	}

	fmt.Printf("%d valid attempts\n", valid)
	dur := time.Since(start)
	fmt.Printf("Took %.2f seconds\n", dur.Seconds())
}
