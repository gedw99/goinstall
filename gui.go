package goinstall

import (
	"gioui.org/app"
	"image/color"
	"log"
	"os"
	"time"

	"gioui.org/font/gofont"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type status uint8

const (
	awaitAcceptCondition status = iota
	processing
	succeed
	failed
)

func (s *guiState) display(gtx layout.Context) layout.Dimensions {
	var widgets []layout.Widget

	widgets = append(widgets, s.title.Layout)

	switch s.status {
	case awaitAcceptCondition:
		widgets = append(widgets, s.conditions...)
		widgets = append(widgets, func(gtx C) D {
			return s.displayAcceptButton(gtx)
		})
	case processing:
		widgets = append(widgets, s.steps...)
	case succeed:
		widgets = append(widgets, s.succeed)
	case failed:
		widgets = append(widgets, s.failed)
	}
	return displayWidgetList(gtx, s.list, widgets)
}

type guiState struct {
	title  material.LabelStyle
	status status
	window *app.Window
	theme  *material.Theme
	list   *layout.List

	//conditions to be accepted before processing install
	conditions       []layout.Widget
	acceptButton     *widget.Clickable
	acceptButtonText string
	onAccept         func()

	steps   []layout.Widget
	succeed layout.Widget
	failed  layout.Widget
}

type statusIndicator struct {
	loading material.LoaderStyle
	success *widget.Icon
	fail    *widget.Icon
}

type guiStepState struct {
	step            step
	done            bool
	succeed         bool
	statusIndicator statusIndicator
}

func (i *installer) initGUIState() *guiState {
	theme := material.NewTheme(gofont.Collection())
	statusIndicator := getStatusIndicator(theme)
	guiStepStates := getGUIStepStates(i.steps, statusIndicator)

	state := &guiState{
		title:  getMaterialH5Title(i.title, theme),
		status: getInitialInstallStatus(i.conditions),
		theme:  theme,

		list: &layout.List{Axis: layout.Vertical},

		conditions:       getConditionWidgets(i.conditions, theme),
		acceptButton:     new(widget.Clickable),
		acceptButtonText: i.printer.Sprintf(acceptButtonMsg),

		steps: getGUIStepWidgets(guiStepStates, theme, statusIndicator),

		succeed: material.Body1(theme, i.printer.Sprintf(installSuccessMsg)).Layout,
		failed:  material.Body1(theme, i.printer.Sprintf(installFailMsg)).Layout,
	}
	state.newOnAcceptFunc(guiStepStates)
	return state
}

func (s *guiState) displayAcceptButton(gtx layout.Context) D {
	in := layout.UniformInset(unit.Dp(8))
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return in.Layout(gtx, func(gtx C) D {
				for s.acceptButton.Clicked() {
					s.onAccept()
				}
				dims := material.Button(s.theme, s.acceptButton, s.acceptButtonText).Layout(gtx)
				pointer.CursorNameOp{Name: pointer.CursorPointer}.Add(gtx.Ops)
				return dims
			})
		}),
	)
}

func (s *guiState) Start(windowTitle string) {
	s.window = newWindow(windowTitle)
	go func() {
		if err := s.loop(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

func (s *guiState) loop() error {
	var ops op.Ops
	for {
		select {
		case e := <-s.window.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				s.display(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}
}

func getGUIStepWidgets(states []*guiStepState, theme *material.Theme, indicator statusIndicator) []layout.Widget {
	var widgets []layout.Widget
	for _, s := range states {
		widgets = append(widgets, s.getWidget(theme, indicator)...)
	}
	return widgets
}

func (g *guiStepState) getWidget(theme *material.Theme, indicator statusIndicator) []layout.Widget {
	return []layout.Widget{
		g.getDescriptionWidget(theme),
		g.getStatusIndicatorWidget(indicator),
	}
}

func (g *guiStepState) getDescriptionWidget(theme *material.Theme) func(gtx C) D {
	return getFlexLayoutFunc(layout.Rigid(material.Body2(theme, g.step.description).Layout))
}

func getFlexLayoutFunc(content ...layout.FlexChild) func(gtx C) D {
	return func(gtx C) D {
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx, content...)
	}
}

func (g *guiStepState) getStatusIndicatorWidget(indicator statusIndicator) func(gtx C) D {
	return getFlexLayoutFunc(
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx C) D {
				return g.getStatusIndicatorLayout(indicator, gtx)
			})
		}))

}

