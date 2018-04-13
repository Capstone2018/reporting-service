package pages

import (
	"database/sql/driver"
	"encoding/json"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// OpenGraph represents opengraph protocol properties,
// describes objects in the semantic web
type OpenGraph struct {
	isArticle        bool
	isBook           bool
	isProfile        bool
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
	Icon             *Image     `json:"icon,omitempty" db:"icon"`
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
	Width     int    `json:"width" db:"width"`
	Height    int    `json:"height" db:"height"`
	Alt       string `json:"alt" db:"alt"`
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
	Width     int    `json:"width" db:"width"`
	Height    int    `json:"height" db:"height"`
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
	Authors        ProfileSlice `json:"authors" db:"authors"`
	PublishedTime  time.Time    `json:"published_time" db:"published_time"`
	ModifiedTime   time.Time    `json:"modified_time" db:"modified_time"`
	ExpirationTime time.Time    `json:"expiration_time" db:"expiration_time"`
	Section        string       `json:"section" db:"section"`
	Tags           []string     `json:"tags" db:"tags"`
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
	ISBN        string       `json:"isbn" db:"isbn"`
	ReleaseDate time.Time    `json:"release_date" db:"release_date"`
	Tags        []string     `json:"tags" db:"tags"`
	Authors     ProfileSlice `json:"authors" db:"authors"`
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
	baseURL, err := url.Parse(pageURL)
	if err != nil {
		og = nil
		return err
	}

	tokenizer := html.NewTokenizer(htmlStream)
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			// return the error if we encounter it
			og = nil
			return tokenizer.Err()
		case html.EndTagToken:
			// end early
			if tokenizer.Token().Data == "head" {
				// set fallbacks
				if len(og.URL) == 0 {
					og.URL = baseURL.String()
				}
				if len(og.Type) == 0 {
					og.Type = "website"
				}
				return nil
			}
		case html.StartTagToken, html.SelfClosingTagToken:
			tagName, hasAttr := tokenizer.TagName()
			switch string(tagName) {
			case "meta":
				if hasAttr {
					og.parseMeta(baseURL, getAttrs(tokenizer))
				}
			case "link":
				if hasAttr {
					og.parseLink(baseURL, getAttrs(tokenizer))
				}
			case "title":
				if len(og.Title) == 0 {
					if tokenizer.Next() == html.TextToken {
						og.Title = tokenizer.Token().Data
					}
				}
			}
		}
	}
}

// return a map of all the attributes for a tag
func getAttrs(tokenizer *html.Tokenizer) map[string]string {
	attrs := make(map[string]string)
	more := true
	for more {
		var key, val []byte
		key, val, more = tokenizer.TagAttr()
		attrs[string(key)] = string(val)
	}
	return attrs
}

// parse an icon link
func (og *OpenGraph) parseLink(baseURL *url.URL, attrs map[string]string) error {
	if attrs["rel"] == "icon" || attrs["rel"] == "shortcut icon" {
		iconURL, err := url.Parse(attrs["href"])
		if err != nil {
			return err
		}
		og.Icon = &Image{
			URL: baseURL.ResolveReference(iconURL).String(),
		}
		//if there is a "sizes" attribute,
		//and it's set to something other than "any"
		//parse it and set width/height
		sizes := attrs["sizes"]
		if len(sizes) > 0 && sizes != "any" {
			og.Icon.Height, og.Icon.Width = parseSizesAttr(sizes)
		}
		if len(attrs["type"]) > 0 {
			og.Icon.Type = attrs["type"]
		}
	}
	return nil
}

// parse a size when hxw
func parseSizesAttr(value string) (int, int) {
	dims := strings.Split(strings.ToLower(value), "x")
	h, _ := strconv.Atoi(dims[0])
	w, _ := strconv.Atoi(dims[1])
	return h, w
}

