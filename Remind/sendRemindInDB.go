package Remind

import (
	"RemindBot/BotApi/dto"
	"RemindBot/BotApi/openapi"
	"context"
	"fmt"
	"time"
)

type SendRemindInDB struct {
	remindChan chan Remind
	api        *openapi.OpenAPI
}

func NewSendRemindInDB(api *openapi.OpenAPI) *SendRemindInDB {
	service := &SendRemindInDB{
		remindChan: make(chan Remind, 5000),
		api:        api,
	}
	return service
}

func (s *SendRemindInDB) AddRemind(remind Remind) {
	s.remindChan <- remind
}

func (s *SendRemindInDB) AddRemindList(remindList []Remind) {
	for i := range remindList {
		s.remindChan <- remindList[i]
	}
}

func (s *SendRemindInDB) SendRemind() {
	for r := range s.remindChan {
		t1, err := time.Parse(LayoutDateTime, r.Date+" "+r.Time)
		if err != nil {
			continue
		}
		t2, _ := time.Parse(LayoutDateTime, time.Now().Format(LayoutDateTime))
		d := t1.Sub(t2)
		if d > 0 {
			time.Sleep(d)
		}

		remind := r
		go func() {
			var ctx context.Context
			reply, err := s.api.PostMessage(ctx, remind.Cid, (&dto.PostMessage{
				Content: remind.Content,
			}).AddAtUsr(remind.Uid))
			if err != nil {
				fmt.Printf("[Send Message] failed, remind: %v, err: %v\n", remind, err)
			}
			fmt.Printf("[reply] reply: %v", reply)
		}()
	}
}
