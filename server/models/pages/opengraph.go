package pages

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"
)

// OpenGraph represents opengraph protocol properties,
// describes objects in the semantic web
type OpenGraph struct {
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	ID               int64      `json:"id" db:"id"`
	Title            string     `json:"title,omitempty" db:"title"`
	Type             string     `json:"type,omitempty" db:"type"`
	URL              string     `json:"url,omitempty" db:"url"`
	Description      string     `json:"description,omitempty" db:"description"`
	Determiner       string     `json:"determiner,omitempty" db:"determiner"`
	SiteName         string     `json:"siteName,omitempty" db:"siteName"`
	Locale           string     `json:"locale,omitempty" db:"locale"`
	LocalesAlternate []string   `json:"locales_alternate,omitempty" db:"locales_alternate"`
	Images           ImageSlice `json:"images,omitempty" db:"images"`
	Audios           AudioSlice `json:"audios,omitempty" db:"audios"`
	Videos           VideoSlice `json:"videos,omitempty" db:"videos"`
	Article          *Article   `json:"article,omitempty" db:"article"`
	Book             *Book      `json:"book,omitempty" db:"book"`
	Profile          *Profile   `json:"profile,omitempty" db:"profile"`
}

// Value implements driver Valuer interface
func (og OpenGraph) Value() (driver.Value, error) {
	j, err := json.Marshal(og)
	return j, err
}

// Image defines Open Graph Image type
type Image struct {
	URL       string `json:"url" db:"url"`
	SecureURL string `json:"secure_url" db:"secure_url"`
	Type      string `json:"type" db:"type"`
	Width     uint64 `json:"width" db:"width"`
	Height    uint64 `json:"height" db:"height"`
}

// Value implements driver Valuer interface
func (i Image) Value() (driver.Value, error) {
	j, err := json.Marshal(i)
	return j, err
}

// ImageSlice is a slice of pointers to image
type ImageSlice []*Image

// Value implements driver Valuer interface
func (i ImageSlice) Value() (driver.Value, error) {
	j, err := json.Marshal(i)
	return j, err
}

// Video defines Open Graph Video type
type Video struct {
	URL       string `json:"url" db:"url"`
	SecureURL string `json:"secureUrl" db:"video"`
	Type      string `json:"type" db:"type"`
	Width     uint64 `json:"width" db:"width"`
	Height    uint64 `json:"height" db:"height"`
}

// Value implements driver Valuer interface
func (v Video) Value() (driver.Value, error) {
	j, err := json.Marshal(v)
	return j, err
}

//VideoSlice allows us to marshall structs when inserting into the database
type VideoSlice []*Video

// Value implements driver Valuer interface
func (v VideoSlice) Value() (driver.Value, error) {
	j, err := json.Marshal(v)
	return j, err
}

// Audio defines Open Graph Audio Type
type Audio struct {
	URL       string `json:"url" db:"url"`
	SecureURL string `json:"secure_url" db:"secure_url"`
	Type      string `json:"type" db:"type"`
}

// Value implements driver Valuer interface
func (a Audio) Value() (driver.Value, error) {
	j, err := json.Marshal(a)
	return j, err
}

// AudioSlice allows us to insert jsonb
type AudioSlice []*Audio

// Value implements driver Valuer interface
func (a AudioSlice) Value() (driver.Value, error) {
	j, err := json.Marshal(a)
	return j, err
}

// Article represents opengraph article properties
type Article struct {
	Authors        []*Profile `json:"authors" db:"authors"`
	PublishedTime  time.Time  `json:"published_time" db:"published_time"`
	ModifiedTime   time.Time  `json:"modified_time" db:"modified_time"`
	ExpirationTime time.Time  `json:"expiration_time" db:"expiration_time"`
	Section        string     `json:"section" db:"section"`
	Tags           []string   `json:"tags" db:"tags"`
}

// Value implements driver Valuer interface
func (a Article) Value() (driver.Value, error) {
	j, err := json.Marshal(a)
	return j, err
}

//ArticleSlice allows us to insert jsonb
type ArticleSlice []*Article

// Value implements driver Valuer interface
func (a ArticleSlice) Value() (driver.Value, error) {
	j, err := json.Marshal(a)
	return j, err
}

// Profile contains Open Graph Profile structure
type Profile struct {
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Username  string `json:"username" db:"username"`
	Gender    string `json:"gender" db:"gender"`
}

// Value implements driver Valuer interface
func (p Profile) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

// ProfileSlice allows us to insert jsonb
type ProfileSlice []*Profile

// Value implements driver Valuer interface
func (p ProfileSlice) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

// Book contains Open Graph Book structure
type Book struct {
	ISBN        string     `json:"isbn" db:"isbn"`
	ReleaseDate *time.Time `json:"release_date" db:"release_date"`
	Tags        []string   `json:"tags" db:"tags"`
	Authors     []*Profile `json:"authors" db:"authors"`
}

// Value implements driver Valuer interface
func (b Book) Value() (driver.Value, error) {
	j, err := json.Marshal(b)
	return j, err
}

// BookSlice is a slice of Books
type BookSlice []*Book

// Value implements driver Valuer interface
func (b BookSlice) Value() (driver.Value, error) {
	j, err := json.Marshal(b)
	return j, err
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
