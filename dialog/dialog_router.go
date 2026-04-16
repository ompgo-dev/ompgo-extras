package dialog

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/ompgo-dev/ompgo/pkg/handle"
	"github.com/ompgo-dev/ompgo/pkg/omp"
	ompdialog "github.com/ompgo-dev/ompgo/pkg/omp/dialog"
	"github.com/ompgo-dev/ompgo/pkg/omp/players"
	"github.com/ompgo-dev/ompgo/pkg/runtime"
)

// DefaultDialogID is used by NewDialogRouter.
const DefaultDialogID int32 = 1000

// DialogContext contains data captured when the dialog was shown.
type DialogContext struct {
	Player   *omp.Player
	Data     any
	DialogID int32
	Title    string
	Body     string
	Button1  string
	Button2  string
	Style    omp.DialogStyle
}

// DialogHandler handles a dialog response for a player. Return true if handled.
type DialogHandler func(ctx *DialogContext, event *omp.DialogResponseEvent) bool

// MessageHandler handles a message-box dialog response.
type MessageHandler func(ctx *DialogContext, accepted bool) bool

// InputHandler handles an input or password dialog response.
type InputHandler func(ctx *DialogContext, accepted bool, input string) bool

// ListHandler handles a list, tablist, or tablist-headers dialog response.
type ListHandler func(ctx *DialogContext, accepted bool, listItem int32, itemText string) bool

type dialogEntry struct {
	ctx     *DialogContext
	handler DialogHandler
}

// DialogRouter dispatches OnDialogResponse for dialogs shown via this router.
// It uses a single dialog ID and tracks the latest dialog per player.
type DialogRouter struct {
	mu       sync.Mutex
	dialogID atomic.Int32
	entries  map[handle.Handle]*dialogEntry
}

var (
	dialogRouterOnce sync.Once
	dialogRouter     *DialogRouter
)

// DialogRouterOptions configures the dialog router.
type DialogRouterOptions struct {
	DialogID int32
}

// NewDialogRouter creates a new DialogRouter using DefaultDialogID.
func NewDialogRouter() *DialogRouter {
	return NewDialogRouterWithID(DefaultDialogID)
}

// NewDialogRouterWithID creates a new DialogRouter with a custom dialog ID.
func NewDialogRouterWithID(dialogID int32) *DialogRouter {
	if dialogID == 0 {
		dialogID = DefaultDialogID
	}
	router := &DialogRouter{
		entries: map[handle.Handle]*dialogEntry{},
	}
	router.dialogID.Store(dialogID)
	return router
}

// GetDialogRouter returns the singleton DialogRouter.
func GetDialogRouter() *DialogRouter {
	dialogRouterOnce.Do(func() {
		dialogRouter = NewDialogRouter()
		_ = runtime.RegisterOnDialogResponse(func(ctx context.Context, event *omp.DialogResponseEvent) error {
			_ = dialogRouter.Handle(event)
			return nil
		})
	})
	return dialogRouter
}

// WithDialogRouter configures the singleton router (dialog ID, etc).
func WithDialogRouter(opts DialogRouterOptions) runtime.Option {
	return runtime.WithSetup(func(ctx context.Context) error {
		_ = ctx
		router := GetDialogRouter()
		if opts.DialogID != 0 {
			router.SetDialogID(opts.DialogID)
		}
		return nil
	})
}

// Show displays a dialog and stores the handler and data for the player.
func (r *DialogRouter) Show(player *omp.Player, style omp.DialogStyle, title string, body string, button1 string, button2 string, data any, handler DialogHandler) bool {
	if r == nil || player == nil || !player.Valid() || handler == nil {
		return false
	}
	dialogID := r.dialogID.Load()
	if !ompdialog.Show(player, dialogID, int32(style), title, body, button1, button2) {
		return false
	}

	ctx := &DialogContext{
		Player:   player,
		Data:     data,
		DialogID: dialogID,
		Title:    title,
		Body:     body,
		Button1:  button1,
		Button2:  button2,
		Style:    style,
	}
	populateDialogContext(ctx)

	entry := &dialogEntry{ctx: ctx, handler: handler}

	r.mu.Lock()
	r.entries[player.Handle()] = entry
	r.mu.Unlock()

	return true
}

// ShowMessage shows a DIALOG_STYLE_MSGBOX dialog.
func (r *DialogRouter) ShowMessage(player *omp.Player, title string, body string, button1 string, button2 string, data any, handler MessageHandler) bool {
	if handler == nil {
		return false
	}
	return r.Show(player, omp.DialogStyleMsgBox, title, body, button1, button2, data, func(ctx *DialogContext, event *omp.DialogResponseEvent) bool {
		return handler(ctx, dialogAccepted(event))
	})
}

// ShowInput shows a DIALOG_STYLE_INPUT dialog.
func (r *DialogRouter) ShowInput(player *omp.Player, title string, body string, button1 string, button2 string, data any, handler InputHandler) bool {
	if handler == nil {
		return false
	}
	return r.Show(player, omp.DialogStyleInput, title, body, button1, button2, data, func(ctx *DialogContext, event *omp.DialogResponseEvent) bool {
		return handler(ctx, dialogAccepted(event), dialogInputText(event))
	})
}

