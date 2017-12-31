package app

/*
#cgo LDFLAGS: -L${SRCDIR}/../../bin -lplatform
#cgo windows LDFLAGS: -luser32 -lole32 -ld2d1 -ldwrite -lstdc++
#include "platform/platform.h"
#include <stdlib.h>

void goPaintCallback(app_t *app, int width, int height);
void goResizeCallback(app_t *app, int width, int height);
void goMouseMoveCallback(app_t *app, int x, int y);
void goMouseEventCallback(app_t *app, int button, int event, int width, int height);

static void init_go_callbacks(app_t *app)
{
	app_set_on_paint_callback(app, goPaintCallback);
	app_set_on_resize_callback(app, goResizeCallback);
	app_set_on_mouse_move(app, goMouseMoveCallback);
	app_set_on_mouse_event(app, goMouseEventCallback);
}
*/
import "C"

import (
	"errors"
	"unsafe"
)

type Application interface {
	OnPaint(w, h int)
	OnResize(w, h int)
	OnMouseMove(x, y float32)
	OnMouseEvent(b MouseButton, e MouseEvent, x, y float32)
}

type Brush uint32
type Font int32

type Color struct {
	R, G, B, A float32
}

type Rect struct {
	Left, Top, Right, Bottom float32
}

type Vec struct {
	X, Y float32
}

type MouseButton int
type MouseEvent int
type FontWeight uint32
type Alignment uint32

const (
	TransparentBrush uint32 = C.APP_TRANSPARENT_BRUSH
	WhiteBrush              = C.APP_WHITE_BRUSH
	BlackBrush              = C.APP_BLACK_BRUSH

	FontWeightNormal FontWeight = C.APP_FONT_WEIGHT_NORMAL
	FontWeightNarrow            = C.APP_FONT_WEIGHT_NARROW
	FontWeightBold              = C.APP_FONT_WEIGHT_BOLD

	AlignLeft    Alignment = C.APP_ALIGN_LEFT
	AlignRight             = C.APP_ALIGN_RIGHT
	AlignVCenter           = C.APP_ALIGN_VCENTER
	AlignTop               = C.APP_ALIGN_TOP
	AlignBottom            = C.APP_ALIGN_BOTTOM
	AlignHCenter           = C.APP_ALIGN_HCENTER

	ButtonLeft   MouseButton = C.APP_BUTTON_LEFT
	ButtonMiddle             = C.APP_BUTTON_MIDDLE
	ButtonRight              = C.APP_BUTTON_RIGHT

	Press   MouseEvent = C.APP_PRESS
	Release            = C.APP_RELEASE

	NoFont Font = -1
)

var internalApp *C.app_t
var externalApp Application

func NewRect(left, top, right, bottom float32) Rect {
	return Rect{left, top, right, bottom}
}

func NewRectI(left, top, right, bottom int) Rect {
	return Rect{float32(left), float32(top), float32(right), float32(bottom)}
}

func (r *Rect) Move(x, y float32) {
	r.Left += x
	r.Right += x
	r.Top += y
	r.Bottom += y
}

func (r Rect) Hit(x, y float32) bool {
	return r.Left < x && r.Top < y && r.Right > x && r.Bottom > y
}

func (b Brush) Destory() {
	if internalApp != nil {
		C.app_destroy_brush(internalApp, C.app_brush_t(b))
	}
}

func (f Font) Destroy() {
	if internalApp != nil {
		C.app_destroy_font(internalApp, C.app_font_t(f))
	}
}

func Init(title string, width, height int) error {
	if internalApp != nil {
		return errors.New("We only support one running application")
	}

	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))

	internalApp = C.app_init(cTitle, C.int(width), C.int(height))
	C.init_go_callbacks(internalApp)

	return nil
}

func Shutdown() {
	if internalApp != nil {
		C.app_shutdown(internalApp)
		internalApp = nil
	}
}

func Run() {
	if internalApp != nil {
		C.app_run(internalApp)
	}
}

func SetApplication(app Application) {
	externalApp = app
}

func CreateSolidBrush(color Color) Brush {
	if internalApp != nil {
		return Brush(C.app_create_solid_brush(internalApp, *(*C.app_color_t)(unsafe.Pointer(&color))))
	}

	return Brush(0)
}

func CreateFont(familyName string, pointSize float32, fontWeight FontWeight) Font {
	if internalApp != nil {
		name := C.CString(familyName)
		defer C.free(unsafe.Pointer(name))
		return Font(C.app_create_font(internalApp, name, C.float(pointSize), C.uint32_t(fontWeight)))
	}
	return Font(-1)
}

func Repaint() {
	if internalApp != nil {
		C.app_repaint(internalApp)
	}
}

func DrawLine(startX, startY, endX, endY float32, brush Brush, thickness float32) {
	if internalApp != nil {
		C.app_draw_line(internalApp, C.float(startX), C.float(startY), C.float(endX), C.float(endY), C.app_brush_t(brush), C.float(thickness))
	}
}

func DrawRect(rect Rect, brush Brush, strokeThickness float32) {
	if internalApp != nil {
		C.app_draw_rectangle(internalApp, *(*C.app_rect_t)(unsafe.Pointer(&rect)), C.app_brush_t(brush), C.float(strokeThickness))
	}
}

func DrawRoundedRect(rect Rect, radius float32, brush Brush, strokeThickness float32) {
	if internalApp != nil {
		C.app_draw_rounded_rectangle(internalApp, *(*C.app_rect_t)(unsafe.Pointer(&rect)), C.float(radius), C.app_brush_t(brush), C.float(strokeThickness))
	}
}

// void app_draw_text(app_t * app, const char * text, app_font_t font, app_brush_t brush, app_rect_t bounds, uint32_t alignment);
func DrawText(text string, font Font, brush Brush, bounds Rect, alignment Alignment) {
	if internalApp != nil {
		t := C.CString(text)
		defer C.free(unsafe.Pointer(t))
		C.app_draw_text(internalApp, t, C.app_font_t(font), C.app_brush_t(brush), *(*C.app_rect_t)(unsafe.Pointer(&bounds)), C.uint32_t(alignment))
	}
}

func FillRect(rect Rect, brush Brush) {
	if internalApp != nil {
		C.app_fill_rectangle(internalApp, *(*C.app_rect_t)(unsafe.Pointer(&rect)), C.app_brush_t(brush))
	}
}

func FillRoundedRect(rect Rect, radius float32, brush Brush) {
	if internalApp != nil {
		C.app_fill_rounded_rectangle(internalApp, *(*C.app_rect_t)(unsafe.Pointer(&rect)), C.float(radius), C.app_brush_t(brush))
	}
}

//export goPaintCallback
func goPaintCallback(app *C.app_t, width, height C.int) {
	if externalApp != nil {
		externalApp.OnPaint(int(width), int(height))
	}
}

//export goResizeCallback
func goResizeCallback(app *C.app_t, width, height C.int) {
	if externalApp != nil {
		externalApp.OnResize(int(width), int(height))
	}
}

//export goMouseMoveCallback
func goMouseMoveCallback(app *C.app_t, x, y C.int) {
	if externalApp != nil {
		externalApp.OnMouseMove(float32(x), float32(y))
	}
}

//export goMouseEventCallback
func goMouseEventCallback(app *C.app_t, button, event, x, y C.int) {
	if externalApp != nil {
		externalApp.OnMouseEvent(MouseButton(button), MouseEvent(event), float32(x), float32(y))
	}
}
