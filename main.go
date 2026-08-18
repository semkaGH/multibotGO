package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Страницы приложения
type Page int

const (
	PageHome Page = iota
	PageCalculator
	PageMinecraftAnalyzer
	PageYouTubeDownloader
	PageSettings
	PageAbout
)

// Режимы анализа Minecraft
type MCMode int

const (
	MCModeNormal MCMode = iota
	MCModeAI
)

// Состояние приложения
type App struct {
	currentPage       Page
	mcMode            MCMode
	theme             *material.Theme
	downloading       bool
	downloadProgress  float32
	autoUpdateEnabled bool

	// Виджеты навигации
	homeBtn, calcBtn, mcBtn, ytBtn, settingsBtn, aboutBtn widget.Clickable

	// Калькулятор
	calcInput  widget.Editor
	calcResult string
	calcBtn    widget.Clickable

	// Minecraft анализатор
	mcLogInput     widget.Editor
	mcResult       string
	mcAnalyzeBtn   widget.Clickable
	mcAIPrompt     widget.Editor
	mcAIAnalyzeBtn widget.Clickable

	// YouTube загрузчик
	ytURLInput    widget.Editor
	ytQuality     int    // 1080, 720, 480, 360
	ytFormat      string // "mp4" или "mp3"
	ytDownloadBtn widget.Clickable

	// Настройки
	autoUpdateSwitch widget.Bool

	// Текстовые виджеты для скролла
	minecraftScroll, ytScroll layout.List
}

func NewApp() *App {
	return &App{
		currentPage:       PageHome,
		mcMode:            MCModeNormal,
		theme:             material.NewTheme(),
		autoUpdateEnabled: true,
		ytQuality:         1080,
		ytFormat:          "mp4",
	}
}

func (a *App) Layout(gtx layout.Context) layout.Dimensions {
	// Обновляем тему
	a.theme.Shaper = text.NewShaper(text.WithCollection(a.theme.Shaper.Collection()))

	switch a.currentPage {
	case PageHome:
		return a.layoutHome(gtx)
	case PageCalculator:
		return a.layoutCalculator(gtx)
	case PageMinecraftAnalyzer:
		return a.layoutMinecraftAnalyzer(gtx)
	case PageYouTubeDownloader:
		return a.layoutYouTubeDownloader(gtx)
	case PageSettings:
		return a.layoutSettings(gtx)
	case PageAbout:
		return a.layoutAbout(gtx)
	default:
		return a.layoutHome(gtx)
	}
}

func (a *App) layoutHome(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(50), Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					lbl := material.H4(a.theme, "MultiTool GO")
					lbl.Alignment = text.Middle
					return lbl.Layout(gtx)
				})
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.theme, &a.homeBtn, "🏠 Главная")
				btn.Background = color.NRGBA{R: 63, G: 81, B: 181, A: 255}
				return btn.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.theme, &a.calcBtn, "🧮 Калькулятор")
				return btn.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.theme, &a.mcBtn, "⛏️ Анализ Minecraft")
				return btn.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.theme, &a.ytBtn, "📺 YouTube Загрузчик")
				return btn.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.theme, &a.settingsBtn, "⚙️ Настройки")
				return btn.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.theme, &a.aboutBtn, "ℹ️ О приложении")
				return btn.Layout(gtx)
			})
		}),
	)
}

func (a *App) layoutCalculator(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(50), Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					lbl := material.H5(a.theme, "Калькулятор")
					return lbl.Layout(gtx)
				})
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				a.calcInput.SingleLine = true
				editor := material.Editor(a.theme, &a.calcInput, "Введите выражение (например: 2+2*2)")
				return editor.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.theme, &a.calcBtn, "Вычислить")
				return btn.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				result := material.Body1(a.theme, "Результат: "+a.calcResult)
				return result.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.theme, &a.homeBtn, "← Назад")
				return btn.Layout(gtx)
			})
		}),
	)
}