func (g *guiStepState) getStatusIndicatorLayout(indicator statusIndicator, gtx C) D {
	if !g.done {
		return indicator.loading.Layout(gtx)
	}
	if g.succeed {
		return indicator.success.Layout(gtx)
	}
	return indicator.fail.Layout(gtx)
}

func (s *guiState) newOnAcceptFunc(steps []*guiStepState) {
	s.onAccept = func() {
		s.status = processing
		go func() {
			err := s.processToInstallEachStep(steps)
			s.status = getInstallDoneStatus(err)
			s.autoCloseWindow()
		}()
	}
}

func getGUIStepStates(steps []step, indicator statusIndicator) []*guiStepState {
	var states []*guiStepState
	for _, s := range steps {
		states = append(states, &guiStepState{
			statusIndicator: indicator,
			step:            s,
		})
	}
	return states
}

func (s *guiState) processToInstallEachStep(steps []*guiStepState) error {
	if len(steps) == 0 {
		return nil
	}
	nextStep := steps[0]
	err := nextStep.step.process()
	nextStep.done = true
	nextStep.succeed = err == nil
	if err != nil {
		return err
	}
	return s.processToInstallEachStep(steps[1:])
}

func getInstallDoneStatus(err error) status {
	if err != nil {
		return failed
	}
	return succeed
}

func (s *guiState) autoCloseWindow() {
	time.Sleep(time.Second * 5)
	s.window.Close()
}

func getConditionWidgets(conditions []condition, theme *material.Theme) []layout.Widget {
	var displayConditions []layout.Widget
	for _, c := range conditions {
		displayConditions = append(displayConditions, getSingleConditionWidget(c, theme)...)
	}
	return displayConditions
}

func getSingleConditionWidget(cond condition, th *material.Theme) []layout.Widget {
	title := material.H6(th, cond.title)
	title.Color = th.Palette.ContrastBg
	return []layout.Widget{
		title.Layout,
		material.Body1(th, cond.body).Layout,
	}
}

func getMaterialH5Title(title string, theme *material.Theme) material.LabelStyle {
	t := material.H5(theme, title)
	t.Color = theme.Palette.ContrastBg
	return t
}

func newWindow(title string) *app.Window {
	windowTitle := app.Title(title)
	size := app.Size(unit.Dp(300), unit.Dp(500))
	return app.NewWindow(size, windowTitle, app.Centered(true))
}

func getInitialInstallStatus(conditions []condition) status {
	if len(conditions) == 0 {
		return processing
	}
	return awaitAcceptCondition
}

func getStatusIndicator(theme *material.Theme) statusIndicator {
	success, _ := widget.NewIcon(icons.ActionDone)
	success.Color = color.NRGBA{R: 0, G: 128, B: 0, A: 0xff}
	fail, _ := widget.NewIcon(icons.AlertError)
	fail.Color = color.NRGBA{R: 255, G: 0, B: 0, A: 0xff}
	return statusIndicator{
		success: success,
		fail:    fail,
		loading: material.Loader(theme),
	}
}

func displayWidgetList(gtx layout.Context, list *layout.List, widgets []layout.Widget) layout.Dimensions {
	return list.Layout(gtx, len(widgets), func(gtx C, i int) D {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, widgets[i])
	})
}

type (
	D = layout.Dimensions
	C = layout.Context
)
