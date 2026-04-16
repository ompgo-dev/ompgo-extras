## Command routing

Use `CommandRouter` to dispatch `OnPlayerCommandText` in a consistent way.
The router auto-registers the event handler, so you don't need to implement `OnPlayerCommandText` yourself.

### Runtime bootstrap

Register the package by passing `command.WithCommandRouter(...)` to `runtime.Bootstrap`:

```go
import (
	"github.com/ompgo-dev/ompgo-extras/command"
	"github.com/ompgo-dev/ompgo/pkg/runtime"
)

func main() {
	runtime.Bootstrap(
		runtime.WithGamemode(func() runtime.Gamemode { return &MyGamemode{} }),
		command.WithCommandRouter(command.CommandRouterOptions{Prefix: "/"}),
	)
}
```

```go
import (
	"github.com/ompgo-dev/ompgo-extras/command"
	"github.com/ompgo-dev/ompgo/pkg/omp"
	"github.com/ompgo-dev/ompgo/pkg/omp/players"
)

var commands = command.GetCommandRouter()

func init() {
	commands.Register("/help", func(p *omp.Player, args []string) bool {
		_ = players.SendClientMessage(p, uint32(omp.ColorWhite), "Available commands: /help, /heal")
		return true
	})
}
```
