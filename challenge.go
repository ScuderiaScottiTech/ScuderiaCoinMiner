package main

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ChallengeApiResponse struct {
	ChallengeString string `json:"challenge_current"`
	Reward          int    `json:"challenge_reward"`
	Difficulty      int    `json:"challenge_difficulty"`
}

var currentChallenge ChallengeApiResponse

func GetMiningData() ChallengeApiResponse {
	resp, err := http.Get(*mineapi + "/mine/getChallenge")
	if err != nil {
		panic("Error in GetMiningData:\n" + err.Error())
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	apiResponse := &ChallengeApiResponse{}
	json.Unmarshal(body, apiResponse)

	return *apiResponse
}

func GetFakeMiningData() ChallengeApiResponse {
	return ChallengeApiResponse{
		ChallengeString: "123test123test",
		Reward:          100,
		Difficulty:      7,
	}
}

func RefreshCurrentChallenge() {
	fmt.Print("Getting new challenge info: ")

	if *testMode {
		currentChallenge = GetFakeMiningData()
	} else {
		currentChallenge = GetMiningData()
	}

	fmt.Println("Random:", currentChallenge.ChallengeString, "Difficulty:", currentChallenge.Difficulty, "Reward:", currentChallenge.Reward)
}

func CheckChallenge(magic string) (bool, [64]byte) {
	hash := sha512.Sum512([]byte(currentChallenge.ChallengeString + magic))

	stringhash := hex.EncodeToString(hash[:])
	hexarray := []rune(stringhash)

	for i := 0; i < currentChallenge.Difficulty; i++ {
		if hexarray[i] != '0' {
			return false, hash
		}
	}
	return true, hash
}

func PeriodicChallengeRefresher() {
	for {
		time.Sleep(5 * time.Second)
		RefreshCurrentChallenge()
	}
}

func PostChallengeResult(magic string, telegramid string) bool {
	response, err := http.Get(*mineapi + "/mine/resultChallenge?walletid=" + telegramid + "&magic=" + magic)
	if !*testMode && err != nil {
		panic("Error in PostChallengeResult():\n" + err.Error())
	} else {
		response.StatusCode = 202 // simulate a correct status code when in test mode
	}

	if response.StatusCode == 202 {
		fmt.Println("CHALLENGE GOT ACCEPTED, TOKENS HAVE BEEN DEPOSITED TO THE ADDRESS, getting a new challenge")
		RefreshCurrentChallenge()
	} else {
		errorMessage, _ := ioutil.ReadAll(response.Body)

		fmt.Println("Error! Returning status code: " + response.Status + " (" + string(errorMessage) + ")")
		RefreshCurrentChallenge()
	}

	return false
}
