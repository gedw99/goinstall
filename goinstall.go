package goinstall

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"time"
)

type installer struct {
	title string
	conditions []condition
	steps      []step
	printer    *message.Printer
}

type condition struct {
	title string
	body  string
}

type step struct {
	description string
	process     func() error
}

//New initiates the installer.
//Installer can be completed with conditions that user will need to agree before process installation.
//Queued steps can be added with AddStep or with common functions provided in this package (like create file / or dir / or register a new protocol).
//Installer GUI window can be displayed with StartWithGui function.
func New(title string, lang language.Tag) *installer {
	return &installer{title: title, printer: newIntlPrinter(lang)}
}

//AddCondition adds a condition that user will need to agree before process installation.
func (i *installer) AddCondition(title, body string) {
	i.conditions = append(i.conditions, condition{title: title, body: body})
}

//AddStep add a step to process during the install. A generic process function must be provided.
func (i *installer) AddStep(proc func() error, description string) {
	i.steps = append(i.steps, step{description: description, process: func() error {
		//artificial delay to give user feedback
		time.Sleep(2 * time.Second)
		return proc()
	}})
}

//StartWithGUI starts the installer with graphical user interface.
//If conditions are provided, user will have to accept them before proceeding further.
func (i *installer) StartWithGUI() error {
	state := i.initGUIState()
	state.Start(i.printer.Sprintf(windowTitleMsg))
	return nil
}

