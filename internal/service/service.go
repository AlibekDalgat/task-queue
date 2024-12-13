package service

import "task-queue/internal/models"

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Task interface {
	Create(string) (uint32, error)
	Get(uint32) (models.Task, error)
	StartProcessing()
}

type Service struct {
	Task
}

func NewService() *Service {
	return &Service{
		Task: NewTaskService(),
	}
}
