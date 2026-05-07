package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

/*
	要实现超时控制：任务内部必须配合ctx信号。方式1：若底层库原生支持，直接传递ctx即可。方式2：任务内部自行接收ctx.Done。方式3：用子进程管理。
	超时控制策略有两种：取消当前任务、新任务不再调度等。
*/

type Task struct {
	ID       string
	Interval time.Duration
	Timeout  time.Duration
	Job      func(ctx context.Context)
	cancel   context.CancelFunc
}

type Scheduler struct {
	tasks map[string]*Task
	lock  sync.RWMutex
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		tasks: make(map[string]*Task),
	}
}
func (s *Scheduler) AddTask(t *Task) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	// 如果已存在，先删除
	if _, ok := s.tasks[t.ID]; ok {
		return errors.New("task already exists")
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	s.tasks[t.ID] = t

	go s.runTaskCron(ctx, t)
	return nil
}
func (s *Scheduler) UpdateTask(t *Task) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	// 如果已存在，先删除
	if _, ok := s.tasks[t.ID]; !ok {
		return errors.New("the task does not exist")
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	s.tasks[t.ID] = t

	go s.runTaskCron(ctx, t)
	return nil
}
func (s *Scheduler) DelTask(id string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if t, ok := s.tasks[id]; ok {
		t.cancel()
		delete(s.tasks, id)
	}
}
func (s *Scheduler) runTaskCron(ctx context.Context, t *Task) {
	ticker := time.NewTicker(t.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			go s.executeOnce(t)
		}
	}
}
func (s *Scheduler) executeOnce(t *Task) {
	ctx, cancel := context.WithTimeout(context.Background(), t.Timeout)
	defer cancel()

	done := make(chan struct{})

	go func() {
		t.Job(ctx)
		close(done)
	}()

	select {
	case <-ctx.Done():
		fmt.Printf("task %s timeout\n", t.ID) // 感知到任务超时了
	case <-done:
		// 正常完成
	}
}
func main() {
	s := NewScheduler()

	err := s.AddTask(&Task{
		ID:       "task1",
		Interval: 2 * time.Second,
		Timeout:  1 * time.Second,
		Job: func(ctx context.Context) {
			fmt.Println("task1 running")
			time.Sleep(500 * time.Millisecond)
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	// 动态更新周期
	err = s.UpdateTask(&Task{
		ID:       "task1",
		Interval: 1 * time.Second,
		Timeout:  500 * time.Millisecond,
		Job: func(ctx context.Context) {
			fmt.Println("task1 updated")
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	// 删除任务
	s.DelTask("task1")

	time.Sleep(3 * time.Second)
}
