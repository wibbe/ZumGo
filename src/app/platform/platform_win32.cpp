#include <windows.h>
#include <windowsx.h>

#include <stdio.h>
#include <stdlib.h>
#include <malloc.h>
#include <memory.h>
#include <wchar.h>
#include <cmath>
#include <cstdint>

#include <d2d1.h>
#include <d2d1helper.h>
#include <dwrite.h>
#include <wincodec.h>

#include <vector>
#include <string>

#include "platform.h"


static const int WStrLen = 2048;
static const int StrLen = 2048;

typedef wchar_t WStr[WStrLen];
typedef char Str[StrLen];

template<class Interface>
inline void SafeRelease(Interface **toRelease)
{
   if (*toRelease != nullptr)
   {
      (*toRelease)->Release();
      *toRelease = nullptr;
   }
}

float to_dip_size(float points)
{
   return (points / 72.0f) * 96.0f;
}

std::wstring to_wstr(const char * in)
{
   static const int LEN = 4096;
   static wchar_t BUFFER[LEN];
   int len = MultiByteToWideChar(CP_UTF8, 0, in, -1, BUFFER, LEN);
   return std::wstring(BUFFER, len);
}

static DWRITE_FONT_WEIGHT to_dwrite_font_weight(int weight)
{
   switch (weight)
   {
      case APP_FONT_WEIGHT_NORMAL:
         return DWRITE_FONT_WEIGHT_NORMAL;

      case APP_FONT_WEIGHT_NARROW:
         return DWRITE_FONT_WEIGHT_THIN;

      case APP_FONT_WEIGHT_BOLD:
         return DWRITE_FONT_WEIGHT_BOLD;
   }

   return DWRITE_FONT_WEIGHT_NORMAL;
}

#ifndef HINST_THISCOMPONENT
EXTERN_C IMAGE_DOS_HEADER __ImageBase;
#define HINST_THISCOMPONENT ((HINSTANCE)&__ImageBase)
#endif


struct internal_brush_t
{
   app_color_t color;
   uint32_t id;
   uint32_t next;

   ID2D1SolidColorBrush * brush;


   internal_brush_t()
   {
      color = app_color_rgb(255, 255, 255);
      id = 0;
      next = 0;
      brush = nullptr;
   }

   internal_brush_t(uint32_t id, app_color_t color)
   {
      this->color = color;
      this->id = id;
      this->next = next;
      this->brush = nullptr;
   }
};

struct InternalFont
{
   uint32_t id = 0;
   uint32_t next = 0;
   std::string fontFamily = "";
   float pointSize = 12.0f;
   uint32_t fontWeight = APP_FONT_WEIGHT_NORMAL;

   IDWriteTextFormat * format = nullptr;

   InternalFont(uint32_t id, const char * fontFamily, float pointSize, uint32_t fontWeight)
   {
      this->id = id;
      this->next = 0;
      this->format = nullptr;
      this->fontFamily = fontFamily;
      this->pointSize = pointSize;
      this->fontWeight = fontWeight;
   }
};

struct app_t
{
   app_color_t backgroundColor;
   uint32_t width;
   uint32_t height;
   bool painting;

   void * externalData;

   HWND hwnd;
   ID2D1Factory * factory;
   IDWriteFactory * writeFactory;
   ID2D1HwndRenderTarget * renderTarget;

   app_on_paint_callback_t paintCallback;
   app_on_resize_callback_t resizeCallback;
   app_on_mouse_event_callback_t mouseEventCallback;
   app_on_mouse_move_callback_t mouseMoveCallback;

   int64_t nextBrush;
   int64_t nextFont;
   std::vector<internal_brush_t> brushes;
   std::vector<InternalFont> fonts;

   app_t()
   {
      backgroundColor = app_color_rgb8(255, 255, 255);
      width = 0;
      height = 0;
      painting = false;
      externalData = nullptr;
      hwnd = nullptr;
      factory = nullptr;
      writeFactory = nullptr;
      renderTarget = nullptr;

      paintCallback = nullptr;

      // first null brush
      nextBrush = -1;
      nextFont = -1;
      brushes.push_back(internal_brush_t(APP_TRANSPARENT_BRUSH, app_color_rgba(1.0f, 1.0f, 1.0f, 0.0f)));
      brushes.push_back(internal_brush_t(APP_WHITE_BRUSH, app_color_rgba(1.0f, 1.0f, 1.0f, 1.0f)));
      brushes.push_back(internal_brush_t(APP_BLACK_BRUSH, app_color_rgba(0.0f, 0.0f, 0.0f, 1.0f)));
   }
};


