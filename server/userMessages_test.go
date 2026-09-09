package server

import (
	"encoding/json"
	"github.com/glawscorp/glawscord/db"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupMessages(t *testing.T) {

	err := db.CreateUserMessage(1, 2, "louie.LOUIEE.Go and see your father.")
	if err != nil {
		t.Fatal(err)
	}
	err = db.CreateUserMessage(2, 1, "Uncle X I-")
	if err != nil {
		t.Fatal(err)
	}
	err = db.CreateUserMessage(1, 2, "Shame on you louie. Shame on you.")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetMessages(t *testing.T) {
	dbPath := Setup(t)
	setupMessages(t)

	defer os.Remove(dbPath)

	r := InitServer()
	var msgs []*db.UserMessage

	t.Run("returns all messages exchanged between 2 users", func(t *testing.T) {

		request := httptest.NewRequest(http.MethodGet, "/users/{ID}/messages?sender_id=1&receiver_id=2&limit=10&offset=0", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, request)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
			return
		}

		err := json.NewDecoder(w.Body).Decode(&msgs)
		if err != nil {
			t.Errorf("unable to parse json response: %v", err)
			return
		}

		if len(msgs) != 3 {
			t.Errorf("received unexpected number of messages: %v", len(msgs))
			return
		}

		if msgs[0].Content != "louie.LOUIEE.Go and see your father." {
			t.Error("received unexpected message content")
			return
		}

		if msgs[0].Sender != 1 {
			t.Errorf("received message from unexpected sender: %v", msgs[0].Sender)
			return
		}

	})

}
