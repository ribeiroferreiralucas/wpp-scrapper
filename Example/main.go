package main

import (
	"fmt"
	"log"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	wppscrapper "github.com/ribeiroferreiralucas/wpp-scrapper"
	"github.com/ribeiroferreiralucas/wpp-scrapper/wppscrapperimp"
)

type EventLogger struct {
}

func (logger *EventLogger) OnWppScrapperStarted(wppScrapper wppscrapper.IWppScrapper) {
	fmt.Println("logger OnWppScrapperStarted")
}

func (logger *EventLogger) OnWppScrapperStopped(wppScrapper wppscrapper.IWppScrapper) {
	fmt.Println("logger OnWppScrapperStopped")
}

func (logger *EventLogger) OnWppScrapperFinished(wppScrapper wppscrapper.IWppScrapper) {
	fmt.Println("logger OnWppScrapperFinished")
}

func (logger *EventLogger) OnWppScrapperChatScrapStarted(chat wppscrapper.Chat) {
	fmt.Println("logger OnWppScrapperChatScrapStarted")
}

func (logger *EventLogger) OnWppScrapperChatScrapFinished(chat wppscrapper.Chat) {
	fmt.Println("logger OnWppScrapperChatScrapFinished")
}

func main() {

	scrapper := wppscrapperimp.InitializeConnection().(wppscrapper.IWppScrapper)

	evtLogger := &EventLogger{}

	evtHandler := scrapper.GetWppScrapperEventHandler()

	fmt.Println(evtHandler)

	evtHandler.AddOnScrapperStartedListener(evtLogger)
	evtHandler.AddOnScrapperStoppedListener(evtLogger)
	evtHandler.AddOnScrapperFinishedListener(evtLogger)
	evtHandler.AddOnChatScrapStartedListener(evtLogger)
	evtHandler.AddOnChatScrapFinishedListener(evtLogger)

	qr := make(chan string)
	go func() {
		terminal := qrcodeTerminal.New()
		terminal.Get(<-qr).Print()
	}()

	_, err := scrapper.ReAuth(qr, "other")
	if err != nil {
		log.Fatalf("error scrapper.ReAuth in: %v\n", err)
	}

	if !scrapper.Initialized() {
		<-scrapper.WaitInitialization()
	}

	for k, v := range scrapper.GetChats() {
		fmt.Println("k:", k, "v:", v)
	}

	fmt.Println("---------------\n\n\n\nSTART SCRAPPER\n\n\n\n----------------")
	scrapper.StartScrapper(true)
	<-time.After(time.Second)
	fmt.Println("---------------\n\n\n\nSTOP SCRAPPER\n\n\n\n----------------")
	scrapper.StopScrapper()

}
