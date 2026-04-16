## TextDraw builder

Use the textdraw builder to construct reusable textdraws with options:

```go
import (
	"github.com/ompgo-dev/ompgo-extras/textdraw"
	"github.com/ompgo-dev/ompgo/pkg/omp"
)

title := textdraw.Create(
	320.0, 60.0, "Welcome",
	textdraw.WithFont(omp.TextDrawFontPricedown),
	textdraw.WithAlignment(omp.TextDrawAlignCenter),
	textdraw.WithColor(omp.ColorWhite),
	textdraw.WithLetterSize(0.6, 2.0),
)

_ = textdraw.ShowForPlayer(player, title)
```
