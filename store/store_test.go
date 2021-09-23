package store

import (
	"reflect"
	"testing"
)

func TestGetPendingTasks(t *testing.T) {
	t.Log("getPendingTasks")

	ds := Datastore{
		tasks: []Task{
			{1, "Do housework", true},
			{2, "Buy milk", false},
		},
	}

	got := ds.GetPendingTasks()
	want := []Task{ds.tasks[1]}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("getPendingTasks()=%v, want %v", got, want)
	}
}

func TestSaveTask(t *testing.T) {
	t.Log("saveTask")

	var tests = []struct {
		name string
		ds   *Datastore
		task Task
		want []Task
		err  error
	}{
		{
			name: "Should save new task in datastore",
			ds:   &Datastore{},
			task: Task{Title: "Buy milk"},
			want: []Task{
				{1, "Buy milk", false},
			},
		},
		{
			name: "Should update existing task in the store",
			ds: &Datastore{
				tasks: []Task{
					{1, "Buy milk", false},
				},
			},
			task: Task{Id: 1, Title: "Buy milk", Done: true},
			want: []Task{
				{1, "Buy milk", true},
			},
		},
		{
			name: "Should error when task Id does not exist",
			ds:   &Datastore{},
			task: Task{Id: 1, Title: "Buy milk", Done: true},
			err:  ErrTaskNotFound,
		},
	}

	for _, testcase := range tests {
		t.Log(testcase.name)
		testcase.ds.SaveTask(testcase.task)
		if !reflect.DeepEqual(testcase.ds.tasks, testcase.want) {
			t.Errorf("saveTask()=%v, want %v", testcase.ds.tasks, testcase.want)
		}
	}
}
