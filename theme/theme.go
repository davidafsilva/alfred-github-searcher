package theme

import (
	"fmt"

	aw "github.com/deanishe/awgo"
)

type Icons struct {
	Repository       *aw.Icon
	DraftPullRequest *aw.Icon
	PullRequest      *aw.Icon
}

type Theme struct {
	Icons *Icons
}

var dark = createTheme("dark")
var light = createTheme("light")

const (
	themeDark  = "dark"
	themeLight = "light"

	themeKey     = "ags_theme"
	defaultTheme = themeLight
)

func New(wf *aw.Workflow) *Theme {
	t := wf.Config.GetString(themeKey, defaultTheme)
	switch t {
	case themeDark:
		return dark
	default:
		return light
	}
}

func createTheme(classifier string) *Theme {
	return &Theme{
		Icons: &Icons{
			Repository: &aw.Icon{
				Value: fmt.Sprintf("repository-%s.png", classifier),
				Type:  aw.IconTypeImage,
			},
			DraftPullRequest: &aw.Icon{
				Value: fmt.Sprintf("draft-pr-%s.png", classifier),
				Type:  aw.IconTypeImage,
			},
			PullRequest: &aw.Icon{
				Value: fmt.Sprintf("pr-%s.png", classifier),
				Type:  aw.IconTypeImage,
			},
		},
	}
}
