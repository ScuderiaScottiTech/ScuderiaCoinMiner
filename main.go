package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/paulbellamy/ratecounter"
)

var (
	mineapi            = flag.String("api", "", "API endpoint of the mining server")
	minerid            = flag.String("id", "", "YOUR telegram ID")
	goroutines         = flag.Int("goroutines", 1, "Number of goroutines to mine onto")
	ratecounterenabled = flag.Bool("ratecounter", false, "Rate counter enabled (may degrade performance)")
	testMode           = flag.Bool("testmode", false, "Test your hash rate without the use of an API")
)

var hashCounter *ratecounter.RateCounter
var correctCounter *ratecounter.RateCounter

//
// MINER GOROUTINE
//
// MINER GOROUTINE
//

func MinerThread(result chan string, offset uint64) {
	fmt.Println("Starting with offset", offset)

	for {
		offset++
		magic := strconv.FormatUint(offset, 10)
		valid, _ := CheckChallenge(magic)

		if *ratecounterenabled {
			hashCounter.Incr(1)
		}

		if valid {
			correctCounter.Incr(1)
			result <- magic
		}
	}
}

func CounterPrinter() {
	for {
		time.Sleep(time.Second)
		fmt.Println("Hash rate: ", hashCounter.Rate(), "/ second")
		fmt.Println("Correct hash rate: ", correctCounter.Rate(), "/ hour")
	}
}

func main() {
	hashCounter = ratecounter.NewRateCounter(1 * time.Second)
	correctCounter = ratecounter.NewRateCounter(1 * time.Hour)

	rand.Seed(time.Now().UnixNano())

	flag.Parse()

	if (*mineapi == "" || *minerid == "") && !*testMode { // ignore checks if in test mode
		fmt.Println("--help for help")
		return
	}

	result := make(chan string)

	go PeriodicChallengeRefresher()
	RefreshCurrentChallenge()

	for i := 0; i < *goroutines; i++ {
		fmt.Println("Spawning goroutine")
		time.Sleep(1500 * time.Millisecond)

		offset := rand.Uint64()
		go MinerThread(result, offset)
	}

	if *ratecounterenabled {
		go CounterPrinter()
	}

	for {
		magicResult := <-result
		fmt.Println("Found a valid magic number, ", magicResult)
		PostChallengeResult(magicResult, *minerid)
		time.Sleep(1 * time.Second) // sleep to avoid 503
	}
}
