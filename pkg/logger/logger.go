package logger

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
	"time"
)

func now() string {
	c := color.New(color.FgHiBlack)
	return c.Sprint(time.Now().Format("15:04:05"))
}

func number(n any) string {
	c := color.New(color.FgYellow, color.Bold)
	return c.Sprintf("%v", n)
}

func numberList(n []uint) string {
	var s string
	for _, v := range n {
		s += number(v) + ", "
	}
	return strings.TrimSuffix(s, ", ")
}

func ip(ip string) string {
	var str string

	c := color.New(color.FgGreen)
	data := strings.Split(ip, ".")
	for _, v := range data {
		str += c.Sprintf("%s", v) + "."
	}
	return strings.TrimSuffix(str, ".")
}

func user(userID any) string {
	return fmt.Sprintf(" : %s", number(userID))
}

func ws(wsID string, userID uint) string {
	c := color.New(color.BgHiMagenta, color.FgHiWhite, color.Bold)
	id := color.New(color.FgHiBlue, color.Italic)
	var u string
	if userID > 0 {
		u = user(userID)
	}
	return fmt.Sprintf("%s ~ %s [%s%s] -", c.Sprintf(" WS "), now(), id.Sprintf(wsID), u)
}

func WSNewConnection(wsID string, ipAddr string, userID uint) {
	fmt.Printf("%s New connection (%s)\n", ws(wsID, userID), ip(ipAddr))
}

func WSNewEvent(wsID, event string, userID uint) {
	fmt.Printf("%s New event (%s)\n", ws(wsID, userID), event)
}

func WSDisconnect(wsID string, userID uint) {
	c := color.New(color.FgRed)
	fmt.Printf("%s %s\n", ws(wsID, userID), c.Sprint("Connection closed"))
}

func WSSend(wsID string, event string, userID uint, userIDs []uint) {
	fmt.Printf("%s Message (%s) Sent to: [%s]\n", ws(wsID, userID), event, numberList(userIDs))
}

func WSPoolSize(wsID string, size int, ids []uint) {
	fmt.Printf("%s Pool size: %s: [%s]\n", ws(wsID, 0), number(size), numberList(ids))
}
