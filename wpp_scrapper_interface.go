package wppscrapper

type IWppScrapper interface {
	Auth(qrChan chan<- string) (string, error)
	ReAuth(qrChan chan<- string, uuid string) (string, error)
	WaitInitialization() chan bool
	Initialized() bool
	GetChats() map[string]Chat
	StartScrapper(resume bool)
	StopScrapper()
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
