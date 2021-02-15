package wppscrapper

import (
	"container/list"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	whatsapp "github.com/Rhymen/go-whatsapp"

	"github.com/Rhymen/go-whatsapp/binary/proto"
)

var headers = []string{"message_id", "timestamp", "chat_name", "chat", "sender", "is_forwarded", "from_me", "quoted_message_id", "message"}

//ChatHandler //TODO
type ChatHandler struct {
	chat                 *whatsapp.Chat
	chatInfo             *ChatInfo
	messagesPerCallCount int
	isScrapping          bool
	isScrapped           bool
	conn                 *whatsapp.Conn
	lastMessage          string
	lastMessageOwner     bool
	shouldStopScrapper   bool
	collectedMessages    *list.List
	writer               *csv.Writer
	file                 *os.File
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

/*
 {
	 "id":"GROUPID@g.us",
	 "owner":"OWNERID@c.us",
	 "subject":"GROUP NAME",
	 "creation":1592654382,
	 "participants":[
		 {"id":"XXX@c.us","isAdmin":false,"isSuperAdmin":false},
		 {"id":"ID@c.us","isAdmin":false,"isSuperAdmin":false}
		 ],
	"subjectTime":1592654382,
	"subjectOwner":"ID@c.us",
	"desc":" DESC",
	"descId":"738D7C3293DA5137DAA2C8C50F944369",
	"descTime":1594328067,
	"descOwner":"5521999034100@c.us"}

*/
type ChatInfo struct {
	Jid               string       `json:"id"`
	OwnerJid          string       `json:"owner"`
	Name              string       `json:"subject"`
	CreationTimestamp uint         `json:"creation"`
	Participants      []ChatMember `json:"participants"`
	Desc              string       `json:"desc"`
}

type ChatMember struct {
	Jid          string `json:"id"`
	isAdmin      bool
	isSuperAdmin bool
}

func createMessageHandler(conn *whatsapp.Conn, chat *whatsapp.Chat) *ChatHandler {

	h := ChatHandler{
		messagesPerCallCount: 1,
		conn:                 conn,
		chat:                 chat,
		isScrapping:          false,
		isScrapped:           false,
		lastMessage:          "",
		lastMessageOwner:     true,
		shouldStopScrapper:   false,
		collectedMessages:    list.New().Init(),
	}

	h.isScrapped = h.hasFinalFile()
	return &h

}

//ShouldCallSynchronously is a MessageHandler interface method thats return if the handler should work async or sync
func (h *ChatHandler) ShouldCallSynchronously() bool {
	return true
}

//HandleError is a method to implements the interface of a minimal handler that is used to dispatch error messages. These errors occur on unexpected behavior by the websocket
//connection or if we are unable to handle or interpret an incoming message. Error produced by user actions are not
//dispatched through this handler. They are returned as an error on the specific function call.
func (h *ChatHandler) HandleError(err error) {

	//TODO: Implementar tratamento e resposta a erros

	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		err := h.conn.Restore()
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
	} else {
		log.Printf("error occoured: %v\n", err)
	}
}

//HandleTextMessage TextMessageHandler interface needs to be implemented to receive text messages dispatched by the dispatcher.
func (h *ChatHandler) HandleTextMessage(message whatsapp.TextMessage) {

	newMessage := toMessage(message, *h)
	h.collectedMessages.PushFront(newMessage)
	log.Println(h.collectedMessages.Len())
}

//HandleRawMessage The RawMessageHandler interface needs to be implemented to receive raw messages dispatched by the dispatcher.
// Raw messages are the raw protobuf structs instead of the easy-to-use structs in TextMessageHandler, ImageMessageHandler, etc..
func (h *ChatHandler) HandleRawMessage(message *proto.WebMessageInfo) {

	if h.collectedMessages.Len() > 1 {
		return
	}

	h.lastMessage = *message.Key.Id
	h.lastMessageOwner = message.Key.FromMe != nil && *message.Key.FromMe
}

