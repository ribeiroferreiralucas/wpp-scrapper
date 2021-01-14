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

//WppScrapper TODO:
type WppScrapper struct {
	WhatsappConnection *whatsapp.Conn
	ChatsToScrap       *list.List
	isScrapping        bool
}

//InitializeConnection TODO
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

//GetAllChats recupera todos os Chats
func (wppscrapper *WppScrapper) GetAllChats() map[string]whatsapp.Chat {
	return wppscrapper.WhatsappConnection.Store.Chats
}

//StartScrapper Inicia a coleta de mensagens
func (wppscrapper *WppScrapper) StartScrapper() {
	messageHandler := MessageHandler{
		c: wppscrapper.WhatsappConnection,
	}

	var duration_Milliseconds time.Duration = 1 * time.Millisecond

	for k := range wppscrapper.WhatsappConnection.Store.Chats {

		wppscrapper.WhatsappConnection.LoadFullChatHistory(k, 1, duration_Milliseconds, &messageHandler)
	}

	wppscrapper.isScrapping = true
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
