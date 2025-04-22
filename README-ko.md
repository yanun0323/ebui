> [!IMPORTANT]
> 이 패키지는 현재 알파 개발 단계입니다

# EBUI

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)
[![日本語](https://img.shields.io/badge/日本語-クリック-青)](README-ja.md)
[![한국어](https://img.shields.io/badge/한국어-클릭-yellow)](README-ko.md)

EBUI는 [SwiftUI](https://developer.apple.com/documentation/swiftui)에서 영감을 받아 [Ebitengine](https://github.com/hajimehoshi/ebiten) 프레임워크 위에 구축된 Go 언어용 선언적 UI 프레임워크입니다. 개발자가 깔끔하고 기능적인 구문으로 대화형 애플리케이션을 만들 수 있게 해줍니다.

## 특징

- **선언적 구문**: SwiftUI와 유사한 깔끔한 선언적 구문으로 UI 구축
- **데이터 바인딩**: UI를 자동으로 업데이트하는 반응형 프로그래밍 모델
- **컴포넌트 기반 아키텍처**: 재사용 가능한 UI 컴포넌트 생성
- **레이아웃 시스템**: 유연한 스택 기반 레이아웃 시스템(VStack, HStack, ZStack)
- **수식어**: 스타일과 동작을 사용자 정의하기 위한 체인 가능한 수식어
- **실시간 미리보기**: VSCode 통합을 통한 실시간 변경 미리보기
- **크로스 플랫폼**: Ebitengine이 지원하는 모든 플랫폼에서 실행
- **통합 애니메이션**: 부드러운 UI 애니메이션을 위한 내장 지원
- **Ebitengine 통합**: 독립 실행형 앱으로 또는 기존 Ebitengine 프로젝트 내에서 작동

## 설치

```
# 패키지 관리자 지원 예정
```

## 빠른 시작

### 콘텐츠 뷰 정의

```go
import (
	ui "github.com/yanun0323/ebui"
	"image/color"
)

type contentView struct {
	title   *ui.Binding[string]
	content *ui.Binding[string]
}

func ContentView(title, content *ui.Binding[string]) ui.View {
	return &contentView{
		title:   title,
		content: content,
	}
}

func (view *contentView) Body() ui.SomeView {
	return ui.HStack(
		ui.Spacer(),
		ui.VStack(
			ui.Spacer(),
			ui.Text(view.title).
				Padding(ui.Bind(ui.Inset{0, 15, 0, 15})).
				ForegroundColor(ui.Bind[color.Color](color.White)).
				BackgroundColor(ui.Bind[color.Color](color.Gray{128})),
			ui.Text(view.content),
			ui.Spacer(),
		).Frame(ui.Bind(200.0), nil),
		ui.Spacer(),
	).
		ForegroundColor(ui.Bind[color.Color](color.RGBA{200, 200, 200, 255})).
		BackgroundColor(ui.Bind[color.Color](color.RGBA{255, 0, 0, 255})).
		Padding(ui.Bind(ui.Inset{5, 5, 5, 5}))
}
```

### 독립 실행형 앱으로 실행

```go
import (
	ui "github.com/yanun0323/ebui"
	"image/color"
	"log"
)

func main() {
	title := ui.Bind("title")
	content := ui.Bind("content")
	contentView := ContentView(title, content)

	app := ui.NewApplication(contentView)
	app.SetWindowBackgroundColor(color.RGBA{100, 0, 0, 0})
	app.SetWindowTitle("EBUI Demo")
	app.SetWindowSize(600, 500)
	app.SetWindowResizingMode(ui.WindowResizingModeEnabled)
	app.SetResourceFolder("resource")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

### Ebitengine 내에서 실행

```go
import (
	ui "github.com/yanun0323/ebui"
	"github.com/hajimehoshi/ebiten/v2"
	"log/slog"
	"image/color"
)

func main() {
	contentView := ui.ContentView("title", "content")
	g := NewGame(contentView)

	if err := ebiten.RunGame(g); err != nil {
		slog.Error("run game", "error", err)
	}
}

func NewGame(contentView ui.View) *Game {
	return &Game{
		contentView: contentView.Body(),
	}
}

type Game struct {
	contentView ui.SomeView
}

func (g *Game) Update() error {
	ui.EbitenUpdate(g.contentView)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ui.EbitenDraw(screen, g.contentView)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	ui.EbitenLayout(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}
```

## 사용 가능한 컴포넌트

### 뷰 (✓ = 구현됨)

| 기본 컴포넌트             | 레이아웃 컴포넌트   | 입력 컴포넌트               | 고급 컴포넌트            |
| ------------------------- | ------------------- | --------------------------- | ------------------------ |
| ✓ Text(텍스트)            | ✓ VStack(수직 스택) | ✓ Button(버튼)              | ✓ ScrollView(스크롤 뷰)  |
| ✓ Image(이미지)           | ✓ HStack(수평 스택) | ✓ Toggle(토글)              | ❑ List(목록)             |
| ✓ Rectangle(사각형)       | ✓ ZStack(중첩 스택) | ✓ Slider(슬라이더)          | ❑ TableView(테이블 뷰)   |
| ✓ Circle(원)              | ✓ Spacer(스페이서)  | ❑ TextField(텍스트 필드)    | ❑ Navigation(내비게이션) |
| ✓ Divider(구분선)         | ✓ EmptyView(빈 뷰)  | ❑ TextEditor(텍스트 에디터) | ❑ Sheet(시트)            |
| ✓ ViewModifier(뷰 수식어) |                     | ❑ Picker(선택기)            | ❑ Menu(메뉴)             |
|                           |                     | ❑ Radio(라디오)             | ❑ ProgressView(진행 뷰)  |
|                           |                     | ❑ DatePicker(날짜 선택기)   |                          |
|                           |                     | ❑ TimePicker(시간 선택기)   |                          |
|                           |                     | ❑ ColorPicker(색상 선택기)  |                          |

### 기능 (✓ = 구현됨)

- ✓ Modifier Stack(수식어 스택)
- ✓ CornerRadius(모서리 반경)
- ✓ Animation(애니메이션)
- ✓ Alignment(정렬)
- ❑ Gesture(제스처)
- ❑ Overlay(오버레이)
- ❑ Mask(마스크)
- ❑ Clip(클립)

## VSCode 실시간 미리보기

EBUI는 SwiftUI와 유사한 개발 경험을 제공하여 VSCode에서 직접 UI의 실시간 미리보기를 할 수 있습니다.

### VSCode 확장 기능 사용

[ebui-vscode 확장 기능](https://github.com/yanun0323/ebui-vscode)은 `Preview_`로 시작하는 함수에 대한 핫 리로드 기능을 제공합니다.

1. VSCode 마켓플레이스에서 확장 기능 설치
2. Go 파일에 `Preview_` 접두사가 있는 함수 생성
3. 저장하면 미리보기 창에서 UI 업데이트를 실시간으로 확인 가능

### 예시

```go
package mypackage

import (
	ui "github.com/yanun0323/ebui"
	"image/color"
)

// 이 함수는 자동으로 미리보기됩니다
func Preview_MyButton() ui.View {
	return ui.Button(ui.Const("클릭하세요")).
		BackgroundColor(ui.Bind[color.Color](color.RGBA{200, 100, 100, 255})).
		Padding(ui.Bind(ui.Inset{10, 10, 10, 10})).
		Center()
}
```

### 작동 방식

확장 기능의 워크플로우:

1. Go 파일의 변경 사항 모니터링
2. 저장 시 EBUI 미리보기 도구 실행
3. 코드에서 `Preview_` 함수 감지
4. 실시간 미리보기 창에서 렌더링

이를 통해 편집기를 떠나지 않고도 UI 개발을 위한 빠른 피드백 루프를 제공합니다.
