package model

import (
	"encoding/json"
	"fmt"
	"github.com/GarnBarn/common-go/model"
	"github.com/GarnBarn/common-go/proto"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func convertReminterTimeToString(reminterTime []int) string {
	reminderTimeByte, _ := json.Marshal(reminterTime)
	return strings.Trim(string(reminderTimeByte), "[]")
}

func ToTagPublic(tag model.Tag, maskSecretKey bool) TagPublic {
	reminderTime := strings.Split(tag.ReminderTime, ",")
	reminterTimeInt := []int{}

	for _, item := range reminderTime {
		result, err := strconv.Atoi(item)
		if err != nil {
			logrus.Warn("Can't convert the result to int: ", item, " for ", tag.ID)
			continue
		}
		reminterTimeInt = append(reminterTimeInt, result)
	}

	secretKey := ""
	if !maskSecretKey {
		secretKey = tag.SecretKeyTotp
	}

	return TagPublic{
		ID:            fmt.Sprint(tag.ID),
		Name:          tag.Name,
		Author:        tag.Author,
		Color:         tag.Color,
		ReminderTime:  reminterTimeInt,
		Subscriber:    strings.Split(tag.Subscriber, ","),
		SecretKeyTotp: secretKey,
	}
}

type TagPublic struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Author        string   `json:"author"`
	Color         string   `json:"color"`
	ReminderTime  []int    `json:"reminderTime"`
	Subscriber    []string `json:"subscriber"`
	SecretKeyTotp string   `json:"secretKeyTotp,omitempty"`
}

func (tp *TagPublic) ToProtoTag() proto.TagPublic {
	var reminderTimes []int32
	for _, time := range tp.ReminderTime {
		reminderTimes = append(reminderTimes, int32(time))
	}

	return proto.TagPublic{
		Id:            tp.ID,
		Name:          tp.Name,
		Author:        tp.Author,
		Color:         tp.Color,
		ReminderTime:  reminderTimes,
		Subscriber:    tp.Subscriber,
		SecretKeyTotp: tp.SecretKeyTotp,
	}
}

type CreateTagRequest struct {
	Name         string   `json:"name" validate:"required"`
	Color        string   `json:"color"`
	ReminderTime []int    `json:"reminderTime,omitempty" validate:"omitempty,max=3"`
	Subscriber   []string `json:"subscriber"`
}

func (ct *CreateTagRequest) ToTag(author string) model.Tag {
	return model.Tag{
		Name:         ct.Name,
		Author:       author,
		Color:        ct.Color,
		ReminderTime: convertReminterTimeToString(ct.ReminderTime),
		Subscriber:   strings.Join(ct.Subscriber, ","),
	}
}

type UpdateTagRequest struct {
	Name         *string   `json:"name,omitempty"`
	Color        *string   `json:"color,omitempty"`
	ReminderTime *[]int    `json:"reminderTime,omitempty" validate:"omitempty,max=3"`
	Subscriber   *[]string `json:"subscribe"`
}

func (utr *UpdateTagRequest) UpdateTag(tag *model.Tag) {
	if utr.Name != nil {
		tag.Name = *utr.Name
	}

	if utr.Color != nil {
		tag.Color = *utr.Color
	}

	if utr.ReminderTime != nil {
		tag.ReminderTime = convertReminterTimeToString(*utr.ReminderTime)
	}

	if utr.Subscriber != nil {
		tag.Subscriber = strings.Join(*utr.Subscriber, ",")
	}
}
