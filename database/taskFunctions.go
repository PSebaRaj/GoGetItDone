package database

import (
	"fmt"

	"github.com/psebaraj/gogetitdone/models"
)

var TYPE_TASK string = "tasks"
var TYPE_EXPIRINGTASK string = "expiring_tasks"
var TYPE_PRIORITYTASK string = "priority_tasks"

//
// Functions for all types of tasks
//

// toggles the complete boolean for any task type
func ToggleTaskComplete(taskType string, taskID uint, currState bool) {
	psqlStatement := `
	UPDATE %s
	SET complete = $2
	WHERE id = $1;`

	toggleStatement := fmt.Sprintf(psqlStatement, taskType)

	err = DB.Exec(toggleStatement, taskID, !currState).Error
	if err != nil {
		fmt.Printf("Error marking %s #%d as %t", taskType, taskID, !currState)
		panic(err)
	}

}

// changes title of any task
func ChangeTaskTitle(taskType string, taskID uint, newTitle string) {
	psqlStatement := `
	UPDATE %s
	SET title = $2
	WHERE id = $1;`

	changeStatement := fmt.Sprintf(psqlStatement, taskType)

	err = DB.Exec(changeStatement, taskID, newTitle).Error
	if err != nil {
		fmt.Printf("Error changing title of %s #%d to %s", taskType, taskID, newTitle)
		panic(err)
	}

}

// changes description of any task
func ChangeTaskDescription(taskType string, taskID uint, newDescription string) {
	psqlStatement := `
	UPDATE %s
	SET description = $2
	WHERE id = $1;`

	changeStatement := fmt.Sprintf(psqlStatement, taskType)

	err = DB.Exec(changeStatement, taskID, newDescription).Error
	if err != nil {
		fmt.Printf("Error changing title of %s #%d to %s", taskType, taskID, newDescription)
		panic(err)
	}

}

//
// Functions for priority tasks
//

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
