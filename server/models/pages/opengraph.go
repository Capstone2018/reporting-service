package pages

import (
	"fmt"
	"io"
	"net/url"
	"time"
)

// OpenGraph represents opengraph protocol properties,
// describes objects in the semantic web
type OpenGraph struct {
	CreatedAt        time.Time `json:"created_at"`
	ID               int64     `json:"id"`
	Title            string    `json:"title,omitempty"`
	Type             string    `json:"type,omitempty"`
	URL              string    `json:"url,omitempty"`
	Description      string    `json:"description,omitempty"`
	Determiner       string    `json:"determiner,omitempty"`
	SiteName         string    `json:"siteName,omitempty"`
	Locale           string    `json:"locale,omitempty"`
	LocalesAlternate []string  `json:"localesAlternate,omitempty"`
	Images           []*Image  `json:"images,omitempty"`
	Audios           []*Audio  `json:"audios,omitempty"`
	Videos           []*Video  `json:"videos,omitempty"`
	Article          *Article  `json:"article,omitempty"`
	Book             *Book     `json:"book,omitempty"`
	Profile          *Profile  `json:"profile,omitempty"`
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

// NewOpenGraph returns new instance of Open Graph structure
func NewOpenGraph() *OpenGraph {
	return &OpenGraph{
		CreatedAt: time.Now(),
	}
}

// ProcessStream parses an HTML stream to generate opengraph properties from it
func (og *OpenGraph) ProcessStream(pageURL string, htmlStream io.ReadCloser) error {

	return nil
}

// Validate the open graph
func (og *OpenGraph) Validate() error {
	// todo, check 140 characters
	// validate the opengraph properties
	if len(og.Title) > 40 {
		return fmt.Errorf("opengraph title len > 40")
	}
	// _ := map[string]bool{
	// 	"article": true,
	// 	"website": true,
	// 	"image":   true,
	// 	"book":    true,
	// 	"profile": true,
	// }
	if !(og.Type == "article" || og.Type == "book" || og.Type == "profile" || og.Type == "website" || og.Type == "video" || og.Type == "tes") {
		return fmt.Errorf("opengraph type not allowed value")
	}
	// validate image and opengraph URL
	// _, err := url.ParseRequestURI(g.Image.URL)
	// if err != nil {
	// 	return fmt.Errorf("invalid opengraph image url: %v", err)
	// }
	_, err := url.ParseRequestURI(og.URL)
	if err != nil {
		return fmt.Errorf("invalid opengraph url: %v", err)
	}
	if len(og.Description) > 300 {
		return fmt.Errorf("opengraph description len > 300")
	}

	return nil
}