func (h *ChatHandler) startChatScrapper(resume bool) {

	log.Println("Starting Scrap for " + h.chat.Jid)
	h.fetchGroupInfo()
	if resume && h.isScrapped {
		log.Println("Stopped Scrap for " + h.chat.Jid + ". Chat already scrapped and resume mode is on")
		return
	}

	if !resume && h.hasFinalFile() {
		h.deleteFinalFile()
		log.Println("Deleted Scrapped file for " + h.chat.Jid + ". Resume mode is off")
	}

	if !resume || !h.hasTempFile() {
		h.createWriter()
		h.writer.Write(headers)
		log.Println("starting Scrap " + h.chat.Jid + ". Resume mode is off")
	} else {
		h.reopenFile()
		log.Println("resuming Scrap " + h.chat.Jid + ". Resume mode is on")
	}

	h.isScrapping = true

	lastMessage := ""
	firstIteration := true
	for {

		if h.shouldStopScrapper {
			h.stoppedScrapper()
			log.Println("paused Scrap " + h.chat.Jid)
			return
		}
		if !firstIteration && lastMessage == h.lastMessage {
			h.finishedScrapper()
			log.Println("finished to Scrap " + h.chat.Jid)
			return
		}
		lastMessage = h.lastMessage

		h.conn.LoadChatMessages(h.chat.Jid, h.messagesPerCallCount, h.lastMessage, h.lastMessageOwner, false, h)
		h.writeMessages()

		firstIteration = false
	}
}
func (h *ChatHandler) fetchGroupInfo() {

	info, _ := h.conn.GetGroupMetaData(h.chat.Jid)
	infoJson := <-info

	jsonBytes := []byte(infoJson)
	var chatInfo ChatInfo
	err := json.Unmarshal(jsonBytes, &chatInfo)
	checkError("Failed to Unmarshal", err)
	log.Println(chatInfo)
}

// pauseChatScrapper //TODO
func (h *ChatHandler) pauseChatScrapper() {
	h.shouldStopScrapper = true
}

func (h *ChatHandler) setMessagesPerCallCount(messagesPerCallCount int) {

	if messagesPerCallCount < 1 {
		messagesPerCallCount = 1
	}
	h.messagesPerCallCount = messagesPerCallCount
}

func (h *ChatHandler) writeMessages() {

	for messageNode := h.collectedMessages.Front(); messageNode != nil; messageNode = messageNode.Next() {
		message := messageNode.Value.(Message)
		data := toCsv(message)
		h.writer.Write(data)
		log.Println(data)
	}

	h.writer.Flush()
	h.collectedMessages.Init()
}

func (h *ChatHandler) finishedScrapper() {

	h.writer.Flush()
	h.file.Close()

	err := os.Rename(h.chat.Jid+"-messages.temp", h.chat.Jid+"-messages.csv")
	checkError("Cannot rename file from .temp to .csv", err)

	// h.ChatInfo
	//TODO: Implementar lógica de recuperar informações do grupo
	h.isScrapped = true
	h.isScrapping = false
}

func (h *ChatHandler) stoppedScrapper() {

	h.writer.Flush()
	h.file.Close()

	h.isScrapped = false
	h.isScrapping = false
}

func (h *ChatHandler) hasTempFile() bool {

	_, err := os.Stat(h.chat.Jid + "-messages.temp")

	return err == nil
}

func (h *ChatHandler) hasFinalFile() bool {

	_, err := os.Stat(h.chat.Jid + "-messages.csv")

	return err == nil
}

func (h *ChatHandler) deleteFinalFile() {
	os.Remove(h.chat.Jid + "-messages.csv")
}

func (h *ChatHandler) reopenFile() {
	file, err := os.OpenFile(h.chat.Jid+"-messages.temp", os.O_APPEND, os.ModeAppend)

	checkError("Cannot open file", err)
	reader := csv.NewReader(file)

	data, _ := reader.ReadAll()

	rowsCount := len(data)

	h.lastMessage = data[rowsCount-1][0]
	h.lastMessageOwner, _ = strconv.ParseBool(data[rowsCount-1][6])

	h.writer = csv.NewWriter(file)
}

func (h *ChatHandler) createWriter() {

	file, err := os.Create(h.chat.Jid + "-messages.temp")
	checkError("Cannot create file", err)

	h.writer = csv.NewWriter(file)
	h.file = file
}

func toCsv(message Message) []string {
	return []string{message.MessageID, strconv.FormatUint(message.Timestamp, 10), message.ChatName, message.ChatID, message.Sender, strconv.FormatBool(message.IsForwarded), strconv.FormatBool(message.FromMe), message.QuotedMessageID, message.Text}
}

func toMessage(wppMessage whatsapp.TextMessage, handler ChatHandler) Message {

	var jid string
	if wppMessage.Info.Source.Participant == nil {
		jid = wppMessage.Info.RemoteJid
	} else {
		jid = *wppMessage.Info.Source.Participant
	}

	message := Message{
		MessageID:       wppMessage.Info.Id,
		ChatID:          wppMessage.Info.RemoteJid,
		ChatName:        handler.conn.Store.Chats[wppMessage.Info.RemoteJid].Name,
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
