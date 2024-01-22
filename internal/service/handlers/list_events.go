package handlers

import "net/http"

func ListEvents(w http.ResponseWriter, r *http.Request) {
	did := r.Header.Get("X-User-DID")

	balance := getBalanceByDID(did, w, r)
	if balance == nil {
		return
	}

	//events, err := EventsQ(r).Select()
}
