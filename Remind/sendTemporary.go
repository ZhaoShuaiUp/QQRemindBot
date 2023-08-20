package Remind

import (
	"RemindBot/BotApi/dto"
	"RemindBot/BotApi/openapi"
	"container/heap"
	"context"
	"fmt"
	"sync"
	"time"
)

type RemindQueue []Remind

func (rq RemindQueue) Len() int { return len(rq) }

func (rq RemindQueue) Less(i, j int) bool {
	return rq[i].Time < rq[j].Time
}

func (rq RemindQueue) Swap(i, j int) {
	rq[i], rq[j] = rq[j], rq[i]
}

func (rq *RemindQueue) Push(x any) {
	item := x.(Remind)
	*rq = append(*rq, item)
}

func (rq *RemindQueue) Pop() any {
	old := *rq
	n := len(old)
	item := old[n-1]
	*rq = old[0 : n-1]
	return item
}

type SendTemporary struct {
	queue *RemindQueue
	lock  *sync.Mutex
	api   *openapi.OpenAPI
}

func NewSendTemporary(api *openapi.OpenAPI) *SendTemporary {
	return &SendTemporary{
		queue: new(RemindQueue),
		lock:  new(sync.Mutex),
		api:   api,
	}
}

func (s *SendTemporary) AddRemind(remind Remind) {
	s.lock.Lock()
	defer s.lock.Unlock()

	heap.Push(s.queue, remind)
}

func (s *SendTemporary) AddRemindList(remindList []Remind) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for i := range remindList {
		heap.Push(s.queue, remindList[i])
	}
}

func (s *SendTemporary) SendRemind() {
	date := time.Now().Format(LayoutDate)
	for s.queue.Len() > 0 || date <= time.Now().Format(LayoutDate) {
		if s.queue.Len() == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		s.lock.Lock()
		remind := heap.Pop(s.queue).(Remind)
		s.lock.Unlock()

		t1, err := time.Parse(LayoutDateTime, remind.Date+" "+remind.Time)
		if err != nil {
			continue
		}
		t2, _ := time.Parse(LayoutDateTime, time.Now().Format(LayoutDateTime))
		d := t1.Sub(t2)
		if d > 0 {
			time.Sleep(d)
		}

		go func() {
			var ctx context.Context
			reply, err := s.api.PostMessage(ctx, remind.Cid, (&dto.PostMessage{
				Content: "提醒您：" + remind.Content,
			}).AddAtUsr(remind.Uid))
			if err != nil {
				fmt.Printf("[Send Message] failed, remind: %v, err: %v\n", remind, err)
			}
			fmt.Printf("[reply] reply: %v", reply)
		}()
	}
}
