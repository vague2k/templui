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
				Text: "How to Use",
				Href: "/docs/how-to-use",
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
				Text: "Card",
				Href: "/docs/components/card",
			},
			{
				Text: "Input",
				Href: "/docs/components/input",
			},
			{
				Text: "Sheet",
				Href: "/docs/components/sheet",
			},
			{
				Text: "Tabs",
				Href: "/docs/components/tabs",
			},
		},
	},
}
