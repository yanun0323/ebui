package main

import (
	"fmt"
	"log"

	ui "github.com/yanun0323/ebui"
	ebitenbackend "github.com/yanun0323/ebui/backend/ebiten"
)

type playgroundScreen struct {
	count      *ui.State[int]
	doubleMode *ui.State[bool]
	progress   *ui.State[float64]
}

func newPlaygroundScreen() *playgroundScreen {
	return &playgroundScreen{
		count:      ui.NewState(0),
		doubleMode: ui.NewState(false),
		progress:   ui.NewState(35.0),
	}
}

func (s *playgroundScreen) Body(ctx ui.ViewContext) ui.Element {
	count := ui.Read(ctx, s.count)
	doubleMode := ui.Read(ctx, s.doubleMode)
	progress := ui.Read(ctx, s.progress)

	titleStyle := ui.DefaultTextStyle()
	titleStyle.Size = 26
	titleStyle.Weight = ui.WeightBold
	titleStyle.LineHeight = 6

	bodyStyle := ui.DefaultTextStyle()
	bodyStyle.Size = 16
	bodyStyle.LineHeight = 4

	panel := func(children ...ui.View) ui.Element {
		return ui.Apply(
			ui.VStack(children...),
			ui.Spacing(12),
			ui.Padding(ui.All(16)),
			ui.Background(ui.RGBA(245, 247, 250, 255)),
			ui.Border(ui.UniformStroke(1, ui.RGBA(214, 220, 232, 255))),
		)
	}

	step := 1
	modeText := "一般模式"
	if doubleMode {
		step = 2
		modeText = "雙倍模式"
	}

	progressItems := make([]ui.View, 0, 24)
	for i := 1; i <= 24; i++ {
		label := fmt.Sprintf("第 %02d 筆資料", i)
		status := "待處理"
		if float64(i*4) <= progress {
			status = "已完成"
		}
		progressItems = append(progressItems,
			ui.Apply(
				ui.HStack(
					ui.Apply(ui.Text(label), ui.Font(bodyStyle), ui.Foreground(ui.RGBA(32, 40, 56, 255))),
					ui.Spacer(),
					ui.Apply(ui.Text(status), ui.Font(bodyStyle), ui.Foreground(ui.RGBA(88, 109, 150, 255))),
				),
				ui.Frame(ui.S(288, 28)),
			),
		)
	}

	header := ui.Apply(
		ui.VStack(
			ui.Apply(ui.Text("EBUI v2 Playground"), ui.Font(titleStyle), ui.Foreground(ui.White)),
			ui.Apply(ui.Text("用 button、toggle、slider、scroll view 驗證目前 API。"), ui.Font(bodyStyle), ui.Foreground(ui.RGBA(222, 229, 244, 255))),
		),
		ui.Spacing(8),
		ui.Padding(ui.All(18)),
		ui.Background(ui.RGBA(46, 74, 143, 255)),
	)

	counterPanel := panel(
		ui.Apply(ui.Text(fmt.Sprintf("目前計數：%d", count)), ui.Font(bodyStyle), ui.Foreground(ui.RGBA(27, 33, 45, 255))),
		ui.Apply(ui.Text(fmt.Sprintf("加減步長：%d，%s", step, modeText)), ui.Font(bodyStyle), ui.Foreground(ui.RGBA(96, 109, 138, 255))),
		ui.Apply(
			ui.HStack(
				ui.Apply(
					ui.Button(
						ui.Apply(ui.Text("增加"), ui.Foreground(ui.White), ui.Font(bodyStyle)),
						func() {
							s.count.Update(func(v int) int { return v + step })
						},
					),
					ui.Padding(ui.Symmetric(10, 14)),
					ui.Background(ui.RGBA(58, 111, 219, 255)),
					ui.Frame(ui.S(120, 44)),
				),
				ui.Apply(
					ui.Button(
						ui.Apply(ui.Text("減少"), ui.Foreground(ui.RGBA(38, 51, 76, 255)), ui.Font(bodyStyle)),
						func() {
							s.count.Update(func(v int) int { return v - step })
						},
					),
					ui.Padding(ui.Symmetric(10, 14)),
					ui.Background(ui.RGBA(228, 234, 244, 255)),
					ui.Frame(ui.S(120, 44)),
				),
			),
			ui.Spacing(10),
		),
	)

	controlPanel := panel(
		ui.HStack(
			ui.Apply(
				ui.VStack(
					ui.Apply(ui.Text("雙倍模式"), ui.Font(bodyStyle), ui.Foreground(ui.RGBA(27, 33, 45, 255))),
					ui.Apply(ui.Text("切換後按鈕每次加減 2。"), ui.Font(bodyStyle), ui.Foreground(ui.RGBA(96, 109, 138, 255))),
				),
				ui.Spacing(6),
			),
			ui.Spacer(),
			ui.Apply(ui.Toggle(s.doubleMode), ui.Frame(ui.S(60, 32))),
		),
		ui.Apply(ui.Text(fmt.Sprintf("完成進度：%.0f%%", progress)), ui.Font(bodyStyle), ui.Foreground(ui.RGBA(27, 33, 45, 255))),
		ui.Apply(ui.Slider(s.progress, 0, 100), ui.Frame(ui.S(320, 32))),
	)

	listPanel := panel(
		ui.Apply(ui.Text("捲動清單"), ui.Font(bodyStyle), ui.Foreground(ui.RGBA(27, 33, 45, 255))),
		ui.Apply(
			ui.ScrollView(ui.Apply(ui.VStack(progressItems...), ui.Spacing(6))),
			ui.Frame(ui.S(320, 220)),
			ui.Background(ui.White),
		),
	)

	content := ui.Apply(
		ui.VStack(
			header,
			ui.Apply(
				ui.HStack(
					ui.Apply(ui.VStack(counterPanel, controlPanel), ui.Spacing(14)),
					listPanel,
				),
				ui.Spacing(14),
			),
		),
		ui.Spacing(14),
		ui.Padding(ui.All(20)),
		ui.Frame(ui.S(760, 540)),
	)

	return ui.Apply(
		ui.ZStack(
			ui.Apply(
				ui.Rectangle(),
				ui.Fill(ui.RGBA(234, 239, 247, 255)),
				ui.Frame(ui.S(760, 540)),
			),
			content,
		),
		ui.Frame(ui.S(760, 540)),
	)
}

func main() {
	app, err := ebitenbackend.NewApp(
		newPlaygroundScreen(),
		ebitenbackend.WithWindowTitle("EBUI v2 Playground"),
		ebitenbackend.WithWindowSize(760, 540),
		ebitenbackend.WithBackground(ui.RGBA(228, 234, 242, 255)),
		ebitenbackend.WithFloating(true),
		ebitenbackend.WithFontPath("resource/Taipei Sans TC Regular.ttf"),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
