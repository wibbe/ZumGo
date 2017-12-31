
#ifndef PLATFORM_H
#define PLATFORM_H

#ifdef __cplusplus
#define EXTERN_C extern "C"
#else
#define EXTERN_C
#endif

#include <stdint.h>

typedef struct app_t app_t;
typedef uint32_t app_brush_t;
typedef uint32_t app_font_t;

#define APP_TRANSPARENT_BRUSH 0
#define APP_WHITE_BRUSH 1
#define APP_BLACK_BRUSH 2

#define APP_FONT_WEIGHT_NORMAL 0
#define APP_FONT_WEIGHT_NARROW 1
#define APP_FONT_WEIGHT_BOLD 2

#define APP_ALIGN_LEFT    0x01
#define APP_ALIGN_RIGHT   0x02
#define APP_ALIGN_VCENTER 0x04
#define APP_ALIGN_TOP     0x10
#define APP_ALIGN_BOTTOM  0x20
#define APP_ALIGN_HCENTER 0x40

#define APP_BUTTON_LEFT   0
#define APP_BUTTON_MIDDLE 1
#define APP_BUTTON_RIGHT  2
#define APP_PRESS   0
#define APP_RELEASE 1


typedef void (* app_on_paint_callback_t)(app_t *, int, int);
typedef void (* app_on_resize_callback_t)(app_t *, int, int);
typedef void (* app_on_mouse_move_callback_t)(app_t *, int x, int y);
typedef void (* app_on_mouse_event_callback_t)(app_t *,  int button, int event, int x, int y);

typedef struct app_color_ {
    float r;
    float g;
    float b;
    float a;
} app_color_t;

typedef struct app_vec_ {
    float x;
    float y;
} app_vec_t;

typedef struct app_rect_ {
    float left;
    float top;
    float right;
    float bottom;
} app_rect_t;



EXTERN_C app_color_t app_color_rgb8(uint8_t r, uint8_t g, uint8_t b);
EXTERN_C app_color_t app_color_rgb(float r, float g, float b);
EXTERN_C app_color_t app_color_rgba(float r, float g, float b, float a);

EXTERN_C app_vec_t app_vec(float x, float y);
EXTERN_C app_vec_t app_vec_add(app_vec_t a, app_vec_t b);
EXTERN_C app_vec_t app_vec_sub(app_vec_t a, app_vec_t b);
EXTERN_C app_vec_t app_vec_mul(app_vec_t a, app_vec_t b);
EXTERN_C app_vec_t app_vec_scale(app_vec_t a, float scale);

EXTERN_C app_rect_t app_rect(float left, float top, float right, float bottom);

EXTERN_C app_t * app_init(const char * title, int width, int height);
EXTERN_C void app_shutdown(app_t* app);
EXTERN_C void app_run(app_t* app);

EXTERN_C void app_repaint(app_t* app);

EXTERN_C void app_set_on_paint_callback(app_t * app, app_on_paint_callback_t callback);
EXTERN_C void app_set_on_resize_callback(app_t * app, app_on_resize_callback_t callback);
EXTERN_C void app_set_on_mouse_event(app_t * app, app_on_mouse_event_callback_t callback);
EXTERN_C void app_set_on_mouse_move(app_t * app, app_on_mouse_move_callback_t callback);

EXTERN_C void app_set_data(app_t * app, void * data);
EXTERN_C void * app_get_data(app_t * app);

EXTERN_C void app_set_background(app_t * app, app_color_t color);

EXTERN_C app_brush_t app_create_solid_brush(app_t * app, app_color_t color);
EXTERN_C void app_destroy_brush(app_t * app, app_brush_t brush);

EXTERN_C app_font_t app_create_font(app_t * app, const char * fontFamily, float pointSize, uint32_t fontWeight);
EXTERN_C void app_destroy_font(app_t * app, app_font_t font);

EXTERN_C void app_draw_line(app_t * app, float startX, float startY, float endX, float endY, app_brush_t brush, float thickness);
EXTERN_C void app_draw_rectangle(app_t * app, app_rect_t rect, app_brush_t brush, float strokeThickness);
EXTERN_C void app_draw_rounded_rectangle(app_t * app, app_rect_t rect, float radius, app_brush_t brush, float strokeThickness);

EXTERN_C void app_fill_rectangle(app_t * app, app_rect_t rect, app_brush_t brush);
EXTERN_C void app_fill_rounded_rectangle(app_t * app, app_rect_t rect, float radius, app_brush_t brush);

EXTERN_C void app_draw_text(app_t * app, const char * text, app_font_t font, app_brush_t brush, app_rect_t bounds, uint32_t alignment);

#endif
