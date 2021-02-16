package wppscrapperimp

import (
	bus "github.com/jackhopner/go-events"
	wppscrapper "github.com/ribeiroferreiralucas/wpp-scrapper"
)

type WppScrapperEventHandler struct {
	onScrapperStartedEvent   bus.EventBus
	onScrapperStoppedEvent   bus.EventBus
	onScrapperFinishedEvent  bus.EventBus
	onChatScrapStartedEvent  bus.EventBus
	onChatScrapfinishedEvent bus.EventBus
}

func newWppScrapperEventHandler() *WppScrapperEventHandler {
	return &WppScrapperEventHandler{
		onScrapperStartedEvent:   bus.NewEventBus(),
		onScrapperStoppedEvent:   bus.NewEventBus(),
		onScrapperFinishedEvent:  bus.NewEventBus(),
		onChatScrapStartedEvent:  bus.NewEventBus(),
		onChatScrapfinishedEvent: bus.NewEventBus(),
	}

}

func (w *WppScrapperEventHandler) AddOnScrapperStartedListenner(listener wppscrapper.IWppScrapperStartedListener) {
	w.onScrapperStartedEvent.Register(listener)
}
func (w *WppScrapperEventHandler) AddOnScrapperStoppedListenner(listener wppscrapper.IWppScrapperStoppedListener) {
	w.onScrapperStoppedEvent.Register(listener)
}
func (w *WppScrapperEventHandler) AddOnScrapperFinishedListenner(listener wppscrapper.IWppScrapperFinishedListener) {
	w.onScrapperFinishedEvent.Register(listener)
}
func (w *WppScrapperEventHandler) AddOnChatScrapStartedListenner(listener wppscrapper.IWppScrapperChatScrapStartedListener) {
	w.onChatScrapStartedEvent.Register(listener)
}
func (w *WppScrapperEventHandler) AddOnChatScrapFinishedListenner(listener wppscrapper.IWppScrapperChatScrapFinishedListener) {
	w.onChatScrapfinishedEvent.Register(listener)
}

func (w *WppScrapperEventHandler) RemoveOnScrapperStartedListenner(listener wppscrapper.IWppScrapperStartedListener) {
	w.onScrapperStartedEvent.Deregister(listener)

}
func (w *WppScrapperEventHandler) RemoveOnScrapperStoppedListenner(listener wppscrapper.IWppScrapperStoppedListener) {
	w.onScrapperStoppedEvent.Deregister(listener)
}
func (w *WppScrapperEventHandler) RemoveOnScrapperFinishedListenner(listener wppscrapper.IWppScrapperFinishedListener) {
	w.onScrapperFinishedEvent.Deregister(listener)
}
func (w *WppScrapperEventHandler) RemoveOnChatScrapStartedListenner(listener wppscrapper.IWppScrapperChatScrapStartedListener) {
	w.onChatScrapStartedEvent.Deregister(listener)
}
func (w *WppScrapperEventHandler) RemoveOnChatScrapFinishedListenner(listener wppscrapper.IWppScrapperChatScrapFinishedListener) {
	w.onChatScrapfinishedEvent.Deregister(listener)
}
