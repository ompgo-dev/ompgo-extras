## Dialog routing

Use `DialogRouter` to show dialogs with a managed dialog ID and handle responses based on the current dialog shown to each player.
The router auto-registers the response handler, so you don't need to implement `OnDialogResponse` yourself.

Helpers are available for the common dialog styles: `ShowMessage`, `ShowInput`, `ShowList`, `ShowPassword`, `ShowTabList`, and `ShowTabListHeaders`.

- `ShowMessage` handlers receive only `accepted`.
- `ShowInput` and `ShowPassword` handlers receive `accepted` and `input`.
- `ShowList`, `ShowTabList`, and `ShowTabListHeaders` handlers receive `accepted`, `listItem`, and `itemText`.

### Runtime bootstrap

Register the package by passing `dialog.WithDialogRouter(...)` to `runtime.Bootstrap`:

```go
import (
	"github.com/ompgo-dev/ompgo-extras/dialog"
	"github.com/ompgo-dev/ompgo/pkg/runtime"
)

func main() {
	runtime.Bootstrap(
		runtime.WithGamemode(func() runtime.Gamemode { return &MyGamemode{} }),
		dialog.WithDialogRouter(dialog.DialogRouterOptions{DialogID: 1000}),
	)
}
```

```go
import "github.com/ompgo-dev/ompgo-extras/dialog"

var dialogs = dialog.GetDialogRouter()

func (gm *MyGamemode) OnPlayerConnect(event *gamemode.PlayerConnectEvent) {
	if event.Player == nil || !event.Player.Valid() {
		return
	}
	dialogs.ShowInput(event.Player, "Login", "Enter your password:", "Login", "Cancel", "login", func(ctx *dialog.DialogContext, accepted bool, input string) bool {
		// ctx.Data == "welcome"
		// ctx.DialogID is the actual ID used for this dialog.
		if !accepted {
			return true
		}
		_ = input
		return true
	})
}
```
