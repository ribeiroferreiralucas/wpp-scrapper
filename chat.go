package wppscrapper

import whatsapp "github.com/Rhymen/go-whatsapp"

type Chat struct {
	handler *ChatHandler
	wppChat *whatsapp.Chat
	status  ChatStatus // TODO: Implementar todas as trocas de status
}

type ChatStatus int

const (
	Idle     = -1
	Queue    = 0
	Running  = 1
	Stoped   = 2
	Finished = 3
)

func (c *Chat) Name() string {
	return c.wppChat.Name
}

func (c *Chat) Jid() string {
	return c.wppChat.Jid
}

func (c *Chat) GetStatus() ChatStatus {
	return c.status
}

func (c *Chat) GetChatInfo() {
	//TODO: Return others chat info
	// return c.chatHandler
}
