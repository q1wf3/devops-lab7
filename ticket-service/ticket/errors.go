package ticket

import "errors"

var (
	ErrTicketNotFound    = errors.New("ticket not found")
	ErrTicketAlreadySold = errors.New("ticket already sold")
	ErrCmdRepository     = errors.New("unable to command repository")
	ErrQueryRepository   = errors.New("unable to query repository")
)
