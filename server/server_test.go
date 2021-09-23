package server

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"task-api/store"
	"testing"
)

func TestGetPendingTasks(t *testing.T) {

	t.Log("GetPendingTasks")

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tasks/pending", nil)

	ds = &mockedStore{}

	// The datastore is restored at the end of the test
	defer func() { ds = &store.Datastore{} }()

	GetPendingTasks(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("GetPendingTasks()=%d, want %d", rec.Code, http.StatusOK)
	}

	want := `[{"id":1,"title":"Do housework","done":false},{"id":2,"title":"Buy milk","done":false}]`

	if got := rec.Body.String(); got != want {
		t.Errorf("GetPendingTasks()=%s, want %s", got, want)
	}
}

func TestAddTask(t *testing.T) {
	t.Log("AddTask")

	var tests = []struct {
		name string
		saveFunc func(task store.Task) error
		body []byte
		want int
	}{
		{
			name: "Should add new task from JSON",
			body: []byte(`{"Title":"Buy bread for breakfast."}`),
			want: http.StatusCreated,
		},
		{
			name: "Should return bad argument when JSON could not be handled",
			body: []byte(""),
			want: http.StatusBadRequest,
		},
		{
			name: "Should return bad argument when datastore returns an error",
			saveFunc: func(task store.Task) error {
				return errors.New("datastore error")
			},
			body: []byte(`{"Title":"Buy bread for breakfast."}`),
			want: http.StatusBadRequest,
		},
		{
			name: "Should return bad argument when task title is empty",
			body: []byte(`{"Title":""}`),
			want: http.StatusBadRequest,
		},
	}

	defer func() { ds = &store.Datastore{} }()

	for _, testcase := range tests {
		t.Log(testcase.name)
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(testcase.body))
		ds = &mockedStore{
			SaveTaskFunc: testcase.saveFunc,
		}
		AddTask(rec, req)
		if rec.Code != testcase.want {
			t.Errorf("AddTask=%d, want %d", rec.Code, testcase.want)
		}
	}
}

func TestUpdateTask(t *testing.T) {

	t.Log("UpdateTask")

	var updateTaskTests = []struct {
		name     string
		saveFunc func(task store.Task) error
		body     []byte
		want     int

	}{
		{
			name: "Should respond with a status 200 OK when the task was updated",
			body: []byte(`{"ID":1, "Title":"Buy bread for breakfast.", "Done":true }`),
			want: http.StatusOK,
		},
		{
			name: "Should respond with a status 400 Bad Request when JSON body could not be handle",
			body: []byte(""),
			want: http.StatusBadRequest,
		},
		{
			name: "Should respond with a status 400 Bad Request when the datastore returned an error",
			saveFunc: func(task store.Task) error {
				return errors.New("datastore error")
			},
			body: []byte(`{"ID":1, "Title":"Buy bread for breakfast.", "Done":true }`),
			want: http.StatusBadRequest,
		},
		{
			name: "Should respond with a status 400 Bad Request when task title is emtpy",
			body: []byte(`{"Title":""}`),
			want: http.StatusBadRequest,
		},
	}

	defer func() { ds = &store.Datastore{} }()

	for _, testcase := range updateTaskTests {
		t.Logf(testcase.name)

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/tasks/1", bytes.NewBuffer(testcase.body))

		ds = &mockedStore{
			SaveTaskFunc: testcase.saveFunc,
		}

		UpdateTask(rec, req)

		if rec.Code != testcase.want {
			t.Errorf("KO => Got %d wanted %d", rec.Code, testcase.want)
		}
	}
}

type mockedStore struct{ SaveTaskFunc func(task store.Task) error }

func (ms *mockedStore) GetPendingTasks() []store.Task {
	return []store.Task{
		{1, "Do housework", false},
		{2, "Buy milk", false},
	}
}

func (ms *mockedStore) SaveTask(task store.Task) error {
	if ms.SaveTaskFunc != nil {
		return ms.SaveTaskFunc(task)
	}
	return nil
}
