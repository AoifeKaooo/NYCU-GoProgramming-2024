package main

import (
	"bufio"
	"context"
	"fmt"
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
			// Broadcast incoming message to all
			// clients' outgoing message channels.
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
		conn.WriteMessage(1, []byte(msg))
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
	messages <- rxgo.Of(who + " 來到了現場" + "\n")
	entering <- ch

	defer func() {
		log.Println("disconnect !!")
		leaving <- ch
		messages <- rxgo.Of(who + " 離開了" + "\n")
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
	// 读取 swear_word.txt 和 sensitive_name.txt 的内容
	swearWords := readLines("swear_word.txt")
	sensitiveNames := readLines("sensitive_name.txt")

	// Filter 和 Map 的实现
	ObservableMsg = ObservableMsg.
		Filter(func(item interface{}) bool {
			// 类型断言为 rxgo.Item
			msg, ok := item.(rxgo.Item)
			if !ok {
				return false
			}
			text, ok := msg.V.(string)
			if !ok {
				return false
			}

			// 如果消息包含任意一个 swear word，返回 false 以过滤掉
			for _, word := range swearWords {
				if containsWord(text, word) {
					return false
				}
			}
			return true
		}).
		Map(func(_ context.Context, item interface{}) (interface{}, error) {
			// 类型断言为 rxgo.Item
			msg, ok := item.(rxgo.Item)
			if !ok {
				return nil, fmt.Errorf("invalid item type")
			}
			text, ok := msg.V.(string)
			if !ok {
				return nil, fmt.Errorf("invalid message content")
			}

			// 如果消息中包含 sensitive name，替换第二个字为 '*'
			for _, name := range sensitiveNames {
				if containsWord(text, name) {
					maskedName := maskName(name)
					text = replaceWord(text, name, maskedName)
				}
			}
			return rxgo.Of(text), nil
		})
}

// 读取文件内容到字符串切片
func readLines(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("无法打开文件 %s: %v", filePath, err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("读取文件 %s 出错: %v", filePath, err)
	}
	return lines
}

// 检查消息是否包含特定的单词
func containsWord(msg, word string) bool {
	return strings.Contains(msg, word)
}

// 替换敏感名字的第二个字为 '*'
func maskName(name string) string {
	if len(name) >= 3 {
		runes := []rune(name)
		runes[1] = '*'
		return string(runes)
	}
	return name
}

// 替换消息中的单词
func replaceWord(msg, oldWord, newWord string) string {
	return strings.ReplaceAll(msg, oldWord, newWord)
}

func main() {
	InitObservable()
	go broadcaster()
	http.HandleFunc("/wschatroom", wshandle)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("server start at :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
