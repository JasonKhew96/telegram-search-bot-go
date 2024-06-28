package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JasonKhew96/telegram-search-bot-go/entity"
)

func importData(filename string) {
	db, err := NewDatabase(true)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	if _, err := os.Stat(filename); err != nil {
		log.Fatalln(err)
	}
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)

	var dump entity.Dump

	count := 0

	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		if count == 2 {
			dump.Name = t.(string)
		} else if count == 4 {
			dump.Type = t.(string)
		} else if count == 6 {
			chatId, err := strconv.ParseInt(fmt.Sprintf("-100%d", int64(t.(float64))), 10, 64)
			if err != nil {
				log.Fatalln(err)
			}
			dump.Id = chatId
		} else if count >= 8 { // 8
			timeNow := time.Now().Unix()

			// _, err := db.Upser
			err = db.UpsertChat(dump.Id, dump.Name, true)
			if err != nil {
				log.Fatalln(err)
			}

			cachedPeer := make(map[int64]struct{})
			var rwMutex sync.RWMutex

			messageCount := 0
			for dec.More() {
				var msg entity.Message
				if err := dec.Decode(&msg); err != nil {
					log.Fatalln(err)
				}
				if msg.Type == "message" && msg.FullText != "" {
					fromId := int64(0)
					switch {
					case strings.HasPrefix(msg.FromId, "channel"):
						fromId, err = strconv.ParseInt(fmt.Sprintf("-100%s", strings.TrimPrefix(msg.FromId, "channel")), 10, 64)
						if err != nil {
							log.Fatalln(err)
						}
					case strings.HasPrefix(msg.FromId, "user"):
						fromId, err = strconv.ParseInt(strings.TrimPrefix(msg.FromId, "user"), 10, 64)
						if err != nil {
							log.Fatalln(err)
						}
					default:
						log.Fatalf("unknown from_id: %s", msg.FromId)
					}

					rwMutex.RLock()
					_, ok := cachedPeer[fromId]
					rwMutex.RUnlock()
					if !ok {
						rwMutex.Lock()
						cachedPeer[fromId] = struct{}{}
						rwMutex.Unlock()
						if err := db.UpsertPeer(fromId, msg.From, ""); err != nil {
							log.Fatalln(err)
						}
					}

					msgId := msg.Id
					fullText := msg.FullText
					timestamp, err := strconv.ParseInt(msg.DateUnixTime, 10, 64)
					if err != nil {
						log.Fatalln(err)
					}
					if err = db.UpsertMessage(dump.Id, fromId, msgId, fullText, timestamp); err != nil {
						log.Fatalln(err)
					}

					messageCount++
					if messageCount%10000 == 0 {
						log.Printf("imported %d messages", messageCount)
					}
				}
			}

			elapsedSeconds := time.Now().Unix() - timeNow
			log.Printf("chat %d imported in %d seconds", dump.Id, elapsedSeconds)

			break
		}
		count++
	}
}
