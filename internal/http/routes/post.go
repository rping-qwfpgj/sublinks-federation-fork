package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"participating-online/sublinks-federation/internal/activitypub"
	"participating-online/sublinks-federation/internal/lemmy"

	"github.com/gorilla/mux"
)

func SetupPostRoutes(r *mux.Router) {
	r.HandleFunc("/post/{postId}", getPostHandler).Methods("GET")
}

func getPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := context.Background()
	c := lemmy.GetLemmyClient(ctx)
	post, err := c.GetPost(ctx, vars["postId"])
	if err != nil {
		log.Println("Error reading post", err)
		return
	}
	postLd := activitypub.ConvertPostToApub(post)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
	content, _ := json.MarshalIndent(postLd, "", "  ")
	w.Write(content)
}