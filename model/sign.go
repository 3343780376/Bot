package model

import (
	"github.com/3343780376/go-bot"
	"log"
	"strconv"
	"time"
)

type Sign struct {
	UserId string `json:"user_id" db:"userId"`
	Num    int    `json:"num" db:"num"`
	Card   string `json:"card" db:"card"`
	Day    int    `json:"day" db:"day"`
}

func (con *Connect) SelctSign(event go_bot.Event) int {
	con.Exist(event)
	var sign Sign
	err := con.Db.Get(&sign, "select * from sign where userId=?", strconv.Itoa(event.UserId))
	if err != nil {
		log.Println(err)
		return 0
	} else {
		return sign.Num
	}
}

func (con *Connect) Exist(event go_bot.Event) {
	var sign Sign
	err := con.Db.Get(&sign, "select * from sign where userId=?", strconv.Itoa(event.UserId))
	if err != nil {
		_, _ = con.Db.Exec("insert into sign (userId, num, card, day) VALUES (?,?,?,?)",
			strconv.Itoa(event.UserId), 10, event.Sender.Card, time.Now().Day()-1)
	}
}

/*
	return:
		已经签到 ：true
		未签到： false
*/
func (con *Connect) IsSign(event go_bot.Event) bool {
	con.Exist(event)
	var sign Sign
	_ = con.Db.Get(&sign, "select * from sign where userId=?", strconv.Itoa(event.UserId))
	if sign.Day == time.Now().Day() {
		return true
	} else {
		con.Update(2, event)
		con.Db.Exec("update sign set day = ? where userId = ?", time.Now().Day(), strconv.Itoa(event.UserId))
		return false
	}
}

func (con *Connect) Update(n int, event go_bot.Event) {
	con.Exist(event)
	if n > 0 {
		_, err := con.Db.Exec("update sign set num = num + ? where userId = ?", n, strconv.Itoa(event.UserId))
		if err != nil {
			log.Println(err)
		}
	} else {
		_, err := con.Db.Exec("update sign set num = num - ? where userId = ?", -n, strconv.Itoa(event.UserId))
		if err != nil {
			log.Println(err)
		}
	}

}
