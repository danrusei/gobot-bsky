package gobotbsky

import (
	"net/url"
	"time"

	appbsky "github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/util"
)

var EmbedExternal appbsky.EmbedExternal
var EmbedExternal_External appbsky.EmbedExternal_External
var EmbedImages appbsky.EmbedImages
var EmbedImages_Image appbsky.EmbedImages_Image
var FeedPost_Embed appbsky.FeedPost_Embed

// construct the post
type PostBuilder struct {
	Text           string
	Link           Link
	Images         []Image
	UploadedImages []lexutil.LexBlob
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

func NewPostBuilder(text string) *PostBuilder {
	return &PostBuilder{
		Text: text,
	}
}

func (pb *PostBuilder) WithExternalLink(title string, link url.URL, description string) *PostBuilder {

	return &PostBuilder{
		Link: Link{
			Title:       title,
			Uri:         link,
			Description: description,
		},
	}
}

func (pb *PostBuilder) WithImages(blobs []lexutil.LexBlob, images ...Image) *PostBuilder {

	return &PostBuilder{
		Images:         images,
		UploadedImages: blobs,
	}
}

func (pb *PostBuilder) Build() appbsky.FeedPost {

	post := appbsky.FeedPost{}

	post.LexiconTypeID = "app.bsky.feed.post"
	post.CreatedAt = time.Now().Format(util.ISO8601)

	if pb.Link != (Link{}) {

		EmbedExternal_External.Title = pb.Link.Title
		EmbedExternal_External.Uri = pb.Link.Uri.String()
		EmbedExternal_External.Description = pb.Link.Description

		EmbedExternal.LexiconTypeID = "app.bsky.embed.external"
		EmbedExternal.External = &EmbedExternal_External

	}

	if len(pb.Images) != 0 && len(pb.Images) == len(pb.UploadedImages) {
		images := []*appbsky.EmbedImages_Image{}

		for i, img := range pb.Images {
			EmbedImages_Image.Alt = img.Title
			EmbedImages_Image.Image = &pb.UploadedImages[i]
			images = append(images, &EmbedImages_Image)
		}

		EmbedImages.LexiconTypeID = "app.bsky.embed.images"
		EmbedImages.Images = images
	}

	FeedPost_Embed.EmbedExternal = &EmbedExternal
	FeedPost_Embed.EmbedImages = &EmbedImages

	post.Embed = &FeedPost_Embed

	return post
}
