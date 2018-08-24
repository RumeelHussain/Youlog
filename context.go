<<<<<<< HEAD
package terminallogging
=======
package termlog
>>>>>>> dabee0d... Add Streams

import "golang.org/x/net/context"

// NewContext creates a new context with an included Logger
func NewContext(ctx context.Context, logger Logger) context.Context {
<<<<<<< HEAD
	return context.WithValue(ctx, "terminallogging", logger)
=======
	return context.WithValue(ctx, "termlog", logger)
>>>>>>> dabee0d... Add Streams
}

// FromContext retrieves a Logger from a context. If no logger is present, we
// return a new silenced logger that will produce no output.
func FromContext(ctx context.Context) Logger {
<<<<<<< HEAD
	logger, ok := ctx.Value("terminallogging").(Logger)
=======
	logger, ok := ctx.Value("termlog").(Logger)
>>>>>>> dabee0d... Add Streams
	if !ok {
		l := NewLog()
		l.Quiet()
		return l
	}
	return logger
}