// set a URL string pointer with a resolved absolute url
func setURL(baseURL *url.URL, structURL *string, content string) {
	cannonicalURL, err := url.Parse(content)
	if err == nil {
		*structURL = baseURL.ResolveReference(cannonicalURL).String()
	}
}

// parse the opengraph attributes of a meta tag
func (og *OpenGraph) parseMeta(baseURL *url.URL, attrs map[string]string) {
	// get the property and content
	content := attrs["content"]
	if len(content) == 0 {
		return
	}
	property := attrs["property"]
	if len(property) == 0 {
		return
	}
	// handle all opengraph tags
	if strings.HasPrefix(property, "og:") {
		og.parseOg(baseURL, property, content)
	} else if strings.HasPrefix(property, "twitter:") { // handle all twitter tags
		og.parseTwitter(baseURL, property, content)
	} else { // check to parse object types
		if property == "description" && len(og.Description) == 0 {
			og.Description = content
		} else if og.isArticle {
			og.parseArticle(property, content)
		} else if og.isBook {
			og.parseBook(property, content)
		} else if og.isProfile {
			og.parseProfile(property, content)
		}
	}
}

// parse all opengraph meta tags
func (og *OpenGraph) parseOg(baseURL *url.URL, property, content string) {
	switch property {
	case "og:title":
		og.Title = content
	case "og:type":
		og.Type = content
		switch og.Type {
		case "article":
			og.isArticle = true
		case "book":
			og.isBook = true
		case "profile":
			og.isProfile = true
		}
	case "og:url":
		setURL(baseURL, &og.URL, content)
	case "og:description":
		og.Description = content
	case "og:determiner":
		og.Determiner = content
	case "og:site_name":
		og.SiteName = content
	case "og:locale":
		og.Locale = content
	case "og:locale:alternate":
		og.LocalesAlternate = append(og.LocalesAlternate, content)
	default:
		og.parseStructured(baseURL, property, content)
	}
}

// parse twitter tags
func (og *OpenGraph) parseTwitter(baseURL *url.URL, property, content string) {
	switch property {
	case "twitter:card":
		if len(og.Type) == 0 {
			og.Type = content
		}
	case "twitter:description":
		if len(og.Description) == 0 {
			og.Description = content
		}
	case "twitter:title":
		if len(og.Title) == 0 {
			og.Title = content
		}
	case "twitter:image":
		for _, image := range og.Images {
			if image.URL == content {
				return
			}
		}
		image := &Image{}
		image.URL = content
		og.Images = append(og.Images, image)
	}
}

// parse a structured opengraph object
func (og *OpenGraph) parseStructured(baseURL *url.URL, property, content string) {
	if strings.HasPrefix(property, "og:image") {
		og.parseImage(baseURL, property, content)
	} else if strings.HasPrefix(property, "og:video") {
		og.parseVideo(baseURL, property, content)
	} else if strings.HasPrefix(property, "og:audio") {
		og.parseAudio(baseURL, property, content)
	}
}

// parse an opengraph image
func (og *OpenGraph) parseImage(baseURL *url.URL, property, content string) {
	if property == "og:image" {
		og.Images = append(og.Images, &Image{})
	}
	// get the last image
	image := og.Images[len(og.Images)-1]
	switch property {
	case "og:image", "og:image:url":
		setURL(baseURL, &image.URL, content)
	case "og:image:secure_url":
		setURL(baseURL, &image.SecureURL, content)
	case "og:image:type":
		image.Type = content
	case "og:image:width":
		w, _ := strconv.Atoi(content)
		image.Width = w
	case "og:image:height":
		h, _ := strconv.Atoi(content)
		image.Height = h
	case "og:image:alt":
		image.Alt = content
	}
}

