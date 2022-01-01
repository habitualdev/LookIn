package utils

import (
	"github.com/emersion/go-message/mail"
	"net"
	"strconv"
)

type EmailServer struct {
	Addr  net.IP
	Port  int
	Proto string
}

type CacheEntry struct {
	Header     mail.Header
	Body       [][]byte
	Attachment Attachment
}

type Attachment struct {
	Exists    bool
	Filenames []string
}

func JoinAddr(addr net.IP, port int) string {
	return addr.String() + ":" + strconv.Itoa(port)
}
