package pages

type Page uint8

const (
	OverviewPage Page = iota
	ChatPage
	SettingsPage
)

type ezAppState struct {
	CurrentPage Page
}

var AppState = &ezAppState{
	CurrentPage: OverviewPage,
}
