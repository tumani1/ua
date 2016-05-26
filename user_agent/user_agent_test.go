package user_agent

import (
	"os"
	"log"
	"bufio"
	"testing"
	"os/exec"
	"path/filepath"
)

var (
	benchfile string = filepath.Join("data", "ua_warmup.log")
	datafile string = filepath.Join("data", "ua.log")

	testData []string
	randomUA string = getRandomUA()

	ua *UserAgent
)


func init() {
	// read bench file to the array
	readDataToMemory()

	ua, _ = New()
}

func getRandomUA() string {
	//Init benchmarck file
	//brew install coreutils
	//gshuf -n 20000 uniq_ua.txt > uniq_ua_bench.txt

	cmd := "gshuf"
	args := []string{"-n", "1", datafile}
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		log.Println(os.Stderr, err)
		return ""
	}

	return string(out)
}

func readDataToMemory() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	benchfile := filepath.Join(pwd, benchfile)

	file, err := os.Open(benchfile)
	if err != nil {
		log.Fatal("Error:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		testData = append(testData, scanner.Text())
	}
}

func BenchmarkUaParse(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ua.parseUA(randomUA)
	}
}

func BenchmarkUaParseWg(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ua.parseUAWg(randomUA)
	}
}

func BenchmarkParseUAWithLRU(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ua.parseUAWithLRU(randomUA)
	}
}