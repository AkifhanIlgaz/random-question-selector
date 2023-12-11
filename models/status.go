package models

type Status int

const (
	StatusFail Status = iota
	StatusSuccess
)

func StatusMessage(status Status) string {
	switch status {
	case StatusFail:
		return "fail"
	case StatusSuccess:
		return "success"
	default:
		return "invalid status"
	}
}
