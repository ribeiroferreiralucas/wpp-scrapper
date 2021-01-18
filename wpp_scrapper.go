package wppscrapper

import (
	"container/list"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	whatsapp "github.com/Rhymen/go-whatsapp"
	uuid "github.com/pborman/uuid"
)

//WppScrapper //TODO:
type WppScrapper struct {
	WhatsappConnection *whatsapp.Conn
	chatsToScrap       *list.List
	chatsScrapping     *list.List
	chatsScrapped      *list.List
	isScrapping        bool
}

//InitializeConnection //TODO
func InitializeConnection() *WppScrapper {
	wppscrapper := &WppScrapper{}
	var err error
	wppscrapper.WhatsappConnection, err = whatsapp.NewConn(10 * time.Second)
	if err != nil {
		log.Fatalf("error creating connection: %v\n", err)
	}
	return wppscrapper
}

// Auth auths the user and returns id
func (wppscrapper *WppScrapper) Auth(qrChan chan<- string) (string, error) {

	qr := make(chan string)
	go func() {
		qrChan <- <-qr
	}()

	session, err := wppscrapper.WhatsappConnection.Login(qr)
	if err != nil {
		return "", fmt.Errorf("error during login: %v\n", err)
	}

	uuid := uuid.New()
	err = writeSession(session, uuid)
	if err != nil {
		return "", fmt.Errorf("error saving session: %v\n", err)
	}
	return uuid, nil
}

// ReAuth auths the user and returns id
func (wppscrapper *WppScrapper) ReAuth(qrChan chan<- string, uuid string) (string, error) {

	session, err := readSession(uuid)
	if err == nil {
		session, err = wppscrapper.WhatsappConnection.RestoreWithSession(session)
		if err == nil {
			return uuid, nil
		}
	}

	qr := make(chan string)
	go func() {
		qrChan <- <-qr
	}()
	session, err = wppscrapper.WhatsappConnection.Login(qr)
	if err != nil {
		return "", fmt.Errorf("Error during login: %v", err)
	}
	err = writeSession(session, uuid)
	if err != nil {
		return "", fmt.Errorf("Error saving session: %v", err)
	}

	return uuid, nil
}

//TODO: Implementar uma forma boa de saber se a inicialização ja terminou (pegou as informações dos chats)

//GetAllChats recupera todos os Chats
func (wppscrapper *WppScrapper) GetAllChats() map[string]whatsapp.Chat {
	return wppscrapper.WhatsappConnection.Store.Chats
}

//StartScrapper Inicia a coleta de mensagens
func (wppscrapper *WppScrapper) StartScrapper(resume bool) {

	wppscrapper.chatsScrapped = list.New().Init()
	wppscrapper.chatsScrapping = list.New().Init()
	wppscrapper.chatsToScrap = list.New().Init()

	for _, chat := range wppscrapper.WhatsappConnection.Store.Chats {
		chatHandler := CreateMessageHandler(wppscrapper.WhatsappConnection, chat)
		wppscrapper.chatsToScrap.PushFront(chatHandler)
	}
	go wppscrapper.startScrapper(resume)
	return
}

func (wppscrapper *WppScrapper) StopScrapper() {

	wppscrapper.isScrapping = false

	for scrappingElement := wppscrapper.chatsScrapping.Front(); scrappingElement != nil; scrappingElement = scrappingElement.Next() {
		chatHandler := scrappingElement.Value.(ChatHandler)
		chatHandler.PauseChatScrapper()
	}
}

func (wppscrapper *WppScrapper) startScrapper(resume bool) {
	simultSize := 1
	//TODO: Tratar o final e errosx
	wppscrapper.isScrapping = true

	for {

		if !wppscrapper.isScrapping {
			return
		}

		for scrappingElement := wppscrapper.chatsScrapping.Front(); scrappingElement != nil; scrappingElement = scrappingElement.Next() {

			scrapping := scrappingElement.Value.(*ChatHandler)

			if scrapping.IsScrapped() {
				wppscrapper.chatsScrapping.Remove(scrappingElement)
				wppscrapper.chatsScrapped.PushBack(scrapping)
			}
		}

		if wppscrapper.chatsScrapping.Len() < simultSize {
			handlerToStart := wppscrapper.chatsToScrap.Front()
			wppscrapper.chatsScrapping.PushBack(handlerToStart.Value)
			wppscrapper.chatsToScrap.Remove(handlerToStart)

			chat := handlerToStart.Value.(*ChatHandler)
			go chat.StartChatScrapper(resume)
			continue
		}

		<-time.After(100 * time.Millisecond)
	}
}

func readSession(uuid string) (whatsapp.Session, error) {
	session := whatsapp.Session{}

	dir, err := os.UserConfigDir()
	if err != nil {
		return session, err
	}

	filepath := path.Join(dir, "wpp_scrapper", uuid+".gob")
	fmt.Println(filepath)
	file, err := os.Open(filepath)
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func writeSession(session whatsapp.Session, uuid string) error {

	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	filepath := path.Join(dir, "wpp_scrapper", uuid+".gob")
	os.MkdirAll(path.Dir(filepath), os.ModePerm)
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}
