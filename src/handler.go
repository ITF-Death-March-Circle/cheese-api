package main

import (
	"encoding/json"
	"log"
	"time"
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
		r, err := messageHandler(requestObject.Message, requestObject.RoomId, requestObject.UserId)
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

func messageHandler(message string, room_id string, user_id string) ([]byte, error) {
	// コールバックオブジェクトを作成
	messageObject := MessageObject{
		Action:  "NOTIFY_MESSAGE",
		Time:    time.Now().String(),
		Message: "success",
	}
	b, err := json.Marshal(messageObject)
	if err != nil {
		log.Println("cannot marshal struct: %v", err)
		return nil, err
	}
	log.Println("Success Flask Server")
	return b, nil
}
