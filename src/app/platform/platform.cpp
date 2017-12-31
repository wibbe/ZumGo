
#include "platform.h"

#ifdef _WIN32
#include "platform_win32.cpp"
#endif

EXTERN_C app_color_t app_color_rgb8(uint8_t r, uint8_t g, uint8_t b)
{
    return app_color_t { r / 255.0f, g / 255.0f, b / 255.0f, 1.0f };
}

EXTERN_C app_color_t app_color_rgb(float r, float g, float b)
{
    return app_color_t { r, g, b, 1.0f };
}

EXTERN_C app_color_t app_color_rgba(float r, float g, float b, float a)
{
    return app_color_t { r, g, b, a };
}

EXTERN_C app_vec_t app_vec(float x, float y)
{
    return app_vec_t { x, y };
}

EXTERN_C app_vec_t app_vec_add(app_vec_t a, app_vec_t b)
{
    return app_vec_t { a.x + b.x, a.y + b.y };
}

EXTERN_C app_vec_t app_vec_sub(app_vec_t a, app_vec_t b)
{
    return app_vec_t { a.x - b.x, a.y - b.y };
}

EXTERN_C app_vec_t app_vec_mul(app_vec_t a, app_vec_t b)
{
    return app_vec_t { a.x * b.x, a.y * b.y };
}

EXTERN_C app_vec_t app_vec_scale(app_vec_t a, float scale)
{
    return app_vec_t { a.x * scale, a.y * scale };
}

EXTERN_C app_rect_t app_rect(float left, float top, float right, float bottom)
{
    return app_rect_t { left, top, right, bottom };
}