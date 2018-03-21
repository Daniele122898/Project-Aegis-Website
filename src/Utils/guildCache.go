package Utils

import (
	"time"
	"github.com/Daniele122898/Project-Aegis-Website/src/models"
	"github.com/shurcooL/github_flavored_markdown"
	"github.com/microcosm-cc/bluemonday"
	"regexp"
	"strings"
	"log"
)

const MaxAge int64 = 120

type Guild struct{
	Guild *models.GuildWeb
	Created int64
}


var(
	GuildCache = make(map[string]*Guild)
)

func AddToCache(guild *models.GuildWeb) {
	//Do markdown
	for i := range guild.Reports{
		str := strings.Replace(guild.Reports[i].Text, "-- EDIT -- ", "\n -- EDIT -- \n", -1)
		log.Println(str)
		unsafe := github_flavored_markdown.Markdown([]byte(str))
		p := bluemonday.UGCPolicy()
		p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
		html := p.SanitizeBytes(unsafe)
		guild.Reports[i].Text= string(html)
	}

	g := Guild{Guild:guild, Created:time.Now().Unix()}
	GuildCache[guild.Id] = &g
}

func GetGuild(id string) *models.GuildWeb{
	g, ok := GuildCache[id]
	if !ok {
		return nil
	}
	//cache is too old
	if time.Now().Unix() - g.Created > MaxAge {
		//remove from map
		GuildCache[id] = nil
		return nil
	}
	return g.Guild
}