static void create_brush(app_t * app, uint32_t idx)
{
   internal_brush_t & brush = app->brushes[idx];
   if (brush.brush == nullptr)
   {
      //printf("Creating app_brush_t(%d) with color (%f, %f, %f)\n", idx, brush.color.r, brush.color.g, brush.color.b);
      HRESULT result = app->renderTarget->CreateSolidColorBrush(D2D1::ColorF(brush.color.r, brush.color.g, brush.color.b, brush.color.a), &brush.brush);
      if (FAILED(result))
      {
         fprintf(stderr, "Failed to create brush %d\n", idx);
         brush.brush = nullptr;
      }
   }
}

static void create_font(app_t * app, uint32_t idx)
{
   InternalFont & font = app->fonts[idx];
   if (font.format == nullptr)
   {
      //printf("Creating font(%d) %s with size %f\n", idx, font.fontFamily.c_str(), font.pointSize);

      std::wstring fontFamily = to_wstr(font.fontFamily.c_str());
      HRESULT result = app->writeFactory->CreateTextFormat(fontFamily.c_str(),
                                                           nullptr, // font collection, use default one
                                                           to_dwrite_font_weight(font.fontWeight),
                                                           DWRITE_FONT_STYLE_NORMAL,
                                                           DWRITE_FONT_STRETCH_NORMAL,
                                                           to_dip_size(font.pointSize),
                                                           L"en-us",
                                                           &font.format);
      if (FAILED(result))
      {
         fprintf(stderr, "Failed to create font(%d) %s\n", idx, font.fontFamily.c_str());
         font.format = nullptr;
      }
   }
}

static HRESULT CreateDeviceIndependentResources(app_t * app)
{
   HRESULT result = S_OK;
   result = D2D1CreateFactory(D2D1_FACTORY_TYPE_SINGLE_THREADED, &app->factory);

   if (SUCCEEDED(result))
      result = DWriteCreateFactory(DWRITE_FACTORY_TYPE_SHARED, __uuidof(IDWriteFactory), reinterpret_cast<IUnknown**>(&app->writeFactory));
   return result;
}

static HRESULT CreateDeviceResources(app_t * app)
{
   HRESULT result = S_OK;

   if (!app->renderTarget)
   {
      RECT rect;
      GetClientRect(app->hwnd, &rect);

      D2D1_SIZE_U size = D2D1::SizeU(rect.right - rect.left, rect.bottom - rect.top);

      result = app->factory->CreateHwndRenderTarget(D2D1::RenderTargetProperties(),
                                                    D2D1::HwndRenderTargetProperties(app->hwnd, size),
                                                    &app->renderTarget);
      if (SUCCEEDED(result))
      {
         // Recreate brushes here
         for (uint32_t i = 0; i < app->brushes.size(); ++i)
            create_brush(app, i);

         // Recreate fonts here
         for (uint32_t i = 0; i < app->fonts.size(); ++i)
            create_font(app, i);
      }
   }

   return result;
}

static void AppOnResize(app_t * app, uint32_t width, uint32_t height)
{
   app->width = width;
   app->height = height;

   if (app->renderTarget)
   {
      HRESULT result = app->renderTarget->Resize(D2D1::SizeU(width, height));
      if (!SUCCEEDED(result))
         fprintf(stderr, "Could not resize window!\n");
   }

   if (app->resizeCallback)
   {
      app->resizeCallback(app, app->width, app->height);
   }
}

static void paint_app(app_t * app)
{
   HRESULT result = CreateDeviceResources(app);
   if (SUCCEEDED(result) && app->renderTarget)
   {
      app->painting = true;
      app->renderTarget->BeginDraw();
      app->renderTarget->SetTransform(D2D1::Matrix3x2F::Identity());
      app->renderTarget->Clear(D2D1::ColorF(app->backgroundColor.r, app->backgroundColor.g, app->backgroundColor.b, app->backgroundColor.a));

      if (app->paintCallback)
      {
         app->paintCallback(app, app->width, app->height);
      }

      app->renderTarget->EndDraw();

      app->painting = false;
   }
}

