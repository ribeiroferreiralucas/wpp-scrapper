package wppscrapper

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	whatsapp "github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"
)

var headers = []string{"message_id", "timestamp", "chat_name", "chat", "sender", "is_forwarded", "from_me", "quoted_message_id", "message"}

//MessageHandler //TODO
type MessageHandler struct {
	Chat               *whatsapp.Chat
	IsScrapping        bool //TODO: Transformar em readonly
	WasScrapped        bool
	lastMessage        string
	lastMessageOwner   bool
	c                  *whatsapp.Conn
	shouldStopScrapper bool
	message            chan Message
	writer             *csv.Writer
	file               *os.File
}

//Message //TODO
type Message struct {
	MessageID       string
	Timestamp       uint64
	ChatName        string
	ChatID          string
	Sender          string
	IsForwarded     bool
	FromMe          bool
	QuotedMessageID string
	Text            string
}

// CreateMessageHandler //TODO
func CreateMessageHandler(conn *whatsapp.Conn, chat *whatsapp.Chat) *MessageHandler {

	return &MessageHandler{
		c:                  conn,
		Chat:               chat,
		lastMessage:        "",
		lastMessageOwner:   true,
		shouldStopScrapper: false,
	}
}

// StartChatScrapper //TODO
func (h *MessageHandler) StartChatScrapper(resume bool) {

	lastMessage := "trash"
	h.IsScrapping = true

	if !h.hasTempFile() {
		h.createWriter()
		h.writer.Write(headers)
	} else {
		h.reopenFile()
	}

	for {

		if h.shouldStopScrapper {
			return
		}

		if lastMessage == h.lastMessage {
			h.finishedScrapper()
			return
		}
		lastMessage = h.lastMessage
		h.c.LoadChatMessages(h.Chat.Jid, 1, h.lastMessage, h.lastMessageOwner, false, h)
	}
}

// StopChatScrapper //TODO
func (h *MessageHandler) StopChatScrapper() {

	//TODO: Implementar função adequadamente
	h.shouldStopScrapper = true
	log.Println("STOPED")
}

//ShouldCallSynchronously //TODO
func (h *MessageHandler) ShouldCallSynchronously() bool {
	return true
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
	data := toCsv(newMessage)
	h.writer.Write(data)
	h.writer.Flush()
	log.Println(data)
}

//HandleRawMessage //TODO
func (h *MessageHandler) HandleRawMessage(message *proto.WebMessageInfo) {
	h.lastMessage = *message.Key.Id
	h.lastMessageOwner = message.Key.FromMe != nil && *message.Key.FromMe
	log.Println("lastMessage> " + h.lastMessage)
}

func (h *MessageHandler) finishedScrapper() {

	h.writer.Flush()
	h.file.Close()

	err := os.Rename(h.Chat.Jid+"-messages.temp", h.Chat.Jid+"-messages.csv")
	checkError("Cannot rename file from .temp to .csv", err)

	h.WasScrapped = true
	h.IsScrapping = false
}

func (h *MessageHandler) hasTempFile() bool {

	_, err := os.Stat(h.Chat.Jid + "-messages.temp")

	return err == nil
}

func (h *MessageHandler) reopenFile() {
	file, err := os.OpenFile(h.Chat.Jid+"-messages.temp", os.O_APPEND, os.ModeAppend)

	checkError("Cannot open file", err)
	reader := csv.NewReader(file)

	data, _ := reader.ReadAll()

	rowsCount := len(data)

	h.lastMessage = data[rowsCount-1][0]
	h.lastMessageOwner, _ = strconv.ParseBool(data[rowsCount-1][6])

	h.writer = csv.NewWriter(file)

}

func (h *MessageHandler) createWriter() {

	file, err := os.Create(h.Chat.Jid + "-messages.temp")
	checkError("Cannot create file", err)

	h.writer = csv.NewWriter(file)
	h.file = file
}

func toCsv(message Message) []string {
	return []string{message.MessageID, strconv.FormatUint(message.Timestamp, 10), message.ChatName, message.ChatID, message.Sender, strconv.FormatBool(message.IsForwarded), strconv.FormatBool(message.FromMe), message.QuotedMessageID, message.Text}
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
		FromMe:          wppMessage.Info.FromMe,
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
