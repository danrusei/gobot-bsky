package gobotbsky

import (
	"net/url"
	"time"

	appbsky "github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/util"
)

var EmbedExternal appbsky.EmbedExternal
var EmbedExternal_External appbsky.EmbedExternal_External
var FeedPost_Embed appbsky.FeedPost_Embed

// construct the post
type PostBuilder struct {
	Text      string
	EmbedLink EmbedLink
	//EmbedImage EmbedImage
}

type EmbedLink struct {
	Title       string
	Uri         url.URL
	Description string
}

func NewPostBuilder(text string) *PostBuilder {
	return &PostBuilder{
		Text: text,
	}
}

func (pb *PostBuilder) WithExternalLink(title string, link url.URL, description string) *PostBuilder {

	return &PostBuilder{
		EmbedLink: EmbedLink{
			Title:       title,
			Uri:         link,
			Description: description,
		},
	}
}

func (pb *PostBuilder) Build() appbsky.FeedPost {

	post := appbsky.FeedPost{}

	post.LexiconTypeID = "app.bsky.feed.post"
	post.CreatedAt = time.Now().Format(util.ISO8601)

	if pb.EmbedLink != (EmbedLink{}) {

		EmbedExternal_External.Title = pb.EmbedLink.Title
		EmbedExternal_External.Uri = pb.EmbedLink.Uri.String()
		EmbedExternal_External.Description = pb.EmbedLink.Description

		EmbedExternal.LexiconTypeID = "app.bsky.embed.external"
		EmbedExternal.External = &EmbedExternal_External

	}

	FeedPost_Embed.EmbedExternal = &EmbedExternal

	post.Embed = &FeedPost_Embed

	return post
}