static LRESULT AppWndProc(HWND hwnd, UINT message, WPARAM wParam, LPARAM lParam)
{
   LRESULT result = 0;

   if (message == WM_CREATE)
   {
      LPCREATESTRUCT create = (LPCREATESTRUCT)lParam;
      app_t * app = static_cast<app_t *>(create->lpCreateParams);
      SetWindowLongPtr(hwnd, GWLP_USERDATA, PtrToUlong(app));
      result = 1;
   }
   else
   {
      app_t * app = reinterpret_cast<app_t *>(static_cast<LONG_PTR>(GetWindowLongPtrW(hwnd, GWLP_USERDATA)));

      bool wasHandled = false;

      if (app)
      {
         switch (message)
         {
            case WM_SIZE:
               {
                  uint32_t width = LOWORD(lParam);
                  uint32_t height = HIWORD(lParam);
                  AppOnResize(app, width, height);
                  result = 0;
                  wasHandled = true;
               }
               break;

            case WM_DISPLAYCHANGE:
               {
                  InvalidateRect(hwnd, nullptr, FALSE);
                  result = 0;
                  wasHandled = true;
               }
               break;

            case WM_PAINT:
               {
                  paint_app(app);
                  ValidateRect(hwnd, nullptr);
                  result = 0;
                  wasHandled = true;
               }
               break;

            case WM_DESTROY:
               {
                  PostQuitMessage(0);
                  result = 1;
                  wasHandled = true;
               }
               break;

            case WM_LBUTTONDOWN:
            case WM_RBUTTONDOWN:
            case WM_MBUTTONDOWN:
               {
                  int x = GET_X_LPARAM(lParam);
                  int y = GET_Y_LPARAM(lParam);
                  int button = 0;

                  switch (message)
                  {
                     case WM_LBUTTONDOWN:
                        button = APP_BUTTON_LEFT;
                        break;

                     case WM_RBUTTONDOWN:
                        button = APP_BUTTON_RIGHT;
                        break;

                     case WM_MBUTTONDOWN:
                        button = APP_BUTTON_MIDDLE;
                        break;
                  }

                  if (app->mouseEventCallback)
                     app->mouseEventCallback(app, button, APP_PRESS, x, y);
               }
               break;

            case WM_LBUTTONUP:
            case WM_RBUTTONUP:
            case WM_MBUTTONUP:
               {
                  int x = GET_X_LPARAM(lParam);
                  int y = GET_Y_LPARAM(lParam);
                  int button = 0;

                  switch (message)
                  {
                     case WM_LBUTTONUP:
                        button = APP_BUTTON_LEFT;
                        break;

                     case WM_RBUTTONUP:
                        button = APP_BUTTON_RIGHT;
                        break;

                     case WM_MBUTTONUP:
                        button = APP_BUTTON_MIDDLE;
                        break;
                  }

                  if (app->mouseEventCallback)
                     app->mouseEventCallback(app, button, APP_RELEASE, x, y);
               }
               break;

            case WM_MOUSEMOVE:
               {
                  int x = LOWORD(lParam);
                  int y = HIWORD(lParam);

                  if (app->mouseMoveCallback)
                     app->mouseMoveCallback(app, x, y);
               }
               break;
               
         }
      }

      if (!wasHandled)
         result = DefWindowProcW(hwnd, message, wParam, lParam);
   }

   return result;
}


