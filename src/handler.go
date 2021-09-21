package main

import (
	"encoding/json"
	"log"
	"main/redis"

	"github.com/pkg/errors"
)

const VOTE_BACKGROUND = "VOTE_BACKGROUND"

func handler(s []byte) []byte {
	var requestObject SocketRequest
	if err := json.Unmarshal(s, &requestObject); err != nil {
		return errorSocketResponse
	}

	//各アクションケースに応じて処理を行う
	switch {
	case requestObject.Action == VOTE_BACKGROUND:
		r, err := messageHandler(requestObject.Key)
		if err != nil {
			return errorSocketResponse
		}
		return r
	}
	return errorSocketResponse
}

// func loadMaidMessage() (string, error) {
// 	fileName, err := fileName()
// 	if err != nil {
// 		return "Error!", err
// 	}
// 	f, err := os.Open(fileName)
// 	if err != nil {
// 		log.Println(err)
// 		return "Error!", err
// 	}

// 	defer f.Close()
// 	var strSlice []string
// 	scanner := bufio.NewScanner(f)

// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		//fmt.Println(line)
// 		strSlice = append(strSlice, line)
// 	}
// 	if err := scanner.Err(); err != nil {
// 		return "Error!", err
// 	}
// 	rand.Seed(time.Now().UnixNano())
// 	num := rand.Intn(len(strSlice))
// 	return strSlice[num], nil
// }

// 各アクションタイプ毎のハンドラー

func getVoteValue() (value int, err error) {
	value_1, err := count(VOTE_PATTERNS[0])
	if err != nil {
		return -1, err
	}

	tmp := value_1
	value = 0

	value_2, err := count(VOTE_PATTERNS[1])

	if err != nil {
		return -1, err
	}
	if tmp < value_2 {
		tmp = value_2
		value = 1
	}

	value_3, err := count(VOTE_PATTERNS[2])

	if err != nil {
		return -1, err
	}
	if tmp < value_3 {
		tmp = value_3
		value = 2
	}
	return
}

func count(target string) (value int, err error) {
	err = redis.AddValue(target)
	if err != nil {
		return -1, errors.Wrap(err, "Failed to add connection")
	}

	value, err = redis.DeclValue(target)
	if err != nil {
		return -1, errors.Wrap(err, "Failed to decl connection")
	}
	return
}

func getCurrentValue() ([]byte, error) {
	value, err := getVoteValue()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get VotedValue")
	}
	connection_value, err := count(COUNT_USER)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get connections")
	}

	var patternText string
	if value < 0 || len(PATTERNS) <= value {
		patternText = "分かりません"
	} else {
		patternText = PATTERNS[value]
	}

	messageObject := VoteResponse{
		Action: "RESULT_VOTE",
		Text:   patternText,
		Count:  connection_value,
		Value:  value,
	}
	b, err := json.Marshal(messageObject)
	if err != nil {
		log.Println("cannot marshal struct: %v", err)
		return nil, err
	}
	return b, nil

}

func messageHandler(message string) ([]byte, error) {
	// コールバックオブジェクトを作成
	switch {
	case message == VOTE_PATTERNS[0]:
		err := redis.AddValue(VOTE_PATTERNS[0])
		if err != nil {
			log.Println("failed to call vote: %v", err)
			return nil, err
		}
	case message == VOTE_PATTERNS[1]:
		err := redis.AddValue(VOTE_PATTERNS[1])
		if err != nil {
			log.Println("failed to call vote: %v", err)
			return nil, err
		}
	case message == VOTE_PATTERNS[2]:
		err := redis.AddValue(VOTE_PATTERNS[2])
		if err != nil {
			log.Println("failed to call vote: %v", err)
			return nil, err
		}
	default:
		err := redis.AddValue(VOTE_PATTERNS[0])
		if err != nil {
			log.Println("failed to call vote: %v", err)
			return nil, err
		}
	}

	return getCurrentValue()
}
