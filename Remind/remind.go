package Remind

import (
	"gorm.io/gorm"
	"strings"
	"time"
)

const (
	RemindOnce = iota + 1
	RemindDay
	LayoutDate     = "2006-01-02"
	LayoutTime     = "15:04:05"
	LayoutDateTime = "2006-01-02 15:04:05"
)

type Remind struct {
	// 主键
	Id int `gorm:"column:id;type:int(64) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	// 用户标识
	Uid string `gorm:"column:uid;type:varchar(64);NOT NULL" json:"uid"`
	// QQ频道
	Cid string `gorm:"column:cid;type:varchar(64);NOT NULL" json:"cid"`
	// 提醒类型 1.单次提醒 2.每日提醒
	RemindType uint8 `gorm:"column:remind_type;type:int(8) unsigned;NOT NULL" json:"remind_type"`
	// 提醒日期(单次提醒时不为空)
	Date string `gorm:"column:date;type:varchar(20);" json:"date"`
	// 提醒时间
	Time string `gorm:"column:time;type:varchar(20);NOT NULL" json:"time"`
	// 截至日期(单次提醒时为空)
	EndDate string `gorm:"column:end_date;type:varchar(20);" json:"end_date"`
	// 提醒内容
	Content string `gorm:"column:content;type:varchar(256);NOT NULL" json:"content"`
}

func ExtractMessage(msg string) *Remind {
	var remind *Remind
	info := strings.Split(msg, " ")
	if len(info) == 0 {
		return nil
	}

	// 单次提醒 背单词 2023-08-20 15:30:00
	if info[0] == "单次提醒" && len(info) >= 4 {
		remind = &Remind{}
		remind.RemindType = RemindOnce
		remind.Content = info[1]
		_, err := time.Parse(LayoutDate, info[2])
		if err != nil {
			return nil
		}
		remind.Date = info[2]
		_, err = time.Parse(LayoutTime, info[3])
		if err != nil {
			return nil
		}
		remind.Time = info[3]
		return remind
	}

	// 每日提醒 背单词 15:30:00 2023-09-10
	if info[0] == "每日提醒" && len(info) >= 3 {
		remind = &Remind{}
		remind.RemindType = RemindDay
		remind.Content = info[1]
		_, err := time.Parse(LayoutTime, info[2])
		if err != nil {
			return nil
		}
		remind.Time = info[2]
		if len(info) >= 4 {
			_, err = time.Parse(LayoutDate, info[3])
			if err != nil {
				return nil
			}
			remind.EndDate = info[3]
		}
		return remind
	}

	return remind
}

func (remind *Remind) NeedSendTemporary() bool {
	now := time.Now()
	curDate := now.Format(LayoutDate)
	curTime := now.Format(LayoutTime)
	// 如果今日需要提醒, 特殊处理
	if remind.RemindType == RemindOnce && remind.Date == curDate &&
		remind.Time >= curTime {
		return true
	}
	if remind.RemindType == RemindDay && remind.Time >= curTime {
		return true
	}
	return false
}

func (remind *Remind) Insert(db *gorm.DB) {
	db.Create(remind)
}

func SelectTodayOnce(db *gorm.DB, date string, time string) []Remind {
	var res []Remind
	if time == "" {
		db.Where("remind_type = ? AND date = ?", RemindOnce, date).Find(&res)
	} else {
		db.Where("remind_type = ? AND date = ? AND time > ?", RemindOnce, date, time).Find(&res)
	}
	return res
}

func SelectTodayDay(db *gorm.DB, date string, time string) []Remind {
	var res []Remind
	if time == "" {
		db.Where("remind_type = ? AND end_date >= ?", RemindDay, date).Find(&res)
	} else {
		db.Where("remind_type = ? AND end_date >= ? AND time > ?", RemindDay, date, time).Find(&res)
	}
	return res
}
