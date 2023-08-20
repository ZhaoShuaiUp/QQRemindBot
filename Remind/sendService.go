package Remind

type SendService interface {
	SendRemind()
	AddRemind(remind Remind)
	AddRemindList(remindList []Remind)
}
