<<<<<<< HEAD
// Package terminallogging provides facilities for logging to a terminal geared towards
// interactive use.
package terminallogging
=======
// Package termlog provides facilities for logging to a terminal geared towards
// interactive use.
package termlog
>>>>>>> 3aae1f8... Remove examples from termlog.go, extend README

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"
)

const defaultTimeFmt = "15:04:05"
const indent = "  "

// Palette defines the colour of output
type Palette struct {
	Timestamp *color.Color
	Say       *color.Color
	Notice    *color.Color
	Warn      *color.Color
	Shout     *color.Color
}

// DefaultPalette is a sensbile default palette, with the following foreground
// colours:
//
// 	Say: Terminal default
// 	Notice: Blue
// 	Warn: Yellow
// 	Shout: Red
// 	Timestamp: Cyan
var DefaultPalette = Palette{
	Say:       color.New(),
	Notice:    color.New(color.FgBlue),
	Warn:      color.New(color.FgYellow),
	Shout:     color.New(color.FgRed),
	Timestamp: color.New(color.FgCyan),
}

// Logger logs things
type Logger interface {
	Say(format string, args ...interface{})
	Notice(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Shout(format string, args ...interface{})

	SayAs(name string, format string, args ...interface{})
	NoticeAs(name string, format string, args ...interface{})
	WarnAs(name string, format string, args ...interface{})
	ShoutAs(name string, format string, args ...interface{})

	Done()
	Quiet()
}

// Stream is a stream of log entries with a header
type Stream interface {
	Logger
	Quiet()
	Header()
}

// terminallogging is the top-level terminallogging interface
type terminallogging interface {
	Logger
	Group() Group
	Stream(header string) Stream
	Quiet()
}

type linesource interface {
	getID() string
	getHeader() string
}

type line struct {
	name   string
	color  *color.Color
	format string
	args   []interface{}
}

// Log is the top-level log structure
type Log struct {
	mu      sync.Mutex
	Palette *Palette
	enabled map[string]bool
	quiet   bool
}

// NewLog creates a new Log instance
func NewLog() *Log {
	l := &Log{
		Palette: &DefaultPalette,
		enabled: make(map[string]bool),
	}
	l.enabled[""] = true
	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		l.NoColor()
	}
	return l
}

// NoColor disables colour output
func (*Log) NoColor() {
	color.NoColor = true
}

// Enable logging for a specified name
func (l *Log) Enable(name string) {
	l.enabled[name] = true
}

// Quiet disables all output
func (l *Log) Quiet() {
	l.quiet = true
}

func (l *Log) output(quiet bool, lines ...*line) {
	if quiet {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(lines) == 0 {
		return
	}
	first := true
	for _, line := range lines {
		if _, ok := l.enabled[line.name]; !ok {
			continue
		}
		var format string
		if first {
			l.Palette.Timestamp.Printf("%s", time.Now().Format(defaultTimeFmt))
			l.Palette.Say.Print(": ")
			first = false
			format = line.format + "\n"
		} else {
			format = indent + line.format + "\n"
		}
		line.color.Printf(format, line.args...)
	}
}

// Say logs a line
func (l *Log) Say(format string, args ...interface{}) {
	l.output(l.quiet, &line{"", l.Palette.Say, format, args})
}

// Notice logs a line with the Notice color
func (l *Log) Notice(format string, args ...interface{}) {
	l.output(l.quiet, &line{"", l.Palette.Notice, format, args})
}

// Warn logs a line with the Warn color
func (l *Log) Warn(format string, args ...interface{}) {
	l.output(l.quiet, &line{"", l.Palette.Warn, format, args})
}

// Shout logs a line with the Shout color
func (l *Log) Shout(format string, args ...interface{}) {
	l.output(l.quiet, &line{"", l.Palette.Shout, format, args})
}

// SayAs logs a line
func (l *Log) SayAs(name string, format string, args ...interface{}) {
	l.output(l.quiet, &line{name, l.Palette.Say, format, args})
}

// NoticeAs logs a line with the Notice color
func (l *Log) NoticeAs(name string, format string, args ...interface{}) {
	l.output(l.quiet, &line{name, l.Palette.Notice, format, args})
}

// WarnAs logs a line with the Warn color
func (l *Log) WarnAs(name string, format string, args ...interface{}) {
	l.output(l.quiet, &line{name, l.Palette.Warn, format, args})
}

// ShoutAs logs a line with the Shout color
func (l *Log) ShoutAs(name string, format string, args ...interface{}) {
	l.output(l.quiet, &line{name, l.Palette.Shout, format, args})
}

// Group creates a new log group
func (l *Log) Group() Logger {
	return &Group{
		palette: l.Palette,
		lines:   make([]*line, 0),
		log:     l,
		quiet:   l.quiet,
	}
}

// Done is just a stub to comply with the Logger interface
func (l *Log) Done() {

}

// Group is a group of lines that constitue a single log entry that won't be
// split. Lines in a group are indented.
type Group struct {
	palette *Palette
	lines   []*line
	log     *Log
	quiet   bool
}

func (g *Group) addLine(name string, color *color.Color, format string, args []interface{}) {
	g.lines = append(g.lines, &line{name, color, format, args})
}

// Say logs a line
func (g *Group) Say(format string, args ...interface{}) {
	g.addLine("", g.palette.Say, format, args)
}

// Notice logs a line with the Notice color
func (g *Group) Notice(format string, args ...interface{}) {
	g.addLine("", g.palette.Notice, format, args)
}

// Warn logs a line with the Warn color
func (g *Group) Warn(format string, args ...interface{}) {
	g.addLine("", g.palette.Warn, format, args)
}

// Shout logs a line with the Shout color
func (g *Group) Shout(format string, args ...interface{}) {
	g.addLine("", g.palette.Shout, format, args)
}

// SayAs logs a line
func (g *Group) SayAs(name string, format string, args ...interface{}) {
	g.addLine(name, g.palette.Say, format, args)
}

// NoticeAs logs a line with the Notice color
func (g *Group) NoticeAs(name string, format string, args ...interface{}) {
	g.addLine(name, g.palette.Notice, format, args)
}

// WarnAs logs a line with the Warn color
func (g *Group) WarnAs(name string, format string, args ...interface{}) {
	g.addLine(name, g.palette.Warn, format, args)
}

// ShoutAs logs a line with the Shout color
func (g *Group) ShoutAs(name string, format string, args ...interface{}) {
	g.addLine(name, g.palette.Shout, format, args)
}

// Done outputs the group to screen
func (g *Group) Done() {
	g.log.output(g.quiet, g.lines...)
}

// Quiet disables output for this subgroup
func (g *Group) Quiet() {
	g.quiet = true
}

// Group of a group is just same group - groups don't nest
func (g *Group) Group() Logger {
	return g
}

// NewContext creates a new context with an included logger
func NewContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, "termlog", logger)
}

// FromContext retrieves a logger from a context. If no logger is present, we
// return a new silenced logger that will produce no output.
func FromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value("termlog").(Logger)
	if !ok {
		l := NewLog()
		l.Quiet()
		return l
	}
	return logger
}

// SetOutput sets the output writer for terminallogging (stdout by default).
func SetOutput(w io.Writer) {
	color.Output = w
}
