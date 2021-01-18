package main

import (
	"fmt"
	"log"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	whatsapp "github.com/Rhymen/go-whatsapp"
	wppscrapper "github.com/ribeiroferreiralucas/wpp-scrapper"
)

type MessageHandler struct {
	c *whatsapp.Conn
}

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

	var jid string
	if message.Info.Source.Participant == nil {
		jid = message.Info.RemoteJid
	} else {
		jid = *message.Info.Source.Participant
	}

	fmt.Printf("%v %v %v %v %v\n\t%v\n", h.c.Store.Chats[message.Info.RemoteJid].Name, jid, message.Info.Id, message.Info.Timestamp, message.ContextInfo.QuotedMessageID, message.Text)
}

func main() {

	scrapper := wppscrapper.InitializeConnection()

	qr := make(chan string)
	go func() {
		terminal := qrcodeTerminal.New()
		terminal.Get(<-qr).Print()
	}()

	_, err := scrapper.ReAuth(qr, "other")
	if err != nil {
		log.Fatalf("error scrapper.ReAuth in: %v\n", err)
	}

	var chats map[string]whatsapp.Chat
	for true {
		chats = scrapper.WhatsappConnection.Store.Chats

		if chats != nil && len(chats) > 0 {
			fmt.Println("Chats Fonded")
			break
		}
		fmt.Println("Chats Not Fonded, wainting 100 milliseconds to retry")
		<-time.After(100 * time.Millisecond)
	}

	for k, v := range chats {
		fmt.Println("k:", k, "v:", v)
	}

	fmt.Println("---------------\n\n\n\nSTART SCRAPPER\n\n\n\n----------------")
	scrapper.StartScrapper(true)
	<-time.After(1000000000 * time.Second)
	fmt.Println("---------------\n\n\n\nSTOP SCRAPPER\n\n\n\n----------------")
	scrapper.StopScrapper()

}
