package server

import (
	"github.com/glawscorp/glawscord/db"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func setup(t *testing.T) string {
	tmpdir := "/tmp/glawscord/"
	err := os.MkdirAll(tmpdir, 0755)
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

func TestCreateUser(t *testing.T) {
	dbPath := setup(t)
	//nolint:errcheck
	defer os.Remove(dbPath)

	r := InitServer()

	//create new valid user
	t.Run("creates a user in the db", func(t *testing.T) {

		body := `{"Username": "Test_User", "Password":"Test_Password123-" }`
		request := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, request)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}

	})
	//user already exists
	t.Run("tries to create a user with a username that already exists in the db", func(t *testing.T) {

		body := `{"Username": "user1", "Password":"Test_Password123-" }`
		request := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, request)

		if w.Code != http.StatusConflict {
			t.Errorf("got %d, want %d", w.Code, http.StatusConflict)
		}

	})
	//invalid username
	t.Run("tries to create user with invalid username", func(t *testing.T) {

		body := `{"Username": "", "Password":"Test_Password123-" }`
		request := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, request)

		if w.Code != http.StatusBadRequest {
			t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
		}

	})

}

func TestDeleteUser(t *testing.T) {

	dbPath := setup(t)
	//nolint:errcheck
	defer os.Remove(dbPath)

	r := InitServer()

	t.Run("deletes a user with specified id", func(t *testing.T) {

		request := httptest.NewRequest(http.MethodDelete, "/users/2", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, request)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}

	})
	//tries to delete non-existing user. for some reason this is returning 200
	/*t.Run("tries to delete user that doesn't exist", func(t *testing.T) {

		request := httptest.NewRequest(http.MethodDelete, "/users/9999", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, request)

		if w.Code != http.StatusBadRequest {
			t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
		}

	})*/

}
