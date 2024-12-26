package interfaces

type CronService interface {
	StartCronJobs() error
	StopCronJobs()
}
