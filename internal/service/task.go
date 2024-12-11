package service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
	"task-queue/internal/models"
	"time"
)

type TaskService struct {
	currentTasks map[uint32]*models.Task
	queue        chan *models.Task
	mutex        sync.RWMutex
	idCounter    uint32
}

func NewTaskService() *TaskService {
	return &TaskService{
		currentTasks: make(map[uint32]*models.Task),
		queue:        make(chan *models.Task),
		mutex:        sync.RWMutex{},
		idCounter:    0,
	}
}

func (s *TaskService) generateTaskID() uint32 {
	return atomic.AddUint32(&s.idCounter, 1)
}

func (s *TaskService) Create(input string) uint32 {
	id := s.generateTaskID()
	task := &models.Task{
		ID:        id,
		InputData: input,
		Status:    "pending",
	}
	s.mutex.Lock()
	s.currentTasks[id] = task
	s.mutex.Unlock()

	s.queue <- task
	logrus.Infof("task created, id: %d", id)
	fmt.Println("aueeeee: ", s.currentTasks)
	return id
}

func (s *TaskService) Get(id uint32) (models.Task, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	task, ok := s.currentTasks[id]
	if !ok {
		return models.Task{}, fmt.Errorf("Задача с id %d не найдена", id)
	}
	logrus.Infof("task gotted, id: %d", id)
	return *task, nil
}

func (s *TaskService) StartProcessing() {
	for task := range s.queue {
		time.Sleep(2 * time.Second)
		s.mutex.Lock()
		task.Status = "completed"
		task.Result = fmt.Sprintf("Результат: %s", task.InputData)
		s.mutex.Unlock()
		logrus.Infof("task completed, id: %d", task.ID)
	}
}
