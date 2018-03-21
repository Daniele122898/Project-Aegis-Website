package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ArgonautDevelopments/SoraDashboard/src/Utility"
	"log"
	"github.com/Daniele122898/Project-Aegis-Website/src/models"
	"net/http"
	"io/ioutil"
	"github.com/ArgonautDevelopments/SessionManager"
	"bytes"
	"io"
	"github.com/shurcooL/github_flavored_markdown"
	"github.com/microcosm-cc/bluemonday"
	"regexp"
	"github.com/Daniele122898/Project-Aegis-Website/src/Utils"
	"encoding/json"
	"strconv"
	"github.com/Daniele122898/Project-Aegis-Website/src/config"
)

func handleAdminPage(c *gin.Context){
	_, user, ok := CheckPermission(c, true)
	if !ok {
		return
	}
	if user.Id != config.Get().OwnerId{
		handlePageNotFound(c)
		return
	}
	html, err := ioutil.ReadFile("./public/admin.html")
	if err!= nil {
		log.Println("COULDNT FIND LOCAL admin FILE")
		c.Redirect(http.StatusTemporaryRedirect, "http://project-aegis.pw/")
		return
	}
	c.Header("Content-Type", "text/html")
	c.String(200, string(html))
}

func handlePageNotFound(c *gin.Context){
	//load not found page
	html, err := ioutil.ReadFile("./public/notfound.html")
	if err!= nil {
		log.Println("COULDNT FIND LOCAL notFound FILE")
		c.Redirect(http.StatusTemporaryRedirect, "http://project-aegis.pw/")
		return
	}
	c.Header("Content-Type", "text/html")
	c.String(200, string(html))
}

func handleGuildList(c*gin.Context){
	_, _, ok := CheckPermission(c, true)
	if !ok {
		return
	}

	resp,ok, err := Utils.AdminGetRequestLong("http://localhost:8200/api/admin/blacklist/guild/list", http.StatusOK)
	if !ok || err != nil {
		handlePageNotFound(c)
		return
	}
	c.Writer.Write(resp)
}

func handleGuild(c *gin.Context){
	_, user, ok := CheckPermission(c, true)
	if !ok {
		return
	}

	//Ratelimit
	limited := Utils.CheckIfRatelimited(user.Id)
	if limited {
		c.JSON(http.StatusTooManyRequests, gin.H{"error":"Too many requests. You are being ratelimited."})
		return
	}
	//invoke ratelimit
	Utils.InvokeRatelimit(user.Id)

	idS := c.Params.ByName("id")
	ok= Utils.ValidId(idS)
	if !ok {
		handlePageNotFound(c)
		return
	}

	resp,ok, err := Utils.AdminGetRequestLong("http://localhost:8200/api/admin/blacklist/guild/infolong/"+idS,http.StatusOK)
	if !ok || err != nil {
		handlePageNotFound(c)
		return
	}

	var guild models.GuildWeb
	err = json.Unmarshal(resp, &guild)
	if err!=nil{
		handlePageNotFound(c)
		return
	}
	//cache guild
	Utils.AddToCache(&guild)
	//load normal page
	html, err := ioutil.ReadFile("./public/guild.html")
	if err!= nil {
		log.Println("COULDNT FIND LOCAL guild FILE")
		c.Redirect(http.StatusTemporaryRedirect, "http://project-aegis.pw/")
		return
	}
	c.Header("Content-Type", "text/html")
	c.String(200, string(html))
}

func handleGuildInfo(c *gin.Context){
	_, _, ok := CheckPermission(c, true)
	if !ok {
		return
	}

	idS := c.Params.ByName("id")

	//check if guild is in cache
	g := Utils.GetGuild(idS)
	if g!= nil{
		//write json
		json.NewEncoder(c.Writer).Encode(*g)
		return
	}

	resp,ok, err := Utils.AdminGetRequestLong("http://localhost:8200/api/admin/blacklist/guild/infolong/"+idS,http.StatusOK)
	if !ok || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"Failed to get info about guild..."})
		return
	}

	var guild models.GuildWeb
	err = json.Unmarshal(resp, &guild)
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Failed to get info about guild..."})
		return
	}
	//add to cache again
	Utils.AddToCache(&guild)
	//write json
	json.NewEncoder(c.Writer).Encode(guild)
	return
}

func handleGetToken(c *gin.Context){
	_, user, ok := CheckPermission(c, true)
	if !ok{
		return
	}

	resp,ok, err := Utils.AdminGetRequestLong("http://localhost:8200/api/admin/blacklist/user/getToken/"+user.Id,http.StatusOK)
	if !ok || err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Something went wrong"})
		return
	}
	c.Writer.Write(resp)
}

func handleGetUserInfo(c *gin.Context){
	_, user, ok := CheckPermission(c, true)
	if !ok{
		return
	}
	//send local data
	json.NewEncoder(c.Writer).Encode(user)
}

func handleGenerateToken(c *gin.Context){
	_, user, ok := CheckPermission(c, true)
	if !ok{
		return
	}

	//Ratelimit
	limited := Utils.CheckIfRatelimited(user.Id)
	if limited {
		c.JSON(http.StatusTooManyRequests, gin.H{"error":"Too many requests. You are being ratelimited."})
		return
	}
	//invoke ratelimit
	Utils.InvokeRatelimit(user.Id)

	//request new token
	resp,ok, err := Utils.AdminGetRequestLong("http://localhost:8200/api/admin/blacklist/user/requestToken/"+user.Id,http.StatusOK)
	if !ok || err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Something went wrong"})
		return
	}
	//send data
	c.Writer.Write(resp)
}