// parse an opengraph video
func (og *OpenGraph) parseVideo(baseURL *url.URL, property, content string) {
	if property == "og:video" {
		og.Videos = append(og.Videos, &Video{})
	}
	// get the last video
	video := og.Videos[len(og.Videos)-1]
	switch property {
	case "og:video", "og:video:url":
		setURL(baseURL, &video.URL, content)
	case "og:video:secure_url":
		setURL(baseURL, &video.SecureURL, content)
	case "og:video:type":
		video.Type = content
	case "og:video:width":
		w, _ := strconv.Atoi(content)
		video.Width = w
	case "og:video:height":
		h, _ := strconv.Atoi(content)
		video.Height = h
	}
}

// parse an opengraph audio
func (og *OpenGraph) parseAudio(baseURL *url.URL, property, content string) {
	if property == "og:audio" {
		og.Audios = append(og.Audios, &Audio{})
	}
	audio := og.Audios[len(og.Audios)-1]
	switch property {
	case "og:audio", "og:audio:url":
		setURL(baseURL, &audio.URL, content)
	case "og:audio:secure_url":
		setURL(baseURL, &audio.SecureURL, content)
	case "og:audio:type":
		audio.Type = content
	}
}

// parse an opengraph article
func (og *OpenGraph) parseArticle(property, content string) {
	if og.Article == nil {
		og.Article = &Article{}
	}
	switch property {
	case "article:published_time":
		if t, err := time.Parse(time.RFC3339, content); err == nil {
			og.Article.PublishedTime = t
		}
	case "article:modified_time":
		if t, err := time.Parse(time.RFC3339, content); err == nil {
			og.Article.ModifiedTime = t
		}
	case "article:expiration_time":
		if t, err := time.Parse(time.RFC3339, content); err == nil {
			og.Article.ExpirationTime = t
		}
	case "article:section":
		og.Article.Section = content
	case "article:tag":
		og.Article.Tags = append(og.Article.Tags, content)
	case "article:author:first_name":
		appendProfile(&og.Article.Authors)
		og.Article.Authors[len(og.Article.Authors)-1].FirstName = content
	case "article:author:last_name":
		appendProfile(&og.Article.Authors)
		og.Article.Authors[len(og.Article.Authors)-1].LastName = content
	case "article:author:username":
		appendProfile(&og.Article.Authors)
		og.Article.Authors[len(og.Article.Authors)-1].Username = content
	case "article:author:gender":
		appendProfile(&og.Article.Authors)
		og.Article.Authors[len(og.Article.Authors)-1].Gender = content
	}
}

// appendProfile appends a profile to a pointer to profile slice if len == 0
func appendProfile(profiles *ProfileSlice) {
	if len(*profiles) == 0 {
		*profiles = append(*profiles, &Profile{})
	}
}

// parse an opengraph book
func (og *OpenGraph) parseBook(property, content string) {
	if og.Book == nil {
		og.Book = &Book{}
	}
	switch property {
	case "book:isbn":
		og.Book.ISBN = content
	case "book:release_date":
		if t, err := time.Parse(time.RFC3339, content); err == nil {
			og.Book.ReleaseDate = t
		}
	case "book:tag":
		og.Book.Tags = append(og.Book.Tags, content)
	case "book:author:first_name":
		appendProfile(&og.Book.Authors)
		og.Book.Authors[len(og.Book.Authors)-1].FirstName = content
	case "book:author:last_name":
		appendProfile(&og.Book.Authors)
		og.Book.Authors[len(og.Book.Authors)-1].LastName = content
	case "book:author:username":
		appendProfile(&og.Book.Authors)
		og.Book.Authors[len(og.Book.Authors)-1].Username = content
	case "book:author:gender":
		appendProfile(&og.Book.Authors)
		og.Book.Authors[len(og.Book.Authors)-1].Gender = content
	}
}

// parse an opengraph profile
func (og *OpenGraph) parseProfile(property, content string) {
	if og.Profile == nil {
		og.Profile = &Profile{}
	}
	switch property {
	case "profile:first_name":
		og.Profile.FirstName = content
	case "profile:last_name":
		og.Profile.LastName = content
	case "profile:username":
		og.Profile.Username = content
	case "profile:gender":
		og.Profile.Gender = content
	}
}
