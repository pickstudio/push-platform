package model

import (
	"time"

	oapiv1 "github.com/pickstudio/push-platform/api/oapi/v1"
	"github.com/pickstudio/push-platform/pkg/er"
	"github.com/pkg/errors"
)

var (
	ErrParsingPlainViewTitleIsNotValid        = errors.New("error parsing 'plain_view', 'title' is not valid")
	ErrParsingPlainViewContentIsNotValid      = errors.New("error parsing 'plain_view', 'content' is not valid")
	ErrParsingPlainViewThumbnailURLIsNotValid = errors.New("error parsing 'plain_view', 'thumbnail_url' is not valid")
	ErrParsingPlainViewSchemaURLIsNotValid    = errors.New("error parsing 'plain_view', 'schema_url' is not valid")
	ErrParsingPlainViewAlarmIsNotValid        = errors.New("error parsing 'plain_view', 'alarm' is not valid")
	ErrParsingPlainViewCreatedAtIsNotValid    = errors.New("error parsing 'plain_view', 'created_at' is not valid")
)

type PlainView struct {
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	ThumbnailURL string    `json:"thumbnail_url"`
	SchemeURL    string    `json:"scheme_url"`
	Alarm        string    `json:"alarm"`
	CreatedAt    time.Time `json:"created_at"`
}

func (v *PlainView) Render() error {
	return nil
}

func ParseOAPIPlainView(v *oapiv1.PlainView) (*PlainView, error) {
	if v.Title == "" {
		return nil, ErrParsingPlainViewTitleIsNotValid
	}
	if v.Content == "" {
		return nil, ErrParsingPlainViewContentIsNotValid
	}
	if v.ThumbnailUrl == "" {
		return nil, ErrParsingPlainViewThumbnailURLIsNotValid
	}
	if v.SchemeUrl == "" {
		return nil, ErrParsingPlainViewSchemaURLIsNotValid
	}
	if v.Alarm == "" {
		return nil, ErrParsingPlainViewAlarmIsNotValid
	}
	if v.CreatedAt == "" {
		return nil, ErrParsingPlainViewCreatedAtIsNotValid
	}
	createdAt, err := time.Parse(time.RFC3339, v.CreatedAt)
	if err != nil {
		return nil, er.WithSourceErr(ErrParsingPlainViewCreatedAtIsNotValid, err)
	}
	return &PlainView{
		Title:        v.Title,
		Content:      v.Content,
		ThumbnailURL: v.ThumbnailUrl,
		SchemeURL:    v.SchemeUrl,
		Alarm:        v.Alarm,
		CreatedAt:    createdAt,
	}, nil
}
