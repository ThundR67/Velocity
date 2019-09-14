package config

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//DefaultErrorHandler is used to log an error and then return a wrappeed version of that error
func DefaultErrorHandler(err error, data ...interface{}) interface{} {
	data = data[0].([]interface{})
	msg := data[0].(string)
	newMsg := msg
	logger := data[1].(*zap.Logger)
	fields := []zap.Field{zap.Error(err)}
	for _, field := range data[2:] {
		toAdd, _ := field.(zap.Field)
		fields = append(fields, toAdd)
		newMsg += " " + toAdd.Key + ":" + toAdd.String
	}
	logger.Error(msg, fields...)
	return errors.Wrap(err, newMsg)
}
