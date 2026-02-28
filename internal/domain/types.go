package domain

import "errors"

var ErrEmptyMessage = errors.New("empty log message")

type ViolationCode string

const (
	ViolationLowercase ViolationCode = "logmsg-lowercase"
	ViolationEnglish   ViolationCode = "logmsg-english"
	ViolationNoSpecial ViolationCode = "logmsg-nospecial"
	ViolationSensitive ViolationCode = "logmsg-sensitive"
)

type Violation struct {
	Code    ViolationCode
	Message string
}

type LoggerKind string

const (
	LoggerSlog LoggerKind = "slog"
	LoggerZap  LoggerKind = "zap"
)

type LogCall struct {
	Kind    LoggerKind
	Method  string
	Message string
}
