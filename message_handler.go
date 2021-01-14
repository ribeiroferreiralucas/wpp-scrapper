package wppscrapper

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	whatsapp "github.com/Rhymen/go-whatsapp"
)

var headers = []string{"message_id", "timestamp", "chat_name", "chat", "sender", "is_forwarded", "quoted_message_id", "message"}

//MessageHandler TODO
type MessageHandler struct {
	c            *whatsapp.Conn
	message      chan Message
	writer       *csv.Writer
	writerChatID string
}

//Message TODO
type Message struct {
	MessageID       string `csv:"message_id"`
	Timestamp       uint64 `csv:"timestamp"`
	ChatName        string `csv:"chat_name"`
	ChatID          string `csv:"chat"`
	Sender          string `csv:"sender"`
	IsForwarded     bool   `csv:"is_forwarded"`
	QuotedMessageID string `csv:"quoted_message_id"`
	Text            string `csv:"message_content"`
}

//HandleError needs to be implemented to be a valid WhatsApp handler
func (h *MessageHandler) HandleError(err error) {

	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		err := h.c.Restore()
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
	} else {
		log.Printf("error occoured: %v\n", err)
	}
}

//HandleTextMessage Optional to be implemented. Implement HandleXXXMessage for the types you need.
func (h *MessageHandler) HandleTextMessage(message whatsapp.TextMessage) {
	newMessage := toMessage(message, *h)

	isNewWriter := h.updateWriterIfNeeded(newMessage.ChatID)
	if isNewWriter {
		h.writer.Write(headers)
	}

	data := toCsv(newMessage)
	h.writer.Write(data)
	log.Println(data)

}

func (h *MessageHandler) updateWriterIfNeeded(chatID string) bool {

	if len(h.writerChatID) == 0 || h.writerChatID != chatID {

		if h.writer != nil {
			h.writer.Flush()
		}
		file, err := os.Create(chatID + "-messages.csv")
		checkError("Cannot create file", err)

		h.writer = csv.NewWriter(file)

		h.writerChatID = chatID
		return true
	}

	return false
}

func toCsv(message Message) []string {
	return []string{message.MessageID, strconv.FormatUint(message.Timestamp, 10), message.ChatName, message.ChatID, message.Sender, strconv.FormatBool(message.IsForwarded), message.QuotedMessageID, message.Text}
}
func toMessage(wppMessage whatsapp.TextMessage, handler MessageHandler) Message {

	var jid string
	if wppMessage.Info.Source.Participant == nil {
		jid = wppMessage.Info.RemoteJid
	} else {
		jid = *wppMessage.Info.Source.Participant
	}

	message := Message{
		MessageID:       wppMessage.Info.Id,
		ChatID:          wppMessage.Info.RemoteJid,
		ChatName:        handler.c.Store.Chats[wppMessage.Info.RemoteJid].Name,
		Sender:          jid,
		Timestamp:       wppMessage.Info.Timestamp,
		IsForwarded:     wppMessage.ContextInfo.IsForwarded,
		QuotedMessageID: wppMessage.Info.Id,
		Text:            wppMessage.Text,
	}

	return message
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
