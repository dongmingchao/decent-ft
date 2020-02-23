package scraper

import (
	"fmt"
	"net/url"
)

type Error struct {
	ErrCode uint
	ErrMsg  string
	ErrUrl  url.URL
}

func (err *Error) String() string {
	return fmt.Sprintf("%s [%d] %s", err.ErrUrl.String(), err.ErrCode, err.ErrMsg)
}

func (err *Error) Error() string {
	return err.String()
}
