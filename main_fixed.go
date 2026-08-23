package main

import (
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeAI
)

type C2S struct {
	th *material.Theme

	// Навигация
	currentTab int
	tabs       []string

	// Minecraft
	mcMode      Mode
	mcModeRadio widget.Enum
	logOutput   widget.Editor
	startBtn    widget.Clickable
	clearBtn    widget.Clickable
	progress    float32
	analyzing   bool
	statusText  string

	// YouTube
	ytURL       widget.Editor
	ytFormat    widget.Enum
	downloadBtn widget.Clickable
	ytProgress  float32
	downloading bool

	// Калькулятор
	calcDisplay widget.Editor
	calcBtns    [19]widget.Clickable

	// Настройки
	apiKey     widget.Editor
	checkBtn   widget.Clickable
	themeRadio widget.Enum

	// Общее
	backBtn widget.Clickable
}

func main() {
	go func() {
		window := app.NewWindow(
			app.Title("MultiBot GO"),
			app.Size(unit.Dp(900), unit.Dp(700)),
		)

		if err := layout(window); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func layout(window *app.Window) error {
	state := &C2S{
		currentTab: 0,
		tabs:       []string{"Minecraft", "YouTube", "Calculator", "Settings", "About"},
	}

	state.th = material.NewTheme()
	state.th.Shaper = text.NewShaper(text.WithCollection(state.th.Shaper.Collection))

	// Инициализация редакторов
	state.logOutput.SingleLine = false
	state.logOutput.ReadOnly = true

	state.ytURL.SingleLine = true
	state.ytURL.Text = ""

	state.calcDisplay.SingleLine = true
	state.calcDisplay.ReadOnly = true
	state.calcDisplay.Text = "0"

	state.apiKey.SingleLine = true
	state.apiKey.Text = ""

	// Инициализация Radio
	state.mcModeRadio.Value = "normal"
	state.ytFormat.Value = "mp4_1080"
	state.themeRadio.Value = "system"

	var ops op.Ops
	for {
		e := <-window.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Обработка кликов навигации
			if state.backBtn.Clicked(gtx) {
				state.currentTab = 0
			}

			// Рендеринг в зависимости от вкладки
			var body layout.Widget
			switch state.currentTab {
			case 0:
				body = state.layoutMinecraftTab(gtx)
			case 1:
				body = state.layoutYouTubeTab(gtx)
			case 2:
				body = state.layoutCalculatorTab(gtx)
			case 3:
				body = state.layoutSettingsTab(gtx)
			case 4:
				body = state.layoutAboutTab(gtx)
			}

			// Основной макет с боковой панелью
			material.Shader(gtx, color.NRGBA{R: 30, G: 30, B: 30, A: 255})
			
			flex := layout.Flex{Axis: layout.Horizontal}
			flex.Layout(gtx,
				layout.Rigid(state.layoutSidebar(gtx)),
				layout.Flexed(1, layout.Inset{Top: unit.Dp(10), Right: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(10)}.Layout(body)),
			)

			e.Operation(&ops)
		}
	}
}

func (c *C2S) layoutSidebar(gtx layout.Context) layout.Dimensions {
	sidebar := material.List(c.th, &widget.List{})
	sidebar.Axis = layout.Vertical
	
	return material.Card{Color: color.NRGBA{R: 20, G: 20, B: 20, A: 255}}.Layout(gtx,
		layout.Inset{Top: unit.Dp(20), Bottom: unit.Dp(20)}.Layout(func(gtx layout.Context) layout.Dimensions {
			children := []layout.FlexChild{}
			for i, tab := range c.tabs {
				idx := i
				btn := widget.Clickable{} // В реальном приложении нужно хранить состояние кнопок
				label := material.Body1(c.th, tab)
				if i == c.currentTab {
					label.Color = color.NRGBA{R: 100, G: 200, B: 255, A: 255}
				}
				
				child := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btnWidget := material.Button(c.th, &btn, tab)
					btnWidget.Background = color.NRGBA{A: 0}
					if i == c.currentTab {
						btnWidget.Background = color.NRGBA{R: 50, G: 50, B: 50, A: 255}
					}
					
					if btn.Clicked(gtx) {
						c.currentTab = idx
					}
					
					return btnWidget.Layout(gtx)
				})
				children = append(children, child)
			}
			return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceStart}.Layout(gtx, children...)
		}),
	)
}

