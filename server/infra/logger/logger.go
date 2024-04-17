package logger

type keyLogger int

// LoggerKey is the key used to retrieve the logger value from the request context.
const LoggerKey keyLogger = 0

// Logger is an interface that describes all the capabilities of the application's logger.
type Logger interface {
	With(args ...any) Logger

	WithoutCaller() Logger

	Debug(args ...any)

	// Info uses fmt.Sprint to construct and log a message at INFO level
	Info(args ...any)

	// Error uses fmt.Sprint to construct and log a message at ERROR level
	Error(args ...any)

	// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
	Fatal(args ...any)

	// Debugf uses fmt.Sprintf to construct and log a message at DEBUG level
	Debugf(format string, args ...any)

	// Infof uses fmt.Sprintf to construct and log a message at INFO level
	Infof(format string, args ...any)

	// Errorf uses fmt.Sprintf to construct and log a message at ERROR level
	Errorf(format string, args ...any)

	// Fatalf uses fmt.Sprintf to construct and log a message, then calls os.Exit.
	Fatalf(format string, args ...any)
}
