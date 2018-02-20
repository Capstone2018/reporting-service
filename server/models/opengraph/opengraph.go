package opengraph

import (
	"fmt"
	"net/url"
	"time"
)

// OpenGraph represents opengraph protocol properties,
// describes objects in the semantic web
type OpenGraph struct {
	Title            string   `json:"title"`
	Type             string   `json:"type"`
	URL              string   `json:"url"`
	Description      string   `json:"description"`
	Determiner       string   `json:"determiner"`
	SiteName         string   `json:"siteName"`
	Locale           string   `json:"locale"`
	LocalesAlternate []string `json:"localesAlternate"`
	Images           []*Image `json:"images"`
	Audios           []*Audio `json:"audios"`
	Videos           []*Video `json:"videos"`
	Article          *Article `json:"article,omitempty"`
	Book             *Book    `json:"book,omitempty"`
	Profile          *Profile `json:"profile,omitempty"`
}

// Image defines Open Graph Image type
type Image struct {
	URL       string `json:"url"`
	SecureURL string `json:"secureUrl"`
	Type      string `json:"type"`
	Width     uint64 `json:"width"`
	Height    uint64 `json:"height"`
}

// Video defines Open Graph Video type
type Video struct {
	URL       string `json:"url"`
	SecureURL string `json:"secureUrl"`
	Type      string `json:"type"`
	Width     uint64 `json:"width"`
	Height    uint64 `json:"height"`
}

// Audio defines Open Graph Audio Type
type Audio struct {
	URL       string `json:"url"`
	SecureURL string `json:"secureUrl"`
	Type      string `json:"type"`
}

// Article represents opengraph article properties
type Article struct {
	Authors        []*Profile `json:"authors"`
	PublishedTime  time.Time  `json:"publishedTime"`
	ModifiedTime   time.Time  `json:"modifiedTime"`
	ExpirationTime time.Time  `json:"expirationTime"`
	Section        string     `json:"section"`
	Tags           []string   `json:"tags"`
}

// Profile contains Open Graph Profile structure
type Profile struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Gender    string `json:"gender"`
}

// Book contains Open Graph Book structure
type Book struct {
	ISBN        string     `json:"isbn"`
	ReleaseDate *time.Time `json:"releaseDate"`
	Tags        []string   `json:"tags"`
	Authors     []*Profile `json:"authors"`
}

// Validate the open graph
func (g *OpenGraph) Validate() error {
	// todo, check 140 characters
	// validate the opengraph properties
	if len(g.Title) > 40 {
		return fmt.Errorf("opengraph title len > 40")
	}
	// _ := map[string]bool{
	// 	"article": true,
	// 	"website": true,
	// 	"image":   true,
	// 	"book":    true,
	// 	"profile": true,
	// }
	if !(g.Type == "article" || g.Type == "book" || g.Type == "profile" || g.Type == "website" || g.Type == "video" || g.Type == "tes") {
		return fmt.Errorf("opengraph type not allowed value")
	}
	// validate image and opengraph URL
	// _, err := url.ParseRequestURI(g.Image.URL)
	// if err != nil {
	// 	return fmt.Errorf("invalid opengraph image url: %v", err)
	// }
	_, err := url.ParseRequestURI(g.URL)
	if err != nil {
		return fmt.Errorf("invalid opengraph url: %v", err)
	}
	if len(g.Description) > 300 {
		return fmt.Errorf("opengraph description len > 300")
	}

	return nil
}
