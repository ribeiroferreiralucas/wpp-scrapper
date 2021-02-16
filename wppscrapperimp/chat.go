package wppscrapperimp

import wppscrapper "github.com/ribeiroferreiralucas/wpp-scrapper"

type Chat struct {
	id     string
	name   string
	desc   string
	status wppscrapper.ChatStatus
}

func (c *Chat) Name() string {
	return c.name
}

func (c *Chat) Jid() string {
	return c.id
}

func (c *Chat) GetStatus() wppscrapper.ChatStatus {
	return c.status
}