EXTERN_C  app_t * app_init(const char * title, int width, int height)
{
   const wchar_t * CLASS_NAME = L"AppMainClass";

   if (FAILED(CoInitialize(nullptr)))
      return nullptr;

   app_t * app = new app_t();


   HRESULT result = CreateDeviceIndependentResources(app);

   if (SUCCEEDED(result))
   {
      WNDCLASSW wndClass = { sizeof(WNDCLASSW) };
      wndClass.style = CS_HREDRAW | CS_VREDRAW;
      wndClass.lpfnWndProc = AppWndProc;
      wndClass.cbClsExtra = 0;
      wndClass.cbWndExtra = sizeof(LONG_PTR);
      wndClass.hInstance = HINST_THISCOMPONENT;
      wndClass.hbrBackground = nullptr;
      wndClass.lpszMenuName = nullptr;
      wndClass.hCursor = LoadCursor(nullptr, IDI_APPLICATION);
      wndClass.lpszClassName = CLASS_NAME;

      RegisterClassW(&wndClass);

      float dpiX, dpiY;
      app->factory->GetDesktopDpi(&dpiX, &dpiY);


      std::wstring wTitle = to_wstr(title);

      app->hwnd = CreateWindowExW(0,
                                  CLASS_NAME,
                                  wTitle.c_str(),
                                  WS_OVERLAPPEDWINDOW,
                                  CW_USEDEFAULT,
                                  CW_USEDEFAULT,
                                  static_cast<uint32_t>(std::ceil(width * dpiX / 96.0f)),
                                  static_cast<uint32_t>(std::ceil(height * dpiY / 96.0f)),
                                  nullptr,
                                  nullptr,
                                  HINST_THISCOMPONENT,
                                  app);

      result = app->hwnd ? S_OK : E_FAIL;

      if (SUCCEEDED(result))
      {
         ShowWindow(app->hwnd, SW_SHOWNORMAL);
         UpdateWindow(app->hwnd);
         InvalidateRect(app->hwnd, nullptr, FALSE);
      }
   }

   if (FAILED(result))
   {
      delete app;
      return nullptr;
   }

   return app;
}

EXTERN_C  void app_shutdown(app_t * app)
{
   delete app;
   CoUninitialize();
}

EXTERN_C  void app_run(app_t * app)
{
   MSG msg;

   while (GetMessageW(&msg, nullptr, 0, 0))
   {
      TranslateMessage(&msg);
      DispatchMessageW(&msg);
   }
}

EXTERN_C void app_repaint(app_t* app)
{
   if (!app->painting)
      paint_app(app);
   //InvalidateRect(app->hwnd, nullptr, FALSE);
}

EXTERN_C void app_set_on_paint_callback(app_t * app, app_on_paint_callback_t callback)
{
   if (app)
      app->paintCallback = callback;
}

EXTERN_C void app_set_on_resize_callback(app_t * app, app_on_resize_callback_t callback)
{
   if (app)
      app->resizeCallback = callback;
}

EXTERN_C void app_set_on_mouse_event(app_t * app, app_on_mouse_event_callback_t callback)
{
   if (app)
      app->mouseEventCallback = callback;
}

EXTERN_C void app_set_on_mouse_move(app_t * app, app_on_mouse_move_callback_t callback)
{
   if (app)
      app->mouseMoveCallback = callback;
}

EXTERN_C  void app_set_background(app_t * app, app_color_t color)
{
   if (!app)
      return;

   app->backgroundColor = color;
   if (app->hwnd)
      InvalidateRect(app->hwnd, nullptr, FALSE);
}

EXTERN_C void app_set_data(app_t * app, void * data)
{
   if (app)
      app->externalData = data;
}

EXTERN_C void * app_get_data(app_t * app)
{
   return app ? app->externalData : nullptr;
}

EXTERN_C  app_brush_t app_create_solid_brush(app_t * app, app_color_t color)
{
   uint32_t brush = 0;

   if (app->nextBrush == -1)
   {
      brush = app->brushes.size();
      app->brushes.push_back(internal_brush_t(brush, color));
      create_brush(app, brush);
   }
   else
   {

   }

   return brush;
}

EXTERN_C  void app_destroy_brush(app_t * app, app_brush_t brush)
{
   if (brush < 1 || brush >= app->brushes.size())
      return;
}

EXTERN_C app_font_t app_create_font(app_t * app, const char * fontFamily, float pointSize, uint32_t fontWeight)
{
   uint32_t font = 0;

   if (app->nextFont == -1)
   {
      font = app->fonts.size();
      app->fonts.push_back(InternalFont(font, fontFamily, pointSize, fontWeight));
      create_font(app, font);
   }

   return font;
}

EXTERN_C void app_destroy_font(app_t * app, app_font_t font)
{

}

EXTERN_C void app_draw_line(app_t * app, float startX, float startY, float endX, float endY, app_brush_t brush, float thickness)
{
   if (app == nullptr || app->renderTarget == nullptr || brush < 1 || brush >= app->brushes.size() || app->brushes[brush].brush == nullptr)
   {
      fprintf(stderr, "Invalid DrawLine setup\n");
      return;
   }

   app->renderTarget->DrawLine(D2D1::Point2F(startX, startY), D2D1::Point2F(endX, endY), app->brushes[brush].brush, thickness);
}

