package wppscrapperimp

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
	wppscrapper "github.com/ribeiroferreiralucas/wpp-scrapper"
)

//WppScrapper //TODO:
type WppScrapper struct {
	// wppscrapper.IWppScrapper
	WhatsappConnection *whatsapp.Conn
	initialized        bool
	initializationChan chan bool
	chats              map[string]wppscrapper.Chat
	chatsToScrap       *list.List
	chatsScrapping     *list.List
	chatsScrapped      *list.List
	isScrapping        bool
	isFinished         bool
}

//InitializeConnection //TODO
func InitializeConnection() wppscrapper.IWppScrapper {
	// var wppscrapper *wppscrapper.IWppScrapper

	wppscrapper := &WppScrapper{
		isScrapping:        false,
		isFinished:         false,
		initialized:        false,
		initializationChan: make(chan bool),
		chats:              make(map[string]wppscrapper.Chat),
	}

	var err error
	wppscrapper.WhatsappConnection, err = whatsapp.NewConn(10 * time.Second)
	if err != nil {
		log.Fatalf("error creating connection: %v\n", err)
	}

	return wppscrapper
}

// Auth auths the user and returns id
func (w *WppScrapper) Auth(qrChan chan<- string) (string, error) {
	uuid := uuid.New()
	uuid, err := w.login(uuid, qrChan)
	return uuid, err
}

// ReAuth auths the user and returns id
func (w *WppScrapper) ReAuth(qrChan chan<- string, uuid string) (string, error) {

	uuid, err := w.restoreLogin(uuid)
	if err == nil {
		return uuid, err
	}
	uuid, err = w.login(uuid, qrChan)
	return uuid, err
}

func (w *WppScrapper) Initialized() bool {
	return w.initialized
}

func (w *WppScrapper) WaitInitialization() chan bool {

	return w.initializationChan
}

//GetChats recupera todos os Chats
func (w *WppScrapper) GetChats() map[string]wppscrapper.Chat {
	return w.chats
}

//StartScrapper Inicia a coleta de mensagens
func (w *WppScrapper) StartScrapper(resume bool) {

	w.chatsScrapped = list.New().Init()
	w.chatsScrapping = list.New().Init()
	w.chatsToScrap = list.New().Init()

	for _, chat := range w.chats {
		wppchat := w.WhatsappConnection.Store.Chats[chat.Jid()]
		chatHandler := createMessageHandler(w.WhatsappConnection, &wppchat)
		chatHandler.setMessagesPerCallCount(100)
		w.chatsToScrap.PushFront(chatHandler)

		chat.(*Chat).status = wppscrapper.Queue

	}
	go w.scrapRoutine(resume)
	return
}

func (w *WppScrapper) StopScrapper() {

	w.isScrapping = false
	w.isFinished = false

	for scrappingElement := w.chatsScrapping.Front(); scrappingElement != nil; scrappingElement = scrappingElement.Next() {
		chatHandler := scrappingElement.Value.(*ChatHandler)
		chatHandler.pauseChatScrapper()

		chat := w.chats[chatHandler.chat.Jid]
		chat.(*Chat).status = wppscrapper.Stoped
	}
}

func (w *WppScrapper) login(uuid string, qrChan chan<- string) (string, error) {

	qr := make(chan string)
	go func() {
		qrChan <- <-qr
	}()
	session, err := w.WhatsappConnection.Login(qr)
	if err != nil {
		return uuid, fmt.Errorf("Error during login: %v", err)
	}
	err = writeSession(session, uuid)
	if err != nil {
		return uuid, fmt.Errorf("Error saving session: %v", err)
	}

	defer func() { go w.waitInitialization() }()
	return uuid, nil
}

func (w *WppScrapper) restoreLogin(uuid string) (string, error) {
	session, err := readSession(uuid)
	if err != nil {
		return "", err
	}
	session, err = w.WhatsappConnection.RestoreWithSession(session)
	if err != nil {
		return "", err
	}
	defer func() { go w.waitInitialization() }()
	return uuid, nil
}

func (w *WppScrapper) waitInitialization() {

	var chats map[string]whatsapp.Chat
	for {
		chats = w.WhatsappConnection.Store.Chats
		if chats != nil && len(chats) > 0 {
			log.Println("Initialized")
			w.initializeChats()
			w.initialized = true
			w.initializationChan <- true
			break
		}
		fmt.Println("Chats Not Fonded, wainting 100 milliseconds to retry")
		<-time.After(100 * time.Millisecond)
	}
}

func (w *WppScrapper) initializeChats() {
	chats := w.WhatsappConnection.Store.Chats
	for key, wppchat := range chats {
		var chat Chat
		chat.id = wppchat.Jid
		chat.name = wppchat.Name
		chat.status = wppscrapper.Idle

		w.chats[key] = &chat
	}
}

func (wppscrapper *WppScrapper) scrapRoutine(resume bool) {
	simultSize := 1
	wppscrapper.isScrapping = true
	wppscrapper.isFinished = false

	for {

		if !wppscrapper.isScrapping {
			return
		}
		handleFinishedScraps(wppscrapper)

		hasNext := wppscrapper.chatsToScrap.Len() > 0
		hasScrapping := wppscrapper.chatsScrapping.Len() > 0

		if wppscrapper.chatsScrapping.Len() < simultSize && hasNext {
			startNext(wppscrapper, resume)
			continue
		}

		if !hasNext && !hasScrapping {
			wppscrapper.isFinished = true
			wppscrapper.isScrapping = false
			log.Println("WppScrapper finished to Scrap")
			return
		}
		<-time.After(100 * time.Millisecond)
	}
}

func handleFinishedScraps(scrapper *WppScrapper) {
	for scrappingElement := scrapper.chatsScrapping.Front(); scrappingElement != nil; scrappingElement = scrappingElement.Next() {

		scrappingHandler := scrappingElement.Value.(*ChatHandler)
		if scrappingHandler.isScrapped {
			log.Println("Removing " + scrappingHandler.chat.Jid + " from scrapping list and adding on scrapped list")
			scrapper.chatsScrapping.Remove(scrappingElement)
			scrapper.chatsScrapped.PushBack(scrappingHandler)

			chat := scrapper.chats[scrappingHandler.chat.Jid]
			chat.(*Chat).status = wppscrapper.Finished
		}
	}
}

func startNext(scrapper *WppScrapper, resume bool) {
	handlerToStart := scrapper.chatsToScrap.Front()
	scrapper.chatsScrapping.PushBack(handlerToStart.Value)
	scrapper.chatsToScrap.Remove(handlerToStart)

	chatHandler := handlerToStart.Value.(*ChatHandler)
	chat := scrapper.chats[chatHandler.chat.Jid]
	chat.(*Chat).status = wppscrapper.Running

	go chatHandler.startChatScrapper(resume)
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
