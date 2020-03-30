package event

type Track struct {
	Key         string `xml:"key,attr"`
	Type        string `xml:"type,attr"`
	Title       string `xml:"title,attr"`
	GrandParent string `xml:"grandparentTitle,attr"`
	Parent      string `xml:"parentTitle,attr"`
}

type Part struct {
}

type Video struct {
	Part        Part
	Key         string `xml:"key,attr"`
	Title       string `xml:"title,attr"`
	GrandParent string `xml:"grandparentTitle,attr"`
	Parent      string `xml:"parentTitle,attr"`
}

type MediaContainer struct {
	LibrarySectionTitle string `xml:"librarySectionTitle,attr"`
	Title1              string `xml:"title1,attr"`
	Title2              string `xml:"title2,attr"`
	Track               Track
	Video               Video
}
