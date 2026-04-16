package textdraw

import (
	"github.com/ompgo-dev/ompgo/pkg/omp"
	"github.com/ompgo-dev/ompgo/pkg/omp/textdraw"
)

// TextDrawOption configures a textdraw during build.
type TextDrawOption func(*textDrawConfig)

type textdrawLetter struct {
	letterSizeSet bool    // letter size is the size of each character in the textdraw.
	letterSizeX   float32 // letterSizeY is the vertical size of each character in the textdraw.
	letterSizeY   float32 // letterSizeX is the horizontal size of each character in the textdraw.
}

type textdrawText struct {
	textSizeSet bool    // text size is the size of the text block. It is not necessarily the same as letter size * number of letters, as there may be extra spacing.
	textSizeX   float32 // textSizeY is the vertical size of the text block.
	textSizeY   float32 // textSizeX is the horizontal size of the text block.
}

type textdrawAlignment struct {
	alignmentSet bool              // alignment is the text alignment within the text block.
	alignment    omp.TextDrawAlign // alignment can be left, center or right.
}

type textdrawColor struct {
	colorSet bool      // color is the text color.
	color    omp.Color // color is a 32-bit RGBA color value, where the highest byte is the alpha component and the lowest byte is the blue component.
}

type textdrawBox struct {
	useBoxSet bool // useBox determines whether the textdraw has a box around it.
	useBox    bool // useBoxSet is true if useBox was set.

	boxColorSet bool      // boxColor is the color of the box around the textdraw, if useBox is enabled.
	boxColor    omp.Color // boxColor is a 32-bit RGBA color value, where the highest byte is the alpha component and the lowest byte is the blue component.
}

type textdrawShadow struct {
	shadowSet bool  // shadow is the size of the shadow. A value of 0 means no shadow, while higher values create a larger shadow.
	shadow    int32 // shadowSet is true if shadow was set.
}

type textdrawOutline struct {
	outlineSet bool  // outline is the size of the outline. A value of 0 means no outline, while higher values create a larger outline.
	outline    int32 // outlineSet is true if outline was set.
}

type textdrawBackground struct {
	backgroundColorSet bool      // backgroundColor is the color of the background behind the text. It is only visible if the text has a shadow or outline, as it is drawn behind them.
	backgroundColor    omp.Color // backgroundColor is a 32-bit RGBA color value, where the highest byte is the alpha component and the lowest byte is the blue component.
}

type textdrawFont struct {
	fontSet bool             // font is the font used for the text. The available fonts are defined in the gamemode package and may vary depending on the game mode.
	font    omp.TextDrawFont // fontSet is true if font was set.
}

type textdrawProportional struct {
	proportionalSet bool // proportional determines whether the textdraw uses proportional spacing. If true, the spacing between characters is adjusted based on the character width, creating a more natural look. If false, all characters have the same spacing, which can make the text look more uniform but less visually appealing.
	proportional    bool // proportionalSet is true if proportional was set.
}

type textdrawSelectable struct {
	selectableSet bool // selectable determines whether the textdraw can be selected by the player. If true, the player can interact with the textdraw by clicking on it, which can trigger events in the game mode. If false, the textdraw is purely decorative and cannot be interacted with.
	selectable    bool // selectableSet is true if selectable was set.
}

type textdrawString struct {
	textSet bool   // text is the string displayed by the textdraw. It can contain multiple lines if it includes newline characters (\n).
	text    string // textSet is true if text was set.
}

type textdrawPreviewModel struct {
	previewModelSet bool  // previewModel is the model ID used for the preview in the textdraw. If set, the textdraw will display a 3D model preview instead of text. The model ID corresponds to the models available in the game, and can be used to show vehicles, objects, or characters as part of the textdraw.
	previewModel    int32 // previewModelSet is true if previewModel was set.
}

type textdrawPreviewRotation struct {
	previewRotSet bool    // previewRotX, previewRotY, previewRotZ and previewZoom define the rotation and zoom level of the 3D model preview in the textdraw. The rotation is specified in degrees for each axis (X, Y, Z), while the zoom is a scaling factor that determines how large the model appears in the textdraw. These settings allow for customizing the appearance of the model preview to fit the design of the textdraw.
	previewRotX   float32 // previewRotX is the rotation around the X-axis in degrees.
	previewRotY   float32 // previewRotY is the rotation around the Y-axis in degrees.
	previewRotZ   float32 // previewRotZ is the rotation around the Z-axis in degrees.
	previewZoom   float32 // previewZoom is the zoom level for the model preview, where 1.0 is the default size.
}

