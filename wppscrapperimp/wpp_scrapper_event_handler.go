package wppscrapperimp

import (
	wppscrapper "github.com/ribeiroferreiralucas/wpp-scrapper"
)

type WppScrapperEventHandler struct {
	onScrapperStartedListeners   map[wppscrapper.IWppScrapperStartedListener]struct{}
	onScrapperStoppedListeners   map[wppscrapper.IWppScrapperStoppedListener]struct{}
	onScrapperFinishedListeners  map[wppscrapper.IWppScrapperFinishedListener]struct{}
	onChatScrapStartedListeners  map[wppscrapper.IWppScrapperChatScrapStartedListener]struct{}
	onChatScrapfinishedListeners map[wppscrapper.IWppScrapperChatScrapFinishedListener]struct{}
}

func newWppScrapperEventHandler() *WppScrapperEventHandler {
	return &WppScrapperEventHandler{
		onScrapperStartedListeners:   map[wppscrapper.IWppScrapperStartedListener]struct{}{},
		onScrapperStoppedListeners:   map[wppscrapper.IWppScrapperStoppedListener]struct{}{},
		onScrapperFinishedListeners:  map[wppscrapper.IWppScrapperFinishedListener]struct{}{},
		onChatScrapStartedListeners:  map[wppscrapper.IWppScrapperChatScrapStartedListener]struct{}{},
		onChatScrapfinishedListeners: map[wppscrapper.IWppScrapperChatScrapFinishedListener]struct{}{},
	}

}

func (w *WppScrapperEventHandler) AddOnScrapperStartedListener(listener wppscrapper.IWppScrapperStartedListener) {
	w.onScrapperStartedListeners[listener] = struct{}{}
}
func (w *WppScrapperEventHandler) AddOnScrapperStoppedListener(listener wppscrapper.IWppScrapperStoppedListener) {
	w.onScrapperStoppedListeners[listener] = struct{}{}

}
func (w *WppScrapperEventHandler) AddOnScrapperFinishedListener(listener wppscrapper.IWppScrapperFinishedListener) {
	w.onScrapperFinishedListeners[listener] = struct{}{}

}
func (w *WppScrapperEventHandler) AddOnChatScrapStartedListener(listener wppscrapper.IWppScrapperChatScrapStartedListener) {
	w.onChatScrapStartedListeners[listener] = struct{}{}
}
func (w *WppScrapperEventHandler) AddOnChatScrapFinishedListener(listener wppscrapper.IWppScrapperChatScrapFinishedListener) {
	w.onChatScrapfinishedListeners[listener] = struct{}{}
}

func (w *WppScrapperEventHandler) RemoveOnScrapperStartedListener(listener wppscrapper.IWppScrapperStartedListener) {
	delete(w.onScrapperStartedListeners, listener)
}
func (w *WppScrapperEventHandler) RemoveOnScrapperStoppedListener(listener wppscrapper.IWppScrapperStoppedListener) {
	delete(w.onScrapperStoppedListeners, listener)
}
func (w *WppScrapperEventHandler) RemoveOnScrapperFinishedListener(listener wppscrapper.IWppScrapperFinishedListener) {
	delete(w.onScrapperFinishedListeners, listener)
}
func (w *WppScrapperEventHandler) RemoveOnChatScrapStartedListener(listener wppscrapper.IWppScrapperChatScrapStartedListener) {
	delete(w.onChatScrapStartedListeners, listener)
}
func (w *WppScrapperEventHandler) RemoveOnChatScrapFinishedListener(listener wppscrapper.IWppScrapperChatScrapFinishedListener) {
	delete(w.onChatScrapfinishedListeners, listener)
}

func (w *WppScrapperEventHandler) RaiseOnScrapperStartedEvent(wppScrapper wppscrapper.IWppScrapper) {

	for listener := range w.onScrapperStartedListeners {
		if listener == nil {
			continue
		}

		listener.OnWppScrapperStarted(wppScrapper)
	}
}

func (w *WppScrapperEventHandler) RaiseOnScrapperStoppedEvent(wppScrapper wppscrapper.IWppScrapper) {

	for listener := range w.onScrapperStoppedListeners {
		if listener == nil {
			continue
		}

		listener.OnWppScrapperStopped(wppScrapper)
	}
}

func (w *WppScrapperEventHandler) RaiseOnScrapperFinishedEvent(wppScrapper wppscrapper.IWppScrapper) {

	for listener := range w.onScrapperFinishedListeners {
		if listener == nil {
			continue
		}

		listener.OnWppScrapperFinished(wppScrapper)
	}
}

func (w *WppScrapperEventHandler) RaiseOnChatScrapStartedEvent(chat wppscrapper.Chat) {

	for listener := range w.onChatScrapStartedListeners {
		if listener == nil {
			continue
		}

		listener.OnWppScrapperChatScrapStarted(chat)
	}
}

func (w *WppScrapperEventHandler) RaiseOnChatScrapFinishedEvent(chat wppscrapper.Chat) {

	for listener := range w.onChatScrapfinishedListeners {
		if listener == nil {
			continue
		}

		listener.OnWppScrapperChatScrapFinished(chat)
	}
}
