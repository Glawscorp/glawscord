package server

import (
	"github.com/glawscorp/glawscord/db"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setup(t *testing.T) string {
	dbPath, err := os.CreateTemp("/tmp/glawscord/", "test-*.db")
	if err != nil {
		t.Fatal(err)
	}
	dbPath.Close()
	db.InitDB(dbPath.Name())
	db.CreateUser("user1", "2Cornishhens")
	db.CreateUser("user2", "AndSomeWater")

	return dbPath.Name()
}

func TestGetUsers(t *testing.T) {
	dbPath := setup(t)
	defer os.Remove(dbPath)

	r := InitServer()

	t.Run("returns all users in db", func(t *testing.T) {

		request := httptest.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, request)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
		}

	})

}
