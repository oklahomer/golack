package webapi

type TimeStamp int64

// APIResponse provides common fields shared by all API response.
// https://api.slack.com/web#basics
type APIResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

// Self contains details on the authenticated user.
// https://api.slack.com/methods/rtm.start#response
type Self struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Created        TimeStamp `json:"created"`
	ManualPresence string    `json:"manual_presence"`
}

type UserProfile struct {
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	RealName           string `json:"real_name"`
	RealNameNormalized string `json:"real_name_normalized"`
	Email              string `json:"email"`
	Skype              string `json:"skype"`
	Phone              string `json:"phone"`
	Image24            string `json:"image_24"`
	Image32            string `json:"image_32"`
	Image48            string `json:"image_48"`
	Image72            string `json:"image_72"`
	Image192           string `json:"image_192"`
	ImageOriginal      string `json:"image_original"`
	Title              string `json:"title"`
}

// User contains all the information of a user
type User struct {
	User              string      `json:"user"`
	Name              string      `json:"name"`
	Deleted           bool        `json:"deleted"`
	Color             string      `json:"color"`
	RealName          string      `json:"real_name"`
	TZ                string      `json:"tz"`
	TZLabel           string      `json:"tz_label"`
	TZOffset          int         `json:"tz_offset"`
	Profile           UserProfile `json:"profile"`
	IsBot             bool        `json:"is_bot"`
	IsAdmin           bool        `json:"is_admin"`
	IsOwner           bool        `json:"is_owner"`
	IsPrimaryOwner    bool        `json:"is_primary_owner"`
	IsRestricted      bool        `json:"is_restricted"`
	IsUltraRestricted bool        `json:"is_ultra_restricted"`
	Has2FA            bool        `json:"has_2fa"`
	HasFiles          bool        `json:"has_files"`
	Presence          string      `json:"presence"`
}

type Team struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

type Topic struct {
	Value   string    `json:"value"`
	Creator string    `json:"creator"`
	LastSet TimeStamp `json:"last_set"`
}

type Purpose struct {
	Value   string    `json:"value"`
	Creator string    `json:"creator"`
	LastSet TimeStamp `json:"last_set"`
}

type Message struct {
	UserID string `json:"user"`
	Text   string `json:"text"`
}

type Channel struct {
	ID                 string    `json:"id"`
	Created            TimeStamp `json:"created"`
	IsOpen             bool      `json:"is_open"`
	LastRead           string    `json:"last_read"`
	Latest             *Message  `json:"latest"`
	UnreadCount        int       `json:"unread_count"`
	UnreadCountDisplay int       `json:"unread_count_display"`
	Name               string    `json:"name"`
	Creator            string    `json:"creator"`
	IsArchived         bool      `json:"is_archived"`
	Members            []string  `json:"members"`
	NumMembers         int       `json:"num_members"`
	Topic              Topic     `json:"topic"`
	Purpose            Purpose   `json:"purpose"`
	IsChannel          bool      `json:"is_channel"`
	IsGeneral          bool      `json:"is_general"`
	IsMember           bool      `json:"is_member"`
}

type Group struct {
	ID                 string    `json:"id"`
	Created            TimeStamp `json:"created"`
	IsOpen             bool      `json:"is_open"`
	LastRead           string    `json:"last_read"`
	Latest             *Message  `json:"latest"`
	UnreadCount        int       `json:"unread_count"`
	UnreadCountDisplay int       `json:"unread_count_display"`
	Name               string    `json:"name"`
	Creator            string    `json:"creator"`
	IsArchived         bool      `json:"is_archived"`
	Members            []string  `json:"members"`
	NumMembers         int       `json:"num_members"`
	Topic              Topic     `json:"topic"`
	Purpose            Purpose   `json:"purpose"`
	IsGroup            bool      `json:"is_group"`
}

type Icons struct {
	Image48 string `json:"image_48"`
}

type Bot struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Deleted bool   `json:"deleted"`
	Icons   Icons  `json:"icons"`
}

type IM struct {
	ID                 string    `json:"id"`
	Created            TimeStamp `json:"created"`
	IsOpen             bool      `json:"is_open"`
	LastRead           string    `json:"last_read"`
	Latest             *Message  `json:"latest"`
	UnreadCount        int       `json:"unread_count"`
	UnreadCountDisplay int       `json:"unread_count_display"`
	IsIM               bool      `json:"is_im"`
	UserID             string    `json:"user"`
	IsUserDeleted      bool      `json:"is_user_deleted"`
}

// https://api.slack.com/methods/rtm.start
type RTMStart struct {
	APIResponse

	// TODO consider net/url
	URL string `json:"url"`

	Self     *Self     `json:"self"`
	Team     *Team     `json:"team"`
	Users    []User    `json:"users"`
	Channels []Channel `json:"channels"`
	Groups   []Group   `json:"groups"`
	Bots     []Bot     `json:"bots"`
	IMs      []IM      `json:"ims"`
}
