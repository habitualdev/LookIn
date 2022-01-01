package main

import (
	"LookIn/imapd"
	"LookIn/ui"
	"LookIn/utils"
	"fmt"
	_ "github.com/emersion/go-message/charset"
	"github.com/muesli/termenv"
	"github.com/patrickmn/go-cache"
	"net"
	"strconv"
)


var banner =`
******************************************
*                 LookIn                 *
*                                        *
*  An IMAP client for the cli inclined   *
******************************************`

func main() {

	var address string
	var proto string
	var port string
	var username string
	var password string

	termenv.ClearScreen()

	fmt.Println(banner)
	fmt.Print("IMAP server address: ")
	fmt.Scanln(&address)
	fmt.Print("IMAP server port: ")
	fmt.Scanln(&port)
	fmt.Print("Connection protocol (tls, startls, plain): ")
	fmt.Scanln(&proto)
	fmt.Print("Username: ")
	fmt.Scanln(&username)
	fmt.Print("Password: ")
	fmt.Scanln(&password)

	portInt, _ := strconv.Atoi(port)

	info := utils.EmailServer{Addr: net.ParseIP(address), Port: portInt, Proto: proto}

	ca := cache.New(cache.NoExpiration, cache.NoExpiration)

	imapd.LoginRetrieve(ca, info, username, password)

	ui.StartUi(ca)

}