type textdrawPreviewVehicleColors struct {
	previewVehColSet bool  // previewVehCol1 and previewVehCol2 define the primary and secondary colors for vehicle previews in the textdraw. These colors are used when the preview model is a vehicle, allowing for customization of the vehicle's appearance in the textdraw. The color values correspond to the color IDs used in the game, which determine the specific colors applied to the vehicle's primary and secondary color slots.
	previewVehCol1   int32 // previewVehCol1 is the primary color ID for vehicle previews.
	previewVehCol2   int32 // previewVehCol2 is the secondary color ID for vehicle previews.
}

type textdrawPosition struct {
	posSet bool    // posX and posY define the position of the textdraw on the screen. The coordinates are specified in pixels, where (0, 0) is the top-left corner of the screen. Setting the position allows for moving the textdraw to a specific location after it has been created, providing flexibility in its placement within the game interface.
	posX   float32 // posX is the horizontal position of the textdraw in pixels.
	posY   float32 // posY is the vertical position of the textdraw in pixels.
}

type textDrawConfig struct {
	textdrawLetter
	textdrawText
	textdrawAlignment
	textdrawColor
	textdrawBox
	textdrawShadow
	textdrawOutline
	textdrawBackground
	textdrawFont
	textdrawProportional
	textdrawSelectable
	textdrawString
	textdrawPreviewModel
	textdrawPreviewRotation
	textdrawPreviewVehicleColors
	textdrawPosition
}

// TextDrawBuilder builds reusable textdraws.
type TextDrawBuilder struct {
	x    float32
	y    float32
	text string
	opts []TextDrawOption
}

// NewTextDrawBuilder creates a builder for a textdraw.
func NewTextDrawBuilder(x, y float32, text string, opts ...TextDrawOption) *TextDrawBuilder {
	return &TextDrawBuilder{x: x, y: y, text: text, opts: opts}
}

// Build creates the textdraw and applies all options.
func (b *TextDrawBuilder) Build() *omp.TextDraw {
	if b == nil {
		return nil
	}
	return Create(b.x, b.y, b.text, b.opts...)
}

// Create creates a textdraw with options without constructing a builder explicitly.
func Create(x, y float32, text string, opts ...TextDrawOption) *omp.TextDraw {
	return CreateWithID(x, y, text, nil, opts...)
}

// CreateWithID creates a textdraw, optionally storing the underlying open.mp textdraw ID.
func CreateWithID(x, y float32, text string, id *int32, opts ...TextDrawOption) *omp.TextDraw {
	td := textdraw.Create(x, y, text, id)
	if td == nil {
		return nil
	}
	Apply(td, opts...)
	return td
}

// Apply applies builder options to an existing textdraw.
func Apply(td *omp.TextDraw, opts ...TextDrawOption) bool {
	if td == nil {
		return false
	}
	cfg := textDrawConfig{}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	applyTextDrawConfig(td, &cfg)
	return true
}

// Destroy removes a textdraw.
func Destroy(td *omp.TextDraw) bool {
	return textdraw.Destroy(td)
}

// ShowForPlayer shows a textdraw to one player.
func ShowForPlayer(player *omp.Player, td *omp.TextDraw) bool {
	return textdraw.ShowForPlayer(player, td)
}

// HideForPlayer hides a textdraw from one player.
func HideForPlayer(player *omp.Player, td *omp.TextDraw) bool {
	return textdraw.HideForPlayer(player, td)
}

// ShowForAll shows a textdraw to all players.
func ShowForAll(td *omp.TextDraw) bool {
	return textdraw.ShowForAll(td)
}

// HideForAll hides a textdraw from all players.
func HideForAll(td *omp.TextDraw) bool {
	return textdraw.HideForAll(td)
}

// SetString updates the text shown by a textdraw.
func SetString(td *omp.TextDraw, text string) bool {
	return textdraw.SetString(td, text)
}

// SetStringForPlayer overrides the text shown to one player.
func SetStringForPlayer(td *omp.TextDraw, player *omp.Player, text string) bool {
	return textdraw.SetStringForPlayer(td, player, text)
}

