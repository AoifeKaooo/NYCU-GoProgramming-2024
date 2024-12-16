package main

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/reactivex/rxgo/v2"
)

type client chan<- string // an outgoing message channel

var (
	entering      = make(chan client)
	leaving       = make(chan client)
	messages      = make(chan rxgo.Item) // all incoming client messages
	ObservableMsg = rxgo.FromChannel(messages)
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	MessageBroadcast := ObservableMsg.Observe()
	for {
		select {
		case msg := <-MessageBroadcast:
			// Broadcast incoming message to all clients
			for cli := range clients {
				cli <- msg.V.(string)
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func clientWriter(conn *websocket.Conn, ch <-chan string) {
	for msg := range ch {
		conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}

func wshandle(w http.ResponseWriter, r *http.Request) {
	upgrader := &websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "你是 " + who + "\n"
	messages <- rxgo.Of(who + " 來到了現場\n")
	entering <- ch

	defer func() {
		log.Println("disconnect !!")
		leaving <- ch
		messages <- rxgo.Of(who + " 離開了\n")
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		messages <- rxgo.Of(who + " 表示: " + string(msg))
	}
}

func InitObservable() {
	// 載入髒話與敏感字詞清單
	swearWords := loadWords("swear_word.txt")
	sensitiveWords := loadWords("sensitive_name.txt")


	ObservableMsg = ObservableMsg.
		Filter(func(i interface{}) bool {
			msg := i.(string)
			return !containsSwearWords(msg, swearWords)
		}).
		Map(func(_ context.Context, i interface{}) (interface{}, error) {
			msg := i.(string)
			return replaceSensitiveWords(msg, sensitiveWords), nil
		})
}

func loadWords(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening %s: %v", filename, err)
		return nil
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			words = append(words, word)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading %s: %v", filename, err)
	}
	return words
}

func containsSwearWords(msg string, swearWords []string) bool {
	for _, w := range swearWords {
		if w != "" && strings.Contains(msg, w) {
			return true
		}
	}
	return false
}

func replaceSensitiveWords(msg string, sensitiveWords []string) string {
	for _, w := range sensitiveWords {
		if strings.Contains(msg, w) {
			runes := []rune(w)
			if len(runes) > 1 {
				masked := string(runes[0]) + "*"
				if len(runes) > 2 {
					masked += string(runes[2:])
				}
				msg = strings.ReplaceAll(msg, w, masked)
			}
		}
	}
	return msg
}

func main() {
	InitObservable()
	go broadcaster()
	http.HandleFunc("/wschatroom", wshandle)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("server start at :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}