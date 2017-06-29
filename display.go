package issuesapp

import (
	"fmt"
	"html/template"
	"time"

	"github.com/shurcooL/htmlg"
	"github.com/shurcooL/issues"
	"github.com/shurcooL/issuesapp/component"
	"github.com/shurcooL/octiconssvg"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// issue is an issues.Issue wrapper with display augmentations.
type issue struct {
	issues.Issue

	Unread bool // Unread indicates whether the issue contains unread notifications for authenticated user.
}

// issueItem represents an issue item for display purposes.
type issueItem struct {
	// IssueItem can be one of issues.Comment, event.
	IssueItem interface{}
}

func (i issueItem) TemplateName() string {
	switch i.IssueItem.(type) {
	case issues.Comment:
		return "comment"
	case event:
		return "event"
	default:
		panic(fmt.Errorf("unknown item type %T", i.IssueItem))
	}
}

func (i issueItem) CreatedAt() time.Time {
	switch i := i.IssueItem.(type) {
	case issues.Comment:
		return i.CreatedAt
	case event:
		return i.CreatedAt
	default:
		panic(fmt.Errorf("unknown item type %T", i))
	}
}

func (i issueItem) ID() uint64 {
	switch i := i.IssueItem.(type) {
	case issues.Comment:
		return i.ID
	case event:
		return i.ID
	default:
		panic(fmt.Errorf("unknown item type %T", i))
	}
}

// byCreatedAtID implements sort.Interface.
type byCreatedAtID []issueItem

func (s byCreatedAtID) Len() int { return len(s) }
func (s byCreatedAtID) Less(i, j int) bool {
	if s[i].CreatedAt().Equal(s[j].CreatedAt()) {
		// If CreatedAt time is equal, fall back to ID as a tiebreaker.
		return s[i].ID() < s[j].ID()
	}
	return s[i].CreatedAt().Before(s[j].CreatedAt())
}
func (s byCreatedAtID) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// event is an issues.Event wrapper with display augmentations.
type event struct {
	issues.Event
}

func (e event) Text() template.HTML {
	switch e.Event.Type {
	case issues.Reopened, issues.Closed:
		return template.HTML(htmlg.Render(htmlg.Text(fmt.Sprintf("%s this", e.Event.Type))))
	case issues.Renamed:
		return template.HTML(htmlg.Render(htmlg.Text("changed the title from "), htmlg.Strong(e.Event.Rename.From), htmlg.Text(" to "), htmlg.Strong(e.Event.Rename.To)))
	case issues.Labeled:
		return template.HTML(htmlg.Render(htmlg.Text("added the "), component.Label{Label: *e.Event.Label}.Render()[0], htmlg.Text(" label")))
	case issues.Unlabeled:
		return template.HTML(htmlg.Render(htmlg.Text("removed the "), component.Label{Label: *e.Event.Label}.Render()[0], htmlg.Text(" label")))
	case issues.CommentDeleted:
		return template.HTML(htmlg.Render(htmlg.Text("deleted a comment")))
	default:
		return template.HTML(htmlg.Render(htmlg.Text(string(e.Event.Type))))
	}
}

func (e event) Icon() template.HTML {
	var (
		icon            *html.Node
		color           = "#767676"
		backgroundColor = "#f3f3f3"
	)
	switch e.Event.Type {
	case issues.Reopened:
		icon = octiconssvg.PrimitiveDot()
		color, backgroundColor = "#fff", "#6cc644"
	case issues.Closed:
		icon = octiconssvg.CircleSlash()
		color, backgroundColor = "#fff", "#bd2c00"
	case issues.Renamed:
		icon = octiconssvg.Pencil()
	case issues.Labeled, issues.Unlabeled:
		icon = octiconssvg.Tag()
	case issues.CommentDeleted:
		icon = octiconssvg.X()
	default:
		icon = octiconssvg.PrimitiveDot()
	}
	span := &html.Node{
		Type: html.ElementNode, Data: atom.Span.String(),
		Attr: []html.Attribute{
			{Key: atom.Class.String(), Val: "event-icon"},
			{Key: atom.Style.String(), Val: fmt.Sprintf("color: %s; background-color: %s;", color, backgroundColor)},
		},
		FirstChild: icon,
	}
	return template.HTML(htmlg.Render(span))
}
