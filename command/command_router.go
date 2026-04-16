package command

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/ompgo-dev/ompgo/pkg/omp"
	"github.com/ompgo-dev/ompgo/pkg/runtime"
)

// CommandHandler handles a command for a player. Return true if the command was handled.
type CommandHandler func(player *omp.Player, args []string) bool

// CommandRouter dispatches OnPlayerCommandText to registered handlers.
// Commands are matched case insensitive.
type CommandRouter struct {
	mu       sync.Mutex
	snapshot atomic.Pointer[commandRouterSnapshot]
}

type commandRouterSnapshot struct {
	prefix   string
	handlers map[string]CommandHandler
}

var (
	commandRouterOnce sync.Once
	commandRouter     *CommandRouter
)

// CommandRouterOptions configures the command router.
type CommandRouterOptions struct {
	Prefix string
}

// NewCommandRouter creates a new CommandRouter.
func NewCommandRouter() *CommandRouter {
	return NewCommandRouterWithPrefix("/")
}

// NewCommandRouterWithPrefix creates a new CommandRouter with a custom command prefix.
func NewCommandRouterWithPrefix(prefix string) *CommandRouter {
	if prefix == "" {
		prefix = "/"
	}
	router := &CommandRouter{}
	router.snapshot.Store(&commandRouterSnapshot{
		prefix:   prefix,
		handlers: map[string]CommandHandler{},
	})
	return router
}

func (r *CommandRouter) loadSnapshot() *commandRouterSnapshot {
	if r == nil {
		return nil
	}
	snapshot := r.snapshot.Load()
	if snapshot != nil {
		return snapshot
	}
	snapshot = &commandRouterSnapshot{
		prefix:   "/",
		handlers: map[string]CommandHandler{},
	}
	r.snapshot.CompareAndSwap(nil, snapshot)
	return r.snapshot.Load()
}

func cloneHandlers(src map[string]CommandHandler) map[string]CommandHandler {
	if len(src) == 0 {
		return map[string]CommandHandler{}
	}
	dst := make(map[string]CommandHandler, len(src))
	for cmd, handler := range src {
		dst[cmd] = handler
	}
	return dst
}

// GetCommandRouter returns the singleton CommandRouter.
func GetCommandRouter() *CommandRouter {
	commandRouterOnce.Do(func() {
		commandRouter = NewCommandRouter()
		_ = runtime.RegisterOnPlayerCommandText(func(ctx context.Context, event *omp.PlayerCommandTextEvent) (bool, error) {
			_ = ctx
			return commandRouter.Handle(event), nil
		})
	})
	return commandRouter
}

// WithCommandRouter configures the singleton router (prefix, etc).
func WithCommandRouter(opts CommandRouterOptions) runtime.Option {
	return runtime.WithSetup(func(ctx context.Context) error {
		_ = ctx
		router := GetCommandRouter()
		if opts.Prefix != "" {
			router.SetPrefix(opts.Prefix)
		}
		return nil
	})
}

// Register adds a handler for a command (e.g. "/help" or "help").
// It returns a function that unregisters the handler.
func (r *CommandRouter) Register(command string, handler CommandHandler) func() {
	if r == nil || handler == nil {
		return func() {}
	}
	snapshot := r.loadSnapshot()
	cmd := normalizeRegisteredCommand(command, snapshot.prefix)
	if cmd == "" {
		return func() {}
	}

	r.mu.Lock()
	snapshot = r.loadSnapshot()
	handlers := cloneHandlers(snapshot.handlers)
	handlers[cmd] = handler
	r.snapshot.Store(&commandRouterSnapshot{
		prefix:   snapshot.prefix,
		handlers: handlers,
	})
	r.mu.Unlock()

	return func() {
		r.mu.Lock()
		snapshot := r.loadSnapshot()
		if _, ok := snapshot.handlers[cmd]; !ok {
			r.mu.Unlock()
			return
		}
		handlers := cloneHandlers(snapshot.handlers)
		delete(handlers, cmd)
		r.snapshot.Store(&commandRouterSnapshot{
			prefix:   snapshot.prefix,
			handlers: handlers,
		})
		r.mu.Unlock()
	}
}

// SetPrefix sets the command prefix (default "/").
func (r *CommandRouter) SetPrefix(prefix string) {
	if r == nil {
		return
	}
	if prefix == "" {
		prefix = "/"
	}
	r.mu.Lock()
	snapshot := r.loadSnapshot()
	r.snapshot.Store(&commandRouterSnapshot{
		prefix:   prefix,
		handlers: snapshot.handlers,
	})
	r.mu.Unlock()
}

// Handle routes an OnPlayerCommandText event. Returns true if handled.
func (r *CommandRouter) Handle(event *omp.PlayerCommandTextEvent) bool {
	if r == nil || event == nil || event.Player == nil || !event.Player.Valid() {
		return false
	}
	snapshot := r.loadSnapshot()
	if !event.Command.HasPrefix(snapshot.prefix) {
		return false
	}
	cmd, args := splitCommand(event.Command.String(), snapshot.prefix)
	if cmd == "" {
		return false
	}

	handler := snapshot.handlers[cmd]
	if handler == nil {
		return false
	}
	return handler(event.Player, args)
}

func normalizeRegisteredCommand(command string, prefix string) string {
	if prefix == "" {
		prefix = "/"
	}
	command = strings.TrimSpace(command)
	if command == "" {
		return ""
	}
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return ""
	}
	command = fields[0]
	command = strings.TrimPrefix(command, prefix)
	command = strings.TrimSpace(command)
	if command == "" {
		return ""
	}
	return strings.ToLower(command)
}

func splitCommand(input string, prefix string) (string, []string) {
	if prefix == "" {
		prefix = "/"
	}
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", nil
	}
	if !strings.HasPrefix(trimmed, prefix) {
		return "", nil
	}
	fields := strings.Fields(trimmed)
	if len(fields) == 0 {
		return "", nil
	}
	cmd := normalizeRegisteredCommand(fields[0], prefix)
	if cmd == "" {
		return "", nil
	}
	return cmd, fields[1:]
}