func (a *App) layoutMinecraftAnalyzer(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(50), Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					lbl := material.H5(a.theme, "Анализ ошибок Minecraft")
					return lbl.Layout(gtx)
				})
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEvenly}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						radio := material.RadioButton(a.theme, &a.mcMode, MCModeNormal, "Обычный")
						return radio.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						radio := material.RadioButton(a.theme, &a.mcMode, MCModeAI, "ИИ (Gemini/OpenRouter)")
						return radio.Layout(gtx)
					}),
				)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if a.mcMode == MCModeNormal {
				return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					editor := material.Editor(a.theme, &a.mcLogInput, "Вставьте лог ошибки Minecraft")
					editor.SetText(a.mcLogInput.Text())
					return editor.Layout(gtx)
				})
			} else {
				return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					editor := material.Editor(a.theme, &a.mcAIPrompt, "Опишите проблему для ИИ анализа")
					editor.SetText(a.mcAIPrompt.Text())
					return editor.Layout(gtx)
				})
			}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				var btn widget.Clickable
				var btnText string
				if a.mcMode == MCModeNormal {
					btn = a.mcAnalyzeBtn
					btnText = "Анализировать лог"
				} else {
					btn = a.mcAIAnalyzeBtn
					btnText = "Анализировать с ИИ"
				}
				button := material.Button(a.theme, &btn, btnText)
				return button.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				result := material.Body1(a.theme, "Результат:\n"+a.mcResult)
				return result.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.theme, &a.homeBtn, "← Назад")
				return btn.Layout(gtx)
			})
		}),
	)
}

func (a *App) layoutYouTubeDownloader(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(50), Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					lbl := material.H5(a.theme, "YouTube Загрузчик")
					return lbl.Layout(gtx)
				})
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				a.ytURLInput.SingleLine = true
				editor := material.Editor(a.theme, &a.ytURLInput, "Вставьте ссылку на YouTube видео")
				return editor.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Body1(a.theme, "Формат: ")
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								radio := material.RadioButton(a.theme, &a.ytFormat, "mp4", "MP4")
								return radio.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								radio := material.RadioButton(a.theme, &a.ytFormat, "mp3", "MP3")
								return radio.Layout(gtx)
							}),
						)
					}),
				)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEvenly}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						radio := material.RadioButton(a.theme, &a.ytQuality, 1080, "1080p")
						return radio.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						radio := material.RadioButton(a.theme, &a.ytQuality, 720, "720p")
						return radio.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						radio := material.RadioButton(a.theme, &a.ytQuality, 480, "480p")
						return radio.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						radio := material.RadioButton(a.theme, &a.ytQuality, 360, "360p")
						return radio.Layout(gtx)
					}),
				)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.theme, &a.ytDownloadBtn, "Скачать")
				if a.downloading {
					btn.Text = "Загрузка..."
				}
				return btn.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if a.downloading {
				return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					progress := material.ProgressBar(a.theme, a.downloadProgress)
					return progress.Layout(gtx)
				})
			}
			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.theme, &a.homeBtn, "← Назад")
				return btn.Layout(gtx)
			})
		}),
	)
}

func (a *App) layoutSettings(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(50), Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					lbl := material.H5(a.theme, "Настройки")
					return lbl.Layout(gtx)
				})
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Body1(a.theme, "Автообновление")
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						sw := material.Switch(a.theme, &a.autoUpdateSwitch, "Включить автообновление")
						return sw.Layout(gtx)
					}),
				)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.theme, &a.homeBtn, "← Назад")
				return btn.Layout(gtx)
			})
		}),
	)
}

func (a *App) layoutAbout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(50), Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					lbl := material.H5(a.theme, "О приложении")
					return lbl.Layout(gtx)
				})
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				text := material.Body1(a.theme, "MultiTool GO v1.0.0")
				return text.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				text := material.Body1(a.theme, "Приложение создано на Go с использованием Gio UI")
				return text.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				text := material.Body1(a.theme, "Функции:")
				return text.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(40), Right: unit.Dp(20), Bottom: unit.Dp(5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				text := material.Body1(a.theme, "• Калькулятор")
				return text.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(40), Right: unit.Dp(20), Bottom: unit.Dp(5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				text := material.Body1(a.theme, "• Анализ ошибок Minecraft (обычный и ИИ)")
				return text.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(40), Right: unit.Dp(20), Bottom: unit.Dp(5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				text := material.Body1(a.theme, "• Загрузка видео с YouTube (MP3/MP4)")
				return text.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(40), Right: unit.Dp(20), Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				text := material.Body1(a.theme, "• Автообновление")
				return text.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.theme, &a.homeBtn, "← Назад")
				return btn.Layout(gtx)
			})
		}),
	)
}

