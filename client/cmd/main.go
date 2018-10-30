package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/go-oauth2/oauth2/utils/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/zcong1993/oauth2-store-mysql/client"
	"strconv"
	"strings"
	"time"
)

func newToken(name, appName string) string {
	buf := bytes.NewBufferString(name)

	buf.WriteString(appName)
	buf.WriteString(strconv.FormatInt(time.Now().UnixNano(), 10))

	token := base64.URLEncoding.EncodeToString(uuid.NewMD5(uuid.Must(uuid.NewRandom()), buf.Bytes()).Bytes())
	token = strings.ToUpper(strings.TrimRight(token, "="))

	return token
}

func main() {
	appName := flag.String("app-name", "", "set the app name. ")
	domain := flag.String("domain", "", "set your domain. ")
	mysql := flag.String("mysql", "", "set mysql dsn. ")
	flag.Parse()

	if *appName == "" {
		panic("app name is required. ")
	}
	if *domain == "" {
		panic("domain is required. ")
	}
	if *mysql == "" {
		panic("mysql dsn is required. ")
	}

	c := &client.Client{
		AppName: *appName,
		Domain:  *domain,
		UID:     newToken("uid", *appName),
		Secret:  newToken("secret", *appName),
	}

	db, err := gorm.Open("mysql", *mysql)

	if err != nil {
		panic(err)
	}

	err = db.Save(c).Error
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println("App Name: ", c.AppName)
	fmt.Println("Client ID: ", c.UID)
	fmt.Println("Clieng Secret: ", c.Secret)
}
