package common

import "src.sourcegraph.com/apps/tracker/issues"

type State struct {
	BaseURI     string
	ReqPath     string
	CurrentUser *issues.User
}
