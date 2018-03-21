package main

import (
	"golang.org/x/oauth2"
	"github.com/ArgonautDevelopments/SessionManager"
	_ "github.com/ArgonautDevelopments/SessionManager/providers/memory"
	"log"
	"github.com/gin-gonic/gin"
	"net/http"
	"io/ioutil"
	"bytes"
	"io"
	"github.com/ArgonautDevelopments/SoraDashboard/src/Utility"
	"github.com/Daniele122898/Project-Aegis-Website/src/models"
	"encoding/json"
	"github.com/Daniele122898/Project-Aegis-Website/src/Utils"
	"strings"
	"github.com/Daniele122898/Project-Aegis-Website/src/config"
)

const (
	authURL      string = "https://discordapp.com/api/oauth2/authorize"
	tokenURL     string = "https://discordapp.com/api/oauth2/token"
	userEndpoint string = "https://discordapp.com/api/users/@me"
)

var (
	discordOauthConfig = &oauth2.Config{
		RedirectURL: "http://project-aegis.pw/discordCallback",
		ClientID: config.Get().ClientID,
		ClientSecret: config.Get().ClientSecret,
		Scopes: []string{"identify"},
		Endpoint: oauth2.Endpoint{
			AuthURL: authURL,
			TokenURL: tokenURL,
		},
	}
	// Some random string, random for each request
	oauthStateString = "random"

	//Sesison Manager
	globalSessions *sessionmanager.Manager

	existingUsersCache = make(map[string]bool)

)

func init(){
	gin.SetMode(gin.ReleaseMode)

	var err error
	globalSessions, err = sessionmanager.NewManager("memory", "project-aegis-sessionid", 7200)
	if err != nil {
		log.Fatal("Error initializing session: ",err)
	}
	//start garbage collection
	go globalSessions.GC()
}

func main(){
	router := gin.Default()

	//router.Use(static.Serve("/", static.LocalFile("./public", false)))
	router.Use(Serve("/", LocalFile("./public", false)))
	router.GET("/", handleHomepage)
	router.GET("/dashboard", handleDashboard)
	router.GET("/logout", handleLogout)
	router.GET("/admin", handleAdminPage)
	router.GET("/api/syncProfile", handleSyncProfile)
	router.GET("/api/genNewToken", handleGenerateToken)
	router.GET("/api/getUserInfo", handleGetUserInfo)
	router.GET("/api/getToken", handleGetToken)
	router.GET("/guild/:id", handleGuild)
	router.GET("/api/guild/:id", handleGuildInfo)
	router.GET("/api/guildlist", handleGuildList)

	//POST
	router.POST("/api/generateMarkdown", handleGenerateMarkdown)
	router.POST("/api/guildReport/:id", handleGuildReport)

	//discordlogin
	router.GET("/discordlogin",handleDiscordLogin)
	router.GET("/discordCallback", handleDiscordCallback)

	//if url is invalid
	router.NoRoute(handlePageNotFound)

	router.Run(":8300")
}

func handleDiscordLogin(c *gin.Context){
	url := discordOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(c.Writer,c.Request,url, http.StatusTemporaryRedirect)
}

func handleHomepage(c *gin.Context){
	html, err := ioutil.ReadFile("./public/index.html")
	if err!= nil {
		log.Println("COULDNT FIND LOCAL INDEX FILE")
	}
	c.Header("Content-Type", "text/html")
	c.String(200, string(html))
}

func handleDiscordCallback(c *gin.Context){
	r:= c.Request
	w:= c.Writer
	state := r.FormValue("state")
	if state != oauthStateString {
		log.Printf("Invalid oauth state! expected '%s', got '%s'\n", oauthStateString, state)
		//redirect back
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	code := r.FormValue("code")
	token, err := discordOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("Code exchange failed with '%s'\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	req, err := http.NewRequest("GET", userEndpoint, nil)
	if err != nil{
		log.Println("error doing request to discord: ", err)
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error getting response: ", err)
		return
	}
	//close the body at the end
	defer resp.Body.Close()

	buf := bytes.NewBuffer(nil)

	_, err = io.Copy(buf, resp.Body)
	if err != nil{
		log.Println("Error copy buffer: ",err)
		return
	}

	login := models.LoginUserData{}

	err = json.Unmarshal(buf.Bytes(), &login)
	if err != nil{
		log.Println("Error unmarshalling login data: ",err)
		return
	}
	//fix avatar https://cdn.discordapp.com/avatars/id/avatar .png"
	login.Avatar = "https://cdn.discordapp.com/avatars/"+login.Id+"/"+login.Avatar+".png"


	sess := globalSessions.SessionStart(w,r)
	sess.Set(Utility.LoginData, login)

	//check if user already exists
	if _, ok:= existingUsersCache[login.Id]; !ok {
		//he wasn't cached yet so make a call to the API to check if he exists
		req, err := Utils.AdminGetRequest("http://localhost:8200/api/admin/blacklist/user/exists/"+login.Id)
		if err != nil{
			//smth went wrong. Maybe API is down? dont let them login.
			log.Println("Error checking if user exists, ",err)
			c.Redirect(http.StatusTemporaryRedirect, "http://project-aegis.pw/")
			return
		}
		var exists models.UserExists
		err = json.Unmarshal(req, &exists)
		if err != nil{
			//smth went wrong. Im probably retarded af
			log.Println("Error checking parsing user exists json, ",err)
			c.Redirect(http.StatusTemporaryRedirect, "http://project-aegis.pw/")
			return
		}
		if !exists.Exists{
			userData:=models.GenUserDataPost{Avatar:login.Avatar, Username: login.Username, Discrim:login.Discriminator}
			data, err := json.Marshal(&userData)
			if err != nil{
				//failed to create user for some reason. i fucked up probably
				log.Println("Failed to marshall data for user creation, ",err)
				c.Redirect(http.StatusTemporaryRedirect, "http://project-aegis.pw/")
				return
			}
			resp, err := Utils.AdminPostRequest("http://localhost:8200/api/admin/blacklist/user/createUser/"+login.Id, data)
			if err!=nil{
				//failed to create user
				log.Println("Failed to create user, ",err)
				c.Redirect(http.StatusTemporaryRedirect, "http://project-aegis.pw/")
				return
			}
			if strings.Contains(string(resp), "error"){
				log.Println("Some error happened from user creation \n"+string(resp))
				c.Redirect(http.StatusTemporaryRedirect, "http://project-aegis.pw/")
				return
			}
		}
		//successfully created (or already existed) user so paste in cache
		existingUsersCache[login.Id] = true
	}
	//Go to Dashboard
	c.Redirect(http.StatusTemporaryRedirect,"http://project-aegis.pw/dashboard")
}


