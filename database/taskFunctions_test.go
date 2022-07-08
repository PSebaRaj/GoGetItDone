package database

import (
	"fmt"
	"testing"

	"github.com/psebaraj/gogetitdone/models"
)

//
// Tests for task existence
//
func TestTaskDoesntExist(t *testing.T) {
	var task models.Task
	err := DB.First(&task, "id = ?", "99999999").Error // entry does not exist
	if err != nil {
		fmt.Println("Threw error when task did not exist")
		return
	}

	t.Fatal("Should have threw error, task does not exist")
}

func TestPriorityTaskDoesntExist(t *testing.T) {
	var priorityTask models.PriorityTask
	err := DB.First(&priorityTask, "id = ?", "99999999").Error // entry does not exist
	if err != nil {
		fmt.Println("Threw error when priority task did not exist")
		return
	}

	t.Fatal("Should have threw error, priority task does not exist")
}

func TestExpiringTaskDoesntExist(t *testing.T) {
	var expiringTask models.ExpiringTask
	err := DB.First(&expiringTask, "id = ?", "99999999").Error // entry does not exist
	if err != nil {
		fmt.Println("Threw error when expiring task did not exist")
		return
	}

	t.Fatal("Should have threw error, expiring task does not exist")
}

// Valid task used for testing:
//{"ID":13,"CreatedAt":"2022-07-07T16:53:36.906713-04:00",
//   "UpdatedAt":"2022-07-07T16:53:36.906713-04:00",
//   "DeletedAt":null,"Title":"test task","Description":"a definition!",
//   "Complete":true,"PersonID":5}

//
// Tests for ToggleTaskComplete
//
func TestToggleTaskIsValid(t *testing.T) {
	var task models.Task
	err := DB.First(&task, "id = ?", "13").Error
	if err != nil {
		t.Fatal("Error finding task before running test")
	}
	preToggle := task.Complete

	postToggle, err := ToggleTaskComplete(TYPE_TASK, task.ID, preToggle)
	if err != nil {
		t.Fatal("Error within internal toggle task complete method")
	}

	if postToggle == !preToggle {
		t.Fatal("Completion bool is not toggled")
	}

	fmt.Println("Successfully toggled completion of task")
	return
}

func TestToggleTaskInvalidTaskType(t *testing.T) {
	var task models.Task
	err := DB.First(&task, "id = ?", "13").Error
	if err != nil {
		t.Fatal("Error finding task before running test")
	}
	preToggle := task.Complete

	_, err = ToggleTaskComplete("some invalid type", task.ID, preToggle)
	if err != nil {
		fmt.Println("Successfully toggled completion of task")
		return
	}

	t.Fatal("Should have threw error, task was of invalid type")
}

//
// Tests for ChangeTaskTitle
//
func TestChangeTaskTitleIsValid(t *testing.T) {
	var task models.Task
	err := DB.First(&task, "id = ?", "13").Error
	if err != nil {
		t.Fatal("Error finding task before running test")
	}
	preTitle := task.Title

	postTitle, err := ChangeTaskTitle(TYPE_TASK, task.ID, "new title")
	if err != nil {
		t.Fatal("Error within internal change title method")
	}

	if postTitle == preTitle {
		t.Fatal("Task title was not changed")
	}

	fmt.Println("Successfully changed title of task")
	return
}

func TestChangeTaskTitleInvalidTaskType(t *testing.T) {
	var task models.Task
	err := DB.First(&task, "id = ?", "13").Error
	if err != nil {
		t.Fatal("Error finding task before running test")
	}
	preTitle := task.Title

	_, err = ChangeTaskTitle("some invalid type", task.ID, preTitle)
	if err != nil {
		fmt.Println("Successfully changed title of task")
		return
	}

	t.Fatal("Should have threw error, task was of invalid type")
}

//
// Tests for ChangeTaskDescription
//
func TestChangeTaskDescriptionIsValid(t *testing.T) {
	var task models.Task
	err := DB.First(&task, "id = ?", "13").Error
	if err != nil {
		t.Fatal("Error finding task before running test")
	}
	preDescription := task.Description

	postDescription, err := ChangeTaskDescription(TYPE_TASK, task.ID, "new description")
	if err != nil {
		t.Fatal("Error within internal change description method")
	}

	if postDescription == preDescription {
		t.Fatal("Task description was not changed")
	}

	fmt.Println("Successfully changed description of task")
	return
}

func TestChangeTaskDescriptionInvalidTaskType(t *testing.T) {
	var task models.Task
	err := DB.First(&task, "id = ?", "13").Error
	if err != nil {
		t.Fatal("Error finding task before running test")
	}
	preDescription := task.Description

	_, err = ChangeTaskDescription("some invalid type", task.ID, preDescription)
	if err != nil {
		fmt.Println("Successfully changed description of task")
		return
	}

	t.Fatal("Should have threw error, task was of invalid type")
}

//
// Tests for ChangePriority
//
