package database

import (
	"errors"
	"fmt"

	"github.com/psebaraj/gogetitdone/models"
)

// mapping between usable task types in Go & their table names in Postgres
var TYPE_TASK string = "tasks"
var TYPE_EXPIRINGTASK string = "expiring_tasks"
var TYPE_PRIORITYTASK string = "priority_tasks"
var allTaskTypes []string = []string{TYPE_TASK, TYPE_EXPIRINGTASK, TYPE_PRIORITYTASK}

func isValidTaskType(givenType string) bool {
	for _, v := range allTaskTypes {
		if v == givenType {
			return true
		}
	}

	return false
}

//
// Functions for all types of tasks
//

// internal: toggles the complete boolean for any task type
func ToggleTaskComplete(taskType string, taskID uint, currState bool) (bool, error) {
	psqlStatement := `
	UPDATE %s
	SET complete = $2
	WHERE id = $1;`

	if !isValidTaskType(taskType) {
		fmt.Printf("Error, %s is not a valid task type", taskType)
		fmt.Println("Valid types: ", allTaskTypes)
		return currState, errors.New("Invalid type")
	}

	toggleStatement := fmt.Sprintf(psqlStatement, taskType)

	err = DB.Exec(toggleStatement, taskID, !currState).Error
	if err != nil {
		fmt.Printf("Error marking %s #%d as %t", taskType, taskID, !currState)
		return currState, err
	}

	return !currState, nil

}

// internal: changes title of any task
func ChangeTaskTitle(taskType string, taskID uint, newTitle string) (string, error) {
	psqlStatement := `
	UPDATE %s
	SET title = $2
	WHERE id = $1;`

	if !isValidTaskType(taskType) {
		fmt.Printf("Error, %s is not a valid task type", taskType)
		fmt.Println("Valid types: ", allTaskTypes)
		return "", errors.New("Invalid type")
	}

	changeStatement := fmt.Sprintf(psqlStatement, taskType)

	err = DB.Exec(changeStatement, taskID, newTitle).Error
	if err != nil {
		fmt.Printf("Error changing title of %s #%d to %s", taskType, taskID, newTitle)
		return newTitle, err
		//panic(err)

	}

	return newTitle, nil
}

// internal: changes description of any task
func ChangeTaskDescription(taskType string, taskID uint, newDescription string) (string, error) {
	psqlStatement := `
	UPDATE %s
	SET description = $2
	WHERE id = $1;`

	if !isValidTaskType(taskType) {
		fmt.Printf("Error, %s is not a valid task type", taskType)
		fmt.Println("Valid types: ", allTaskTypes)
		return "", errors.New("Invalid type")
	}

	changeStatement := fmt.Sprintf(psqlStatement, taskType)

	err = DB.Exec(changeStatement, taskID, newDescription).Error
	if err != nil {
		fmt.Printf("Error changing title of %s #%d to %s", taskType, taskID, newDescription)
		return newDescription, err
		//panic(err)
	}

	return newDescription, nil
}

//
// Functions for priority tasks
//

// internal: changes priority level (uint) of a priorityTask
func ChangePriority(priorityTask models.PriorityTask, newPriorityLevel models.PriorityLevelType) {
	psqlStatement := `
	UPDATE priority_tasks
	SET priority_level = $2
	WHERE id = $1;`

	err = DB.Exec(psqlStatement, priorityTask.ID, newPriorityLevel).Error
	if err != nil {
		fmt.Printf("Error changing priority level of %s to %d", priorityTask.Title, newPriorityLevel)
		panic(err)
	}
}