func (c *C2S) layoutMinecraftTab(gtx layout.Context) layout.Dimensions {
	// Обновление состояния радио-кнопок
	if c.mcModeRadio.Update(gtx) {
		if c.mcModeRadio.Value == "normal" {
			c.mcMode = ModeNormal
		} else {
			c.mcMode = ModeAI
		}
	}

	// Кнопка запуска
	if c.startBtn.Clicked(gtx) && !c.analyzing {
		c.analyzing = true
		c.progress = 0
		c.statusText = "Анализ логов..."
		c.logOutput.SetText("")
		
		go func() {
			steps := []string{
				"Чтение latest.log...",
				"Поиск паттернов ошибок...",
				"Проверка модов...",
			}
			for i, step := range steps {
				time.Sleep(600 * time.Millisecond)
				c.logOutput.SetText(c.logOutput.Text + "[INFO] " + step + "\n")
				c.progress = float32(i+1) / float32(len(steps)+1)
			}
			
			// Симуляция ошибки
			c.logOutput.SetText(c.logOutput.Text + "[ERROR] java.lang.OutOfMemoryError in Optifine\n[WARN] Conflict: FabricAPI vs Sodium\n")
			c.statusText = "Найдены критические ошибки!"
			c.progress = 1.0
			c.analyzing = false
		}()
	}

	// Кнопка очистки
	if c.clearBtn.Clicked(gtx) {
		c.logOutput.SetText("")
		c.statusText = "Ожидание"
		c.progress = 0
	}

	return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.H6(c.th, "Анализатор ошибок Minecraft").Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				normalRadio := material.RadioButton(c.th, &c.mcModeRadio, "normal", "Обычный режим")
				aiRadio := material.RadioButton(c.th, &c.mcModeRadio, "ai", "ИИ Режим (Gemini/OpenRouter)")
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(normalRadio.Layout),
					layout.Rigid(layout.Spacer{Width: unit.Dp(20)}.Layout),
					layout.Rigid(aiRadio.Layout),
				)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			startBtn := material.Button(c.th, &c.startBtn, "Запустить анализ")
			if c.analyzing {
				startBtn.Background = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
			}
			
			clearBtn := material.Button(c.th, &c.clearBtn, "Очистить")
			
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if c.analyzing {
						return layout.Dimensions{}
					}
					return startBtn.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Rigid(clearBtn.Layout),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				statusLabel := material.Body1(c.th, "Статус: "+c.statusText)
				return statusLabel.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				progressBar := material.ProgressBar(c.th, c.progress)
				return progressBar.Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			editor := material.Editor(c.th, &c.logOutput, "Лог событий...")
			editor.Color = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
			return editor.Layout(gtx)
		}),
	)
}

func (c *C2S) layoutYouTubeTab(gtx layout.Context) layout.Dimensions {
	if c.ytFormat.Update(gtx) {
		// Обработка выбора формата
	}

	if c.downloadBtn.Clicked(gtx) && !c.downloading {
		if c.ytURL.Text == "" {
			c.ytURL.SetText("Введите URL!")
			return layout.Dimensions{}
		}
		
		c.downloading = true
		c.ytProgress = 0
		
		go func() {
			for i := 0; i <= 10; i++ {
				time.Sleep(200 * time.Millisecond)
				c.ytProgress = float32(i) / 10.0
			}
			c.downloading = false
			c.ytURL.SetText("")
		}()
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.H6(c.th, "Загрузчик YouTube").Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			urlEditor := material.Editor(c.th, &c.ytURL, "Вставьте ссылку на YouTube...")
			return urlEditor.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				mp4_1080 := material.RadioButton(c.th, &c.ytFormat, "mp4_1080", "MP4 1080p")
				mp4_720 := material.RadioButton(c.th, &c.ytFormat, "mp4_720", "MP4 720p")
				mp3 := material.RadioButton(c.th, &c.ytFormat, "mp3", "MP3 Audio")
				
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(mp4_1080.Layout),
					layout.Rigid(mp4_720.Layout),
					layout.Rigid(mp3.Layout),
				)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(c.th, &c.downloadBtn, "Скачать")
			if c.downloading {
				btn.Background = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
			}
			return btn.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if c.downloading || c.ytProgress > 0 {
				return layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return material.ProgressBar(c.th, c.ytProgress).Layout(gtx)
				})
			}
			return layout.Dimensions{}
		}),
	)
}

