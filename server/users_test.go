package server

import (
	"github.com/glawscorp/glawscord/db"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setup(t *testing.T) string {
	tmpdir := "/tmp/glawscord/"
	err := os.Mkdir(tmpdir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	dbPath, err := os.CreateTemp(tmpdir, "test-*.db")
	if err != nil {
		t.Fatal(err)
	}

	//nolint:errcheck
	dbPath.Close()
	err = db.InitDB(dbPath.Name())
	if err != nil {
		t.Fatal(err)
	}

	err = db.CreateUser("user1", "2Cornishhens")
	if err != nil {
		t.Fatal(err)
	}

	err = db.CreateUser("user2", "AndSomeWater")
	if err != nil {
		t.Fatal(err)
	}

	return dbPath.Name()
}

func TestGetUsers(t *testing.T) {
	dbPath := setup(t)
	//nolint:errcheck
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
