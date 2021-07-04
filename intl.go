package goinstall

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	windowTitleMsg     = "title"
	acceptButtonMsg    = "accept"
	installSuccessMsg  = "success"
	installFailMsg = "fail"
	remakeDirMsg        = "rmk_dir"
	registerProtocolMsg = "register_protocol"
	copyingFilesMsg     = "copying_files"
)

func newIntlPrinter(lang language.Tag) *message.Printer {
	setIntlMessages()
	return message.NewPrinter(lang)
}

func setIntlMessages() {
	setWindowTitleMsg()
	setAcceptButtonMsg()
	setInstallSuccessMsg()
	setInstallFailMsg()
	setRemakeDirMsg()
	setRegisterProtocolMsg()
	setCopyingFilesMsg()
}

func setWindowTitleMsg() {
	m := map[language.Tag]string{
		language.English: "Installation",
		language.French:  "Installation",
	}
	setMessage(windowTitleMsg, m)
}

func setAcceptButtonMsg() {
	m := map[language.Tag]string{
		language.English: "Accept",
		language.French:  "Accepter",
	}
	setMessage(acceptButtonMsg, m)
}

func setInstallSuccessMsg() {
	m := map[language.Tag]string{
		language.English: "Install successful. This window may now be closed.",
		language.French:  "L'installation s'est bien déroulée, vous pouvez désormais fermer cette fenêtre.",
	}
	setMessage(installSuccessMsg, m)
}

func setInstallFailMsg() {
	m := map[language.Tag]string{
		language.English: "Install failed. This window will close itself in a few seconds.",
		language.French:  "L'installation a échoué, la fenêtre se fermera dans quelques secondes.",
	}
	setMessage(installFailMsg, m)
}

func setRemakeDirMsg() {
	m := map[language.Tag]string{
		language.English: "Creating new directory : %s",
		language.French:  "Création du dossier : %s",
	}
	setMessage(remakeDirMsg, m)
}

func setRegisterProtocolMsg() {
	m := map[language.Tag]string{
		language.English: "Installing protocol %s.",
		language.French:  "Installation du protocol %s",
	}
	setMessage(registerProtocolMsg, m)
}

func setCopyingFilesMsg() {
	m := map[language.Tag]string{
		language.English: "Extracting files in folder %s.",
		language.French:  "Extraction des fichiers dans le dossier %s",
	}
	setMessage(copyingFilesMsg, m)
}

func setMessage(msg string, translations map[language.Tag]string) {
	for lang, translation := range translations {
		message.SetString(lang, msg, translation)
	}
}
