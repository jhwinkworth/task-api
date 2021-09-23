package store

import "errors"

type Task struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type Datastore struct {
	tasks  []Task
	lastId int
}

// ErrTaskNotFound is returned when a Task ID is not found
var ErrTaskNotFound = errors.New("task was not found")

func (ds *Datastore) GetPendingTasks() []Task {
	var pendingTasks []Task
	for _, task := range ds.tasks {
		if !task.Done {
			pendingTasks = append(pendingTasks, task)
		}
	}
	return pendingTasks
}

func (ds *Datastore) SaveTask(task Task) error {
	if task.Id == 0 {
		ds.lastId++
		task.Id = ds.lastId
		ds.tasks = append(ds.tasks, task)
		return nil
	}
	for i, t := range ds.tasks {
		if t.Id == task.Id {
			ds.tasks[i] = task
			return nil
		}
	}
	return ErrTaskNotFound
}