func (a *App) handleEvents(gtx layout.Context) {
	// Обработка навигации
	if a.homeBtn.Clicked(gtx) {
		a.currentPage = PageHome
	}
	if a.calcBtn.Clicked(gtx) {
		a.currentPage = PageCalculator
	}
	if a.mcBtn.Clicked(gtx) {
		a.currentPage = PageMinecraftAnalyzer
	}
	if a.ytBtn.Clicked(gtx) {
		a.currentPage = PageYouTubeDownloader
	}
	if a.settingsBtn.Clicked(gtx) {
		a.currentPage = PageSettings
	}
	if a.aboutBtn.Clicked(gtx) {
		a.currentPage = PageAbout
	}

	// Обработка калькулятора
	if a.calcBtn.Clicked(gtx) {
		expr := a.calcInput.Text()
		result, err := evaluateExpression(expr)
		if err != nil {
			a.calcResult = "Ошибка: " + err.Error()
		} else {
			a.calcResult = result
		}
	}

	// Обработка Minecraft анализатора
	if a.mcAnalyzeBtn.Clicked(gtx) {
		logText := a.mcLogInput.Text()
		a.mcResult = analyzeMinecraftLog(logText)
	}

	if a.mcAIAnalyzeBtn.Clicked(gtx) {
		prompt := a.mcAIPrompt.Text()
		a.mcResult = "ИИ анализ (требуется API ключ):\n" + prompt + "\n\n[Здесь будет интеграция с Gemini/OpenRouter]"
	}

	// Обработка YouTube загрузчика
	if a.ytDownloadBtn.Clicked(gtx) && !a.downloading {
		url := a.ytURLInput.Text()
		if url != "" {
			a.downloading = true
			a.downloadProgress = 0
			// Здесь должна быть логика загрузки
			go func() {
				// Симуляция загрузки
				for i := 0; i <= 100; i += 10 {
					a.downloadProgress = float32(i) / 100.0
				}
				a.downloading = false
			}()
		}
	}

	// Обработка настроек
	a.autoUpdateEnabled = a.autoUpdateSwitch.Value
}

// Простой парсер выражений для калькулятора
func evaluateExpression(expr string) (string, error) {
	// Упрощенная реализация - в продакшене используйте proper math parser
	expr = strings.ReplaceAll(expr, " ", "")

	// Очень простая реализация только для базовых операций
	// Для полноценного калькулятора используйте github.com/Knetic/govaluate или аналогичную библиотеку

	// Пример простой реализации
	var result float64
	_, err := fmt.Sscanf(expr, "%f", &result)
	if err == nil && len(expr) < 20 {
		return strconv.FormatFloat(result, 'f', -1, 64), nil
	}

	// TODO: Реализовать полноценный парсер математических выражений
	return "Требуется реализация парсера", nil
}

func analyzeMinecraftLog(logText string) string {
	if logText == "" {
		return "Введите лог для анализа"
	}

	// Простой анализ распространенных ошибок
	errors := []string{
		"OutOfMemoryError",
		"NullPointerException",
		"ClassNotFoundException",
		"NoClassDefFoundError",
		"StackOverflowError",
	}

	foundErrors := []string{}
	for _, err := range errors {
		if strings.Contains(logText, err) {
			foundErrors = append(foundErrors, err)
		}
	}

	if len(foundErrors) == 0 {
		return "Явных ошибок не найдено. Проверьте логи внимательнее."
	}

	result := "Найдены ошибки:\n"
	for _, err := range foundErrors {
		switch err {
		case "OutOfMemoryError":
			result += "- OutOfMemoryError: Увеличьте выделенную память для Minecraft\n"
		case "NullPointerException":
			result += "- NullPointerException: Проблема с модом или конфигурацией\n"
		case "ClassNotFoundException":
			result += "- ClassNotFoundException: Отсутствует необходимый класс, проверьте моды\n"
		case "NoClassDefFoundError":
			result += "- NoClassDefFoundError: Конфликт версий модов\n"
		case "StackOverflowError":
			result += "- StackOverflowError: Бесконечная рекурсия в моде\n"
		}
	}

	return result
}

func main() {
	go func() {
		w := app.NewWindow()
		w.Title = "MultiTool GO"
		w.Option(app.Size(unit.Dp(400), unit.Dp(700)))

		appState := NewApp()

		var ops op.Ops

		for e := range w.Events() {
			switch e := e.(type) {
			case system.DestroyEvent:
				os.Exit(0)
			case system.FrameEvent:
				gtx := app.NewContext(&ops, e)

				appState.handleEvents(gtx)

				appState.theme.Palette.Bg = color.NRGBA{R: 245, G: 245, B: 245, A: 255}
				appState.theme.Palette.Fg = color.NRGBA{R: 0, G: 0, B: 0, A: 255}

				appState.Layout(gtx)
				e.Submit(&ops)
			}
		}
	}()

	app.Main()
}
