package models

// Enum ...
var Enum = map[string]interface{}{
	"job_type": map[string]JobType{
		"Cron":   JobTypeCron,
		"Manual": JobTypeManual,
	},
	"job_status": map[string]JobStatus{
		"Running": JobRunning,
		"Succeed": JobSucceed,
		"Failed":  JobFailed,
	},
	"task_status": map[string]TaskStatus{
		"Running": TaskRunning,
		"Succeed": TaskSucceed,
		"Failed":  TaskFailed,
	},
}
