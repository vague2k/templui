package shared

type SideLink struct {
	Text  string
	Href  string
	Icon  string
	Click string
}

type Section struct {
	Title string
	Links []SideLink
}

var Sections = []Section{
	{
		Title: "Getting Started",
		Links: []SideLink{
			{
				Text: "Introduction",
				Href: "/docs/introduction",
			},
			{
				Text: "Installation",
				Href: "/docs/installation",
			},
		},
	},
	{
		Title: "Components",
		Links: []SideLink{
			{
				Text: "Button",
				Href: "/docs/components/button",
			},
			{
				Text: "Sheet",
				Href: "/docs/components/sheet",
			},
		},
	},
}