static D2D1_RECT_F toD2DRect(app_rect_t rect)
{
   return D2D1::RectF(rect.left, rect.top, rect.right, rect.bottom);
}

EXTERN_C void app_draw_rectangle(app_t * app, app_rect_t rect, app_brush_t brush, float strokeThickness)
{
   if (app == nullptr || app->renderTarget == nullptr || brush < 1 || brush >= app->brushes.size() || app->brushes[brush].brush == nullptr)
   {
      fprintf(stderr, "Invalid draw_rectangle setup\n");
      return;
   }

   app->renderTarget->DrawRectangle(toD2DRect(rect),
                                    app->brushes[brush].brush,
                                    strokeThickness);
}

EXTERN_C void app_draw_rounded_rectangle(app_t * app, app_rect_t rect, float radius, app_brush_t brush, float strokeThickness)
{
   if (app == nullptr || app->renderTarget == nullptr || brush < 1 || brush >= app->brushes.size() || app->brushes[brush].brush == nullptr)
   {
      fprintf(stderr, "Invalid draw_rounded_rectangle setup\n");
      return;
   }

   app->renderTarget->DrawRoundedRectangle(D2D1::RoundedRect(toD2DRect(rect), radius, radius),
                                           app->brushes[brush].brush,
                                           strokeThickness);
}

EXTERN_C void app_fill_rectangle(app_t * app, app_rect_t rect, app_brush_t brush)
{
   if (app == nullptr || app->renderTarget == nullptr || brush < 1 || brush >= app->brushes.size() || app->brushes[brush].brush == nullptr)
   {
      fprintf(stderr, "Invalid fill_rectangle setup\n");
      return;
   }

   app->renderTarget->FillRectangle(toD2DRect(rect),
                                    app->brushes[brush].brush);
}

EXTERN_C void app_fill_rounded_rectangle(app_t * app, app_rect_t rect, float radius, app_brush_t brush)
{
   if (app == nullptr || app->renderTarget == nullptr || brush < 1 || brush >= app->brushes.size() || app->brushes[brush].brush == nullptr)
   {
      fprintf(stderr, "Invalid fill_rounded_rectangle setup\n");
      return;
   }

   app->renderTarget->FillRoundedRectangle(D2D1::RoundedRect(toD2DRect(rect), radius, radius),
                                           app->brushes[brush].brush);
}

EXTERN_C void app_draw_text(app_t * app, const char * text, app_font_t font, app_brush_t brush, app_rect_t bounds, uint32_t alignment)
{
   if (app == nullptr || app->renderTarget == nullptr ||
       font >= app->fonts.size() || app->fonts[font].format == nullptr ||
       brush >= app->brushes.size() || app->brushes[brush].brush == nullptr)
   {
      fprintf(stderr, "Invalid draw_text setup\n");
      return;
   }

   DWRITE_TEXT_ALIGNMENT textAlignment = DWRITE_TEXT_ALIGNMENT_CENTER;
   DWRITE_PARAGRAPH_ALIGNMENT paragraphAlignment = DWRITE_PARAGRAPH_ALIGNMENT_CENTER;

   if (alignment & APP_ALIGN_LEFT)
      textAlignment = DWRITE_TEXT_ALIGNMENT_LEADING;
   if (alignment & APP_ALIGN_RIGHT)
      textAlignment = DWRITE_TEXT_ALIGNMENT_TRAILING;
   if (alignment & APP_ALIGN_TOP)
      paragraphAlignment = DWRITE_PARAGRAPH_ALIGNMENT_NEAR;
   if (alignment & APP_ALIGN_BOTTOM)
      paragraphAlignment = DWRITE_PARAGRAPH_ALIGNMENT_FAR;

   app->fonts[font].format->SetTextAlignment(textAlignment);
   app->fonts[font].format->SetParagraphAlignment(paragraphAlignment);

   std::wstring str = to_wstr(text);
   app->renderTarget->DrawText(str.c_str(), str.size(),
                               app->fonts[font].format,
                               toD2DRect(bounds),
                               app->brushes[brush].brush);
}