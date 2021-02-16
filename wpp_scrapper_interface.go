package wppscrapper

type IWppScrapper interface {
	Auth(qrChan chan<- string) (string, error)
	ReAuth(qrChan chan<- string, uuid string) (string, error)
	WaitInitialization() chan bool
	Initialized() bool
	GetChats() map[string]Chat
	StartScrapper(resume bool)
	StopScrapper()
	GetWppScrapperEventHandler() IWppScrapperEventHandler
}

type ChatStatus int

const (
	Idle     = -1
	Queue    = 0
	Running  = 1
	Stoped   = 2
	Finished = 3
)

type Chat interface {
	Name() string
	Jid() string
	GetStatus() ChatStatus
}

type IWppScrapperStartedListener interface {
	OnWppScrapperStarted(wppScrapper IWppScrapper)
}
type IWppScrapperStoppedListener interface {
	OnWppScrapperStopped(wppScrapper IWppScrapper)
}
type IWppScrapperFinishedListener interface {
	OnWppScrapperFinished(wppScrapper IWppScrapper)
}
type IWppScrapperChatScrapStartedListener interface {
	OnWppScrapperChatScrapStarted(chat Chat)
}
type IWppScrapperChatScrapFinishedListener interface {
	OnWppScrapperChatScrapFinished(chat Chat)
}

type IWppScrapperEventHandler interface {
	AddOnScrapperStartedListenner(IWppScrapperStartedListener)
	AddOnScrapperStoppedListenner(IWppScrapperStoppedListener)
	AddOnScrapperFinishedListenner(IWppScrapperFinishedListener)
	AddOnChatScrapStartedListenner(IWppScrapperChatScrapStartedListener)
	AddOnChatScrapFinishedListenner(IWppScrapperChatScrapFinishedListener)

	RemoveOnScrapperStartedListenner(IWppScrapperStartedListener)
	RemoveOnScrapperStoppedListenner(IWppScrapperStoppedListener)
	RemoveOnScrapperFinishedListenner(IWppScrapperFinishedListener)
	RemoveOnChatScrapStartedListenner(IWppScrapperChatScrapStartedListener)
	RemoveOnChatScrapFinishedListenner(IWppScrapperChatScrapFinishedListener)
}
