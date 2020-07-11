package webapi

import (
	"encoding/json"
	"github.com/oklahomer/golack/event"
	"net/url"
	"strconv"
)

// ParseMode defines the parse parameter values for post.message method.
// See https://api.slack.com/docs/message-formatting#parsing_modes
type ParseMode string

const (
	ParseModeNone ParseMode = "none"
	ParseModeFull ParseMode = "full"
)

// String returns a stringified form of BotType
func (mode ParseMode) String() string {
	return string(mode)
}

type AttachmentField struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value"`
	Short bool   `json:"short,omitempty"`
}

type MessageAttachment struct {
	Fallback   string             `json:"fallback"`
	Color      string             `json:"color,omitempty"`
	Pretext    string             `json:"pretext,omitempty"`
	AuthorName string             `json:"author_name,omitempty"`
	AuthorLink string             `json:"author_link,omitempty"`
	AuthorIcon string             `json:"author_icon,omitempty"`
	Title      string             `json:"title,omitempty"`
	TitleLink  string             `json:"title_link,omitempty"`
	Text       string             `json:"text,omitempty"`
	Fields     []*AttachmentField `json:"fields"`
	ImageURL   string             `json:"image_url,omitempty"`
	ThumbURL   string             `json:"thumb_url,omitempty"`
}

// PostMessage is a payload to be sent with chat.postMessage method.
// See https://api.slack.com/methods/chat.postMessage
type PostMessage struct {
	ChannelID       event.ChannelID      `json:"channel"`
	Text            string               `json:"text"`
	Parse           ParseMode            `json:"parse,omitempty"`
	LinkNames       int                  `json:"link_names,omitempty"`
	Attachments     []*MessageAttachment `json:"attachments,omitempty"`
	UnfurlLinks     bool                 `json:"unfurl_links,omitempty"`
	UnfurlMedia     bool                 `json:"unfurl_media,omitempty"`
	UserName        string               `json:"username,omitempty"`
	AsUser          bool                 `json:"as_user,omitempty"`
	IconURL         string               `json:"icon_url,omitempty"`
	IconEmoji       string               `json:"icon_emoji,omitempty"`
	ReplyBroadcast  bool                 `json:"reply_broadcast,omitempty"`
	ThreadTimeStamp string               `json:"thread_ts,omitempty"`
}

// WithAttachments sets/overrides attachments parameter for current PostMessage.
// See https://api.slack.com/docs/message-attachments
func (message *PostMessage) WithAttachments(attachments []*MessageAttachment) *PostMessage {
	message.Attachments = attachments
	return message
}

// WithLinkNames sets/overrides link_names parameter for current PostMessage.
// See https://api.slack.com/methods/chat.postMessage#formatting
func (message *PostMessage) WithLinkNames(linkNames int) *PostMessage {
	message.LinkNames = linkNames
	return message
}

// WithParse sets/overrides parse parameter for current PostMessage.
// See https://api.slack.com/docs/message-formatting#parsing_modes
func (message *PostMessage) WithParse(parse ParseMode) *PostMessage {
	message.Parse = parse
	return message
}

// WithReplyBroadcast sets optional boolean value so the thread response can be broadcasted.
// Thread identifier must be present with WithThreadTimeStamp() to use this option.
// See https://api.slack.com/docs/message-threading#using_the_web_api
func (message *PostMessage) WithReplyBroadcast(broadcast bool) *PostMessage {
	message.ReplyBroadcast = broadcast
	return message
}

// WithThreadTimeStamp sets given ts value to payload.
// See https://api.slack.com/docs/message-threading#using_the_web_api
func (message *PostMessage) WithThreadTimeStamp(ts string) *PostMessage {
	message.ThreadTimeStamp = ts
	return message
}

// WithUnfurlLinks sets/overrides unfurl_links for current PostMessage.
// See https://api.slack.com/docs/message-attachments#unfurling
func (message *PostMessage) WithUnfurlLinks(unfurl bool) *PostMessage {
	message.UnfurlLinks = unfurl
	return message
}

// WithUnfurlLinks sets/overrides unfurl_media for current PostMessage.
// See https://api.slack.com/docs/message-attachments#unfurling
func (message *PostMessage) WithUnfurlMedia(unfurl bool) *PostMessage {
	message.UnfurlMedia = unfurl
	return message
}

// ToURLValues forms requesting parameter for Slack's Rest API endpoint.
// See https://api.slack.com/docs/message-formatting
func (message *PostMessage) ToURLValues() url.Values {
	values := url.Values{}

	values.Add("channel", message.ChannelID.String())
	values.Add("text", message.Text)
	values.Add("parse", message.Parse.String())
	values.Add("link_names", strconv.Itoa(message.LinkNames))
	values.Add("unfurl_links", strconv.FormatBool(message.UnfurlLinks))
	values.Add("unfurl_media", strconv.FormatBool(message.UnfurlMedia))
	values.Add("as_user", strconv.FormatBool(message.AsUser))
	if message.ThreadTimeStamp != "" {
		values.Add("thread_ts", message.ThreadTimeStamp)
		values.Add("reply_broadcast", strconv.FormatBool(message.ReplyBroadcast))
	}
	if message.UserName != "" {
		values.Add("username", message.UserName)
	}
	if message.IconURL != "" {
		values.Add("icon_url", message.IconURL)
	}
	if message.IconEmoji != "" {
		values.Add("icon_emoji", message.IconEmoji)
	}
	if message.Attachments != nil {
		s, _ := json.Marshal(message.Attachments)
		values.Add("attachments", string(s))
	}

	return values
}

// NewPostMessage creates PostMessage instance with given channel and text settings.
// By default this sets commonly used settings as much as possible. e.g. link_names=1, unfurl_links=true, etc...
// To override those settings and add some extra settings including username, icon_url, or icon_emoji, call setter methods start with With***.
func NewPostMessage(channelID event.ChannelID, text string) *PostMessage {
	return &PostMessage{
		ChannelID:   channelID,
		Text:        text,
		Parse:       ParseModeFull,
		LinkNames:   1,
		UnfurlLinks: true,
		UnfurlMedia: true,
	}
}
