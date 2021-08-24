package content_type

import "github.com/microcosm-cc/bluemonday"

func UGCPolicy() {
	p = bluemonday.UGCPolicy()

	// HTML email frequently contains obselete and basic HTML
	p.AllowElements("html", "head", "body", "label", "input", "font", "main", "nav", "header", "footer", "kbd", "legend", "map", "title")

	p.AllowAttrs("type").OnElements("style")

	// Customize attributes on elements
	p.AllowAttrs("type", "media").OnElements("style")
	p.AllowAttrs("face", "size").OnElements("font")
	p.AllowAttrs("name", "content", "http-equiv").OnElements("meta")
	p.AllowAttrs("name", "data-id").OnElements("a")
	p.AllowAttrs("for").OnElements("label")
	p.AllowAttrs("type").OnElements("input")
	p.AllowAttrs("rel", "href").OnElements("link")
	p.AllowAttrs("topmargin", "leftmargin", "marginwidth", "marginheight", "yahoo").OnElements("body")
	p.AllowAttrs("xmlns").OnElements("html")

	p.AllowAttrs("style", "vspace", "hspace", "face", "bgcolor", "color", "border", "cellpadding", "cellspacing").Globally()

	// HTML email tends to see the use of obselete spacing and styling attributes
	p.AllowAttrs("bgcolor", "color", "align").OnElements("basefont", "font", "hr", "table", "td")
	p.AllowAttrs("border").OnElements("img", "table", "basefont", "font", "hr", "td")
	p.AllowAttrs("cellpadding", "cellspacing", "valign", "halign").OnElements("table")

	// Allow "class" attributes on all elements
	p.AllowStyling()

	// Allow images to be embedded via data-uri
	p.AllowDataURIImages()

	// Add "rel=nofollow" to links
	p.RequireNoFollowOnLinks(true)
	p.RequireNoFollowOnFullyQualifiedLinks(true)

	// Open external links in a new window/tab
	p.AddTargetBlankToFullyQualifiedLinks(true)
}
