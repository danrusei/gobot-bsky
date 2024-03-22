package gobotbsky

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	appbsky "github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/util"
)

type Facet_Type int

const (
	Facet_Link Facet_Type = iota + 1
	Facet_Mention
	Facet_Tag
)

var FeedPost_Embed appbsky.FeedPost_Embed

// construct the post
type PostBuilder struct {
	Text           string
	Facet          []Facet
	Link           Link
	Images         []Image
	UploadedImages []lexutil.LexBlob
}

type Facet struct {
	Ftype   Facet_Type
	Value   string
	T_facet string
}

type Link struct {
	Title       string
	Uri         url.URL
	Description string
}

type Image struct {
	Title string
	Uri   url.URL
}

// Create a simple post with text
func NewPostBuilder(text string) PostBuilder {
	return PostBuilder{
		Text:  text,
		Facet: []Facet{},
	}
}

// Create a Richtext Post with facests
func (pb PostBuilder) WithFacet(ftype Facet_Type, value string, text string) PostBuilder {

	pb.Facet = append(pb.Facet, Facet{
		Ftype:   ftype,
		Value:   value,
		T_facet: text,
	})

	return pb
}

// Create a Post with external links
func (pb PostBuilder) WithExternalLink(title string, link url.URL, description string) PostBuilder {

	pb.Link.Title = title
	pb.Link.Uri = link
	pb.Link.Description = description

	return pb
}

// Create a Post with images
func (pb PostBuilder) WithImages(blobs []lexutil.LexBlob, images []Image) PostBuilder {

	pb.Images = images
	pb.UploadedImages = blobs

	return pb
}

// Build the request
func (pb PostBuilder) Build() (appbsky.FeedPost, error) {

	post := appbsky.FeedPost{}

	post.Text = pb.Text
	post.LexiconTypeID = "app.bsky.feed.post"
	post.CreatedAt = time.Now().Format(util.ISO8601)

	// RichtextFacet Section
	// https://docs.bsky.app/docs/advanced-guides/post-richtext

	Facets := []*appbsky.RichtextFacet{}

	for _, f := range pb.Facet {
		facet := &appbsky.RichtextFacet{}
		features := []*appbsky.RichtextFacet_Features_Elem{}
		feature := &appbsky.RichtextFacet_Features_Elem{}

		switch f.Ftype {

		case Facet_Link:
			{
				feature = &appbsky.RichtextFacet_Features_Elem{
					RichtextFacet_Link: &appbsky.RichtextFacet_Link{
						LexiconTypeID: f.Ftype.String(),
						Uri:           f.Value,
					},
				}
			}

		case Facet_Mention:
			{
				feature = &appbsky.RichtextFacet_Features_Elem{
					RichtextFacet_Mention: &appbsky.RichtextFacet_Mention{
						LexiconTypeID: f.Ftype.String(),
						Did:           f.Value,
					},
				}
			}

		case Facet_Tag:
			{
				feature = &appbsky.RichtextFacet_Features_Elem{
					RichtextFacet_Tag: &appbsky.RichtextFacet_Tag{
						LexiconTypeID: f.Ftype.String(),
						Tag:           f.Value,
					},
				}
			}

		}

		features = append(features, feature)
		facet.Features = features

		ByteStart, ByteEnd, err := findSubstring(post.Text, f.T_facet)
		if err != nil {
			return post, fmt.Errorf("unable to find the substring: %v , %v", f.T_facet, err)
		}

		index := &appbsky.RichtextFacet_ByteSlice{
			ByteStart: int64(ByteStart),
			ByteEnd:   int64(ByteEnd),
		}
		facet.Index = index

		Facets = append(Facets, facet)
	}

	post.Facets = Facets

	// Embed Section (either external links or images)
	// As of now it allows only one Embed type per post:
	// https://github.com/bluesky-social/indigo/blob/main/api/bsky/feedpost.go
	if pb.Link != (Link{}) {
		var EmbedExternal appbsky.EmbedExternal
		var EmbedExternal_External appbsky.EmbedExternal_External

		EmbedExternal_External.Title = pb.Link.Title
		EmbedExternal_External.Uri = pb.Link.Uri.String()
		EmbedExternal_External.Description = pb.Link.Description

		EmbedExternal.LexiconTypeID = "app.bsky.embed.external"
		EmbedExternal.External = &EmbedExternal_External

		FeedPost_Embed.EmbedExternal = &EmbedExternal

	} else {
		if len(pb.Images) != 0 && len(pb.Images) == len(pb.UploadedImages) {
			var EmbedImages appbsky.EmbedImages
			var EmbedImages_Image appbsky.EmbedImages_Image
			images := []*appbsky.EmbedImages_Image{}

			for i, img := range pb.Images {
				EmbedImages_Image.Alt = img.Title
				EmbedImages_Image.Image = &pb.UploadedImages[i]
				images = append(images, &EmbedImages_Image)
			}

			EmbedImages.LexiconTypeID = "app.bsky.embed.images"
			EmbedImages.Images = images

			FeedPost_Embed.EmbedImages = &EmbedImages
		}
	}

	// avoid error when trying to marshal empty field (*bsky.FeedPost_Embed)
	if len(pb.Images) != 0 || pb.Link.Title != "" {
		post.Embed = &FeedPost_Embed
	}

	return post, nil
}

func (f Facet_Type) String() string {
	switch f {
	case Facet_Link:
		return "app.bsky.richtext.facet#link"
	case Facet_Mention:
		return "app.bsky.richtext.facet#mention"
	case Facet_Tag:
		return "app.bsky.richtext.facet#tag"
	default:
		return "Unknown"
	}
}
func findSubstring(s, substr string) (int, int, error) {
	index := strings.Index(s, substr)
	if index == -1 {
		return 0, 0, errors.New("substring not found")
	}
	return index, index + len(substr), nil
}
