package models

type LoginUserData struct{
	Username string `json:"username"`
	Discriminator string `json:"discriminator"`
	MfaEnabled bool `json:"mfa_enabled"`
	Id string `json:"id"`
	Avatar string `json:"avatar"`
}

type GuildWeb struct{
	Id string `json:"id"`
	SecurityLevel int `json:"securityLevel"`
	ReportCount int `json:"reportCount"`
	Checked bool `json:"checked"`
	Reports []ReportWeb `json:"reports"`
}

type ReportWeb struct{
	Text string `json:"text"`
	Date int `json:"date"`
	Closed bool`json:"closed"`
	User GenUserDataPost `json:"user"`
}

type GuildLong struct {
	Id string `json:"id"`
	SecurityLevel int `json:"securityLevel"`
	ReportCount int `json:"reportCount"`
	Checked bool `json:"checked"`
	Reports []Report `json:"reports"`
}

type Report struct{
	Id int `json:"id"`
	GuildId string `json:"guildId"`
	UserId string `json:"userId"`
	Text string `json:"text"`
	Date int `json:"date"`
	Closed bool`json:"closed"`
}

type UserExists struct{
	Exists bool `json:"exists"`
}

type GenUserDataPost struct{
	Avatar string `json:"avatar"`
	Username string `json:"username"`
	Discrim string `json:"discrim"`
}

type GuildreportPost struct {
	Reason string `json:"reason"`
}

type GuildreportPostSend struct {
	Reason string `json:"reason"`
	UserId string `json:"userId"`
}