func applyTextDrawConfig(td *omp.TextDraw, cfg *textDrawConfig) {
	if cfg.letterSizeSet {
		_ = textdraw.SetLetterSize(td, cfg.letterSizeX, cfg.letterSizeY)
	}
	if cfg.textSizeSet {
		_ = textdraw.SetTextSize(td, cfg.textSizeX, cfg.textSizeY)
	}
	if cfg.alignmentSet {
		_ = textdraw.SetAlignment(td, int32(cfg.alignment))
	}
	if cfg.colorSet {
		_ = textdraw.SetColor(td, uint32(cfg.color))
	}
	if cfg.useBoxSet {
		_ = textdraw.SetUseBox(td, cfg.useBox)
	}
	if cfg.boxColorSet {
		_ = textdraw.SetBoxColor(td, uint32(cfg.boxColor))
	}
	if cfg.shadowSet {
		_ = textdraw.SetShadow(td, cfg.shadow)
	}
	if cfg.outlineSet {
		_ = textdraw.SetOutline(td, cfg.outline)
	}
	if cfg.backgroundColorSet {
		_ = textdraw.SetBackgroundColor(td, uint32(cfg.backgroundColor))
	}
	if cfg.fontSet {
		_ = textdraw.SetFont(td, int32(cfg.font))
	}
	if cfg.proportionalSet {
		_ = textdraw.SetProportional(td, cfg.proportional)
	}
	if cfg.selectableSet {
		_ = textdraw.SetSelectable(td, cfg.selectable)
	}
	if cfg.textSet {
		_ = textdraw.SetString(td, cfg.text)
	}
	if cfg.previewModelSet {
		_ = textdraw.SetPreviewModel(td, cfg.previewModel)
	}
	if cfg.previewRotSet {
		_ = textdraw.SetPreviewRot(td, cfg.previewRotX, cfg.previewRotY, cfg.previewRotZ, cfg.previewZoom)
	}
	if cfg.previewVehColSet {
		_ = textdraw.SetPreviewVehCol(td, cfg.previewVehCol1, cfg.previewVehCol2)
	}
	if cfg.posSet {
		_ = textdraw.SetPos(td, cfg.posX, cfg.posY)
	}
}

// WithLetterSize sets the letter size.
func WithLetterSize(x, y float32) TextDrawOption {
	return func(c *textDrawConfig) {
		c.letterSizeSet = true
		c.letterSizeX = x
		c.letterSizeY = y
	}
}

// WithTextSize sets the text size.
func WithTextSize(x, y float32) TextDrawOption {
	return func(c *textDrawConfig) {
		c.textSizeSet = true
		c.textSizeX = x
		c.textSizeY = y
	}
}

// WithAlignment sets the alignment.
func WithAlignment(align omp.TextDrawAlign) TextDrawOption {
	return func(c *textDrawConfig) {
		c.alignmentSet = true
		c.alignment = align
	}
}

// WithColor sets the text color.
func WithColor(color omp.Color) TextDrawOption {
	return func(c *textDrawConfig) {
		c.colorSet = true
		c.color = color
	}
}

// WithBox enables/disables the box.
func WithBox(use bool) TextDrawOption {
	return func(c *textDrawConfig) {
		c.useBoxSet = true
		c.useBox = use
	}
}

// WithBoxColor sets the box color.
func WithBoxColor(color omp.Color) TextDrawOption {
	return func(c *textDrawConfig) {
		c.boxColorSet = true
		c.boxColor = color
	}
}

// WithShadow sets the shadow size.
func WithShadow(size int32) TextDrawOption {
	return func(c *textDrawConfig) {
		c.shadowSet = true
		c.shadow = size
	}
}

// WithOutline sets the outline size.
func WithOutline(size int32) TextDrawOption {
	return func(c *textDrawConfig) {
		c.outlineSet = true
		c.outline = size
	}
}

// WithBackgroundColor sets the background color.
func WithBackgroundColor(color omp.Color) TextDrawOption {
	return func(c *textDrawConfig) {
		c.backgroundColorSet = true
		c.backgroundColor = color
	}
}

// WithFont sets the font.
func WithFont(font omp.TextDrawFont) TextDrawOption {
	return func(c *textDrawConfig) {
		c.fontSet = true
		c.font = font
	}
}

// WithProportional sets proportional mode.
func WithProportional(set bool) TextDrawOption {
	return func(c *textDrawConfig) {
		c.proportionalSet = true
		c.proportional = set
	}
}

// WithSelectable sets selectable mode.
func WithSelectable(set bool) TextDrawOption {
	return func(c *textDrawConfig) {
		c.selectableSet = true
		c.selectable = set
	}
}

// WithText sets the text string after creation.
func WithText(text string) TextDrawOption {
	return func(c *textDrawConfig) {
		c.textSet = true
		c.text = text
	}
}

// WithPreviewModel sets the preview model.
func WithPreviewModel(model int32) TextDrawOption {
	return func(c *textDrawConfig) {
		c.previewModelSet = true
		c.previewModel = model
	}
}

// WithPreviewRotation sets the preview rotation and zoom.
func WithPreviewRotation(x, y, z, zoom float32) TextDrawOption {
	return func(c *textDrawConfig) {
		c.previewRotSet = true
		c.previewRotX = x
		c.previewRotY = y
		c.previewRotZ = z
		c.previewZoom = zoom
	}
}

// WithPreviewVehicleColors sets preview vehicle colors.
func WithPreviewVehicleColors(color1, color2 int32) TextDrawOption {
	return func(c *textDrawConfig) {
		c.previewVehColSet = true
		c.previewVehCol1 = color1
		c.previewVehCol2 = color2
	}
}

// WithPosition sets the textdraw position after creation.
func WithPosition(x, y float32) TextDrawOption {
	return func(c *textDrawConfig) {
		c.posSet = true
		c.posX = x
		c.posY = y
	}
}
