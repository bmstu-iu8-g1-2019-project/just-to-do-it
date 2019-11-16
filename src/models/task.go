package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Tags struct {
	Id int       `json:"id"` // ???
	Title string `json:"title"`
	Color string `json:"color"`
}

type Task struct {
	Id               int       `json:"id"`
	AssigneeId       int       `json:"assignee_id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	State            string    `json:"state"`
	Deadline         time.Time `json:"deadline"`
	Priority         int       `json:"priority"`
	CreationDatetime time.Time `json:"creation_datetime"`
	GroupId          int       `json:"group_id"`
	Tag              []Tags    `json:"tags"`
}

func (t *Tags) Value() (driver.Value, error) {
	return fmt.Sprintf("(%d,'%s','%s')", t.Id, t.Title, t.Color), nil
}

type TagsArr []Tags

func (t *TagsArr) Scan(src interface{}) error {
	b := src.([]byte)
	err := json.Unmarshal(b, t)
	if err != nil {
		return err
	}
	return nil
	//log.Println(string(src.([]byte)))
	//s := string(src.([]byte))
	//// Убираем скобки.
	//s = s[1 : len(s)-1]
	//
	//// Разделяем.
	//parts := strings.Split(s, ",")
	//
	//// Парсим части.
	//var err error
	//t.Id, err = strconv.Atoi(parts[0])
	//if err != nil {
	//	return err
	//}
	//
	//t.Title = parts[1]
	//t.Color = parts[2]
	//
	//return nil
}

