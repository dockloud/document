package document

import "errors"

const MAX_TITLE_LEN = 200
const MAX_BODY_LEN = 250_000

var (
	ErrTitleTooShort = errors.New("title too short")
	ErrTitleTooLong  = errors.New("title too long")
	ErrBodyTooLong   = errors.New("body too long")
)

func ValidateName(title string) error {
	switch {
	case len(title) < 3:
		return ErrTitleTooShort
	case len(title) > MAX_TITLE_LEN:
		return ErrTitleTooLong
	default:
		return nil
	}
}

func ValidateContent(body string) error {
	if len(body) > MAX_BODY_LEN {
		return ErrBodyTooLong
	}
	return nil
}