func (c *C2S) layoutCalculatorTab(gtx layout.Context) layout.Dimensions {
	buttons := []string{
		"C", "(", ")", "/",
		"7", "8", "9", "*",
		"4", "5", "6", "-",
		"1", "2", "3", "+",
		"0", ".", "=", "⌫",
	}

	for i, btnText := range buttons {
		if c.calcBtns[i].Clicked(gtx) {
			switch btnText {
			case "C":
				c.calcDisplay.Text = "0"
			case "⌫":
				if len(c.calcDisplay.Text) > 1 {
					c.calcDisplay.Text = c.calcDisplay.Text[:len(c.calcDisplay.Text)-1]
				} else {
					c.calcDisplay.Text = "0"
				}
			case "=":
				// Простая эмуляция, в реальности нужен парсер
				c.calcDisplay.Text = "Result: " + c.calcDisplay.Text
			default:
				if c.calcDisplay.Text == "0" {
					c.calcDisplay.Text = btnText
				} else {
					c.calcDisplay.Text += btnText
				}
			}
		}
	}

	grid := layout.Grid{Cols: 4}
	children := []layout.GridChild{}
	for i := 0; i < 19; i++ {
		idx := i
		child := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(c.th, &c.calcBtns[idx], buttons[idx])
			return btn.Layout(gtx)
		})
		children = append(children, child)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.H6(c.th, "Калькулятор").Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			editor := material.Editor(c.th, &c.calcDisplay, "")
			editor.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			editor.Background = color.NRGBA{R: 50, G: 50, B: 50, A: 255}
			return editor.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return grid.Layout(gtx, children...)
		}),
	)
}

func (c *C2S) layoutSettingsTab(gtx layout.Context) layout.Dimensions {
	if c.checkBtn.Clicked(gtx) {
		go func() {
			time.Sleep(2 * time.Second)
			// Эмуляция проверки
		}()
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.H6(c.th, "Настройки").Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Body1(c.th, "API Ключ (Gemini/OpenRouter):")
				return label.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			editor := material.Editor(c.th, &c.apiKey, "Введите API ключ...")
			return editor.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(15), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Body1(c.th, "Тема оформления:")
				return label.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			sysTheme := material.RadioButton(c.th, &c.themeRadio, "system", "Системная")
			darkTheme := material.RadioButton(c.th, &c.themeRadio, "dark", "Темная")
			lightTheme := material.RadioButton(c.th, &c.themeRadio, "light", "Светлая")
			
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(sysTheme.Layout),
				layout.Rigid(darkTheme.Layout),
				layout.Rigid(lightTheme.Layout),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(c.th, &c.checkBtn, "Проверить обновления")
				return btn.Layout(gtx)
			})
		}),
	)
}

func (c *C2S) layoutAboutTab(gtx layout.Context) layout.Dimensions {
	content := `## MultiBot GO
Приложение создано на Go с использованием Gio UI.

### Возможности:
- Анализ ошибок Minecraft
- Загрузка YouTube (MP3/MP4)
- Калькулятор
- Настройки

Версия: 1.0.0
Разработчик: Semka`

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(50), Bottom: unit.Dp(50), Left: unit.Dp(50), Right: unit.Dp(50)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(c.th, content)
			return label.Layout(gtx)
		})
	})
}