// ShowList shows a DIALOG_STYLE_LIST dialog.
func (r *DialogRouter) ShowList(player *omp.Player, title string, body string, button1 string, button2 string, data any, handler ListHandler) bool {
	if handler == nil {
		return false
	}
	return r.Show(player, omp.DialogStyleList, title, body, button1, button2, data, func(ctx *DialogContext, event *omp.DialogResponseEvent) bool {
		return handler(ctx, dialogAccepted(event), dialogListItem(event), dialogInputText(event))
	})
}

// ShowPassword shows a DIALOG_STYLE_PASSWORD dialog.
func (r *DialogRouter) ShowPassword(player *omp.Player, title string, body string, button1 string, button2 string, data any, handler InputHandler) bool {
	if handler == nil {
		return false
	}
	return r.Show(player, omp.DialogStylePassword, title, body, button1, button2, data, func(ctx *DialogContext, event *omp.DialogResponseEvent) bool {
		return handler(ctx, dialogAccepted(event), dialogInputText(event))
	})
}

// ShowTabList shows a DIALOG_STYLE_TABLIST dialog.
func (r *DialogRouter) ShowTabList(player *omp.Player, title string, body string, button1 string, button2 string, data any, handler ListHandler) bool {
	if handler == nil {
		return false
	}
	return r.Show(player, omp.DialogStyleTablist, title, body, button1, button2, data, func(ctx *DialogContext, event *omp.DialogResponseEvent) bool {
		return handler(ctx, dialogAccepted(event), dialogListItem(event), dialogInputText(event))
	})
}

// ShowTabListHeaders shows a DIALOG_STYLE_TABLIST_HEADERS dialog.
func (r *DialogRouter) ShowTabListHeaders(player *omp.Player, title string, body string, button1 string, button2 string, data any, handler ListHandler) bool {
	if handler == nil {
		return false
	}
	return r.Show(player, omp.DialogStyleTablistHeaders, title, body, button1, button2, data, func(ctx *DialogContext, event *omp.DialogResponseEvent) bool {
		return handler(ctx, dialogAccepted(event), dialogListItem(event), dialogInputText(event))
	})
}

// Hide clears the router state for a player and hides their current dialog.
func (r *DialogRouter) Hide(player *omp.Player) bool {
	if r == nil || player == nil || !player.Valid() {
		return false
	}

	r.mu.Lock()
	delete(r.entries, player.Handle())
	r.mu.Unlock()

	return ompdialog.Hide(player)
}

// Current returns the active router-managed dialog context for a player.
func (r *DialogRouter) Current(player *omp.Player) (*DialogContext, bool) {
	if r == nil || player == nil || !player.Valid() {
		return nil, false
	}

	r.mu.Lock()
	entry := r.entries[player.Handle()]
	r.mu.Unlock()
	if entry == nil || entry.ctx == nil {
		return nil, false
	}

	ctxCopy := *entry.ctx
	if !populateDialogContext(&ctxCopy) {
		return nil, false
	}
	if ctxCopy.DialogID != entry.ctx.DialogID {
		return nil, false
	}
	return &ctxCopy, true
}

// SetDialogID sets the dialog ID used by the router.
func (r *DialogRouter) SetDialogID(dialogID int32) {
	if r == nil {
		return
	}
	if dialogID == 0 {
		dialogID = DefaultDialogID
	}
	r.dialogID.Store(dialogID)
}

// Handle routes an OnDialogResponse event. Returns true if handled.
func (r *DialogRouter) Handle(event *omp.DialogResponseEvent) bool {
	if r == nil || event == nil || event.Player == nil || !event.Player.Valid() {
		return false
	}

	var entry *dialogEntry
	r.mu.Lock()
	entry = r.entries[event.Player.Handle()]
	if entry != nil && entry.ctx != nil && entry.ctx.DialogID != event.DialogID {
		delete(r.entries, event.Player.Handle())
		r.mu.Unlock()
		return false
	}
	delete(r.entries, event.Player.Handle())
	r.mu.Unlock()
	if entry == nil || entry.handler == nil {
		return false
	}

	return entry.handler(entry.ctx, event)
}

func populateDialogContext(ctx *DialogContext) bool {
	if ctx == nil || ctx.Player == nil || !ctx.Player.Valid() {
		return false
	}

	var dialogID int32
	var style int32
	var title string
	var body string
	var button1 string
	var button2 string
	if !players.GetDialogData(ctx.Player, &dialogID, &style, &title, &body, &button1, &button2) {
		return false
	}

	ctx.DialogID = dialogID
	ctx.Style = omp.DialogStyle(style)
	ctx.Title = title
	ctx.Body = body
	ctx.Button1 = button1
	ctx.Button2 = button2
	return true
}

func dialogAccepted(event *omp.DialogResponseEvent) bool {
	return event != nil && event.Response != 0
}

func dialogInputText(event *omp.DialogResponseEvent) string {
	if event == nil {
		return ""
	}
	return event.InputText.String()
}

func dialogListItem(event *omp.DialogResponseEvent) int32 {
	if event == nil {
		return -1
	}
	return event.ListItem
}
