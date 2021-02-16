package main

import (
	"fmt"
	"log"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/ribeiroferreiralucas/wpp-scrapper/wppscrapperimp"
)

func main() {

	scrapper := wppscrapperimp.InitializeConnection()

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
	<-time.After(1000000000 * time.Second)
	fmt.Println("---------------\n\n\n\nSTOP SCRAPPER\n\n\n\n----------------")
	scrapper.StopScrapper()

}