func handleSyncProfile(c *gin.Context){
	_, user, ok := CheckPermission(c, true)
	if !ok {
		return
	}
	//Ratelimit
	limited := Utils.CheckIfRatelimited(user.Id)
	if limited {
		c.JSON(http.StatusTooManyRequests, gin.H{"error":"Too many requests. You are being ratelimited."})
		return
	}
	//invoke ratelimit
	Utils.InvokeRatelimit(user.Id)
	//Sync Data with DB!
	userToSend := models.GenUserDataPost{Username:user.Username,Avatar:user.Avatar, Discrim:user.Discriminator}
	sendData, err := json.Marshal(userToSend)
	if err != nil{
		log.Println("Failed to marshall user sync data, ",err)
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Something went wrong"})
		return
	}

	//make post request
	_ ,ok, err = Utils.AdminPostRequestLong("http://localhost:8200/api/admin/blacklist/user/syncProfile/"+user.Id, sendData, http.StatusOK)

	if !ok || err!=nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status":"Account successfully synchronized"})
}

func handleGuildReport (c *gin.Context){
	_, user, ok := CheckPermission(c, true)
	if !ok{
		return
	}

	//Ratelimit
	limited := Utils.CheckIfRatelimited(user.Id)
	if limited {
		c.JSON(http.StatusTooManyRequests, gin.H{"error":"Too many requests. You are being ratelimited."})
		return
	}
	//invoke ratelimit
	Utils.InvokeRatelimit(user.Id)

	//first check the fucking ID
	gid, ok := Utils.GetGuildId(&c.Params)
	if !ok{
		log.Println("Bad guild Id")
		c.JSON(http.StatusBadRequest, gin.H{"error":"Bad Guild ID"})
		return
	}

	ok = Utils.ValidIdInt(gid)
	if !ok{
		log.Println("Bad guild Id 2")
		c.JSON(http.StatusBadRequest, gin.H{"error":"Bad Guild ID"})
		return
	}
	//ID SHOULD be fine so lets continue

	defer c.Request.Body.Close()

	buf := bytes.NewBuffer(nil)

	_, err := io.Copy(buf, c.Request.Body)
	if err != nil{
		log.Println("Invalid Json in guildreport")
		c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid JSON: "+err.Error()})
		return
	}
	var post models.GuildreportPost
	err = json.Unmarshal(buf.Bytes(), &post)
	if err != nil{
		log.Println("Invalid Json in guildreport")
		c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid JSON: "+err.Error()})
		return
	}
	if len(post.Reason) > 1000{
		c.JSON(http.StatusBadRequest, gin.H{"error":"You reached the charlimit of 1000!"})
		return
	}
	if len(post.Reason)< 50{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Add a valid reason."})
		return
	}
	//do the fucking post
	send := models.GuildreportPostSend{Reason:post.Reason, UserId:user.Id}
	sendData, err := json.Marshal(send)
	if err!=nil{
		log.Println("failed to marshall send data, ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Something went wrong"})
		return
	}
	_ ,ok, err = Utils.AdminPostRequestLong("http://localhost:8200/api/admin/blacklist/guild/report/"+strconv.FormatInt(gid, 10), sendData, http.StatusOK)
	if !ok || err!=nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status":"Report has been successfully submitted"})
}

func handleDashboard(c *gin.Context){
	_, _, ok := CheckPermission(c, false)
	if !ok {
		return
	}
	html, err := ioutil.ReadFile("./public/dashboard.html")
	if err!= nil {
		log.Println("COULDNT FIND LOCAL dashboard FILE")
		c.Redirect(http.StatusTemporaryRedirect, "http://project-aegis.pw/")
		return
	}
	c.Header("Content-Type", "text/html")
	c.String(200, string(html))
}

func handleGenerateMarkdown(c *gin.Context){
	_, user, ok := CheckPermission(c, false)
	if !ok {
		return
	}

	//Ratelimit
	limited := Utils.CheckIfRatelimited(user.Id)
	if limited {
		c.JSON(http.StatusTooManyRequests, gin.H{"error":"Too many requests. You are being ratelimited."})
		return
	}
	//invoke ratelimit
	Utils.InvokeRatelimit(user.Id)

	buf := bytes.NewBuffer(nil)

	_, err := io.Copy(buf, c.Request.Body)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Failed to read markdown"})
		return
	}
	//Generate markdown
	unsafe := github_flavored_markdown.Markdown(buf.Bytes())

	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	html := p.SanitizeBytes(unsafe)

	c.Writer.Write(html)
}

func handleLogout(c *gin.Context){
	//Destroy Session
	globalSessions.SessionDestroy(c.Writer,c.Request)
	//redirect out
	log.Println("User Logged out")
	c.Redirect(http.StatusTemporaryRedirect, "http://project-aegis.pw/")
}

func CheckPermission(c *gin.Context, convert bool) (*sessionmanager.Session, *models.LoginUserData, bool){
	w := c.Writer
	r := c.Request
	sess := globalSessions.SessionStart(w,r)

	ldI := sess.Get(Utility.LoginData)

	if ldI == nil{
		//no permission to be here
		log.Println("User tried to access Page without permission!")
		c.Redirect(http.StatusTemporaryRedirect, "http://project-aegis.pw/")
		return &sess, nil, false
	}
	if !convert {
		return &sess,nil, true
	} else {
		ld, ok := ldI.(models.LoginUserData)
		if !ok{
			//failed transfo
			log.Println("Error parsing Interface to loginData: ",ldI)
			return &sess,nil, true
		}
		return &sess,&ld, true
	}
}
