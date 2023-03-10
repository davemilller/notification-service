package control

import (
	"fmt"
	"notification-service/domain"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// initLogger initializes a global zap logger for use throughout repo
// call the logger with zap.S().[Info/Warn/Error/Fatal]() from anywhere
func NewLogger() *zap.SugaredLogger {
	consoleWriter := zapcore.AddSync(os.Stdout)

	// cofigure encoder for how log is printed
	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.ConsoleSeparator = " "

	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(fmt.Sprintf("%s |", t.Format(domain.TimeFormat)))
	})
	encoderConfig.EncodeCaller = zapcore.CallerEncoder(func(ec zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(fmt.Sprintf(" | %s |", ec.TrimmedPath()))
	})
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// make logger
	level := zapcore.DebugLevel // this is the min level for messages to be logged
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, consoleWriter, level),
	)

	zapLogger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(zapLogger)

	// AddCallerSkip(1) makes it so that anything using this interface reports the correct caller (not this file)
	return zapLogger.WithOptions(zap.AddCallerSkip(1)).Sugar()
}
