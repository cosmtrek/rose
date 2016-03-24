package main

import (
	"net"
	"strconv"
	"sync"
)

type Global struct {
	sync.RWMutex
	OnlineUsers map[int]*net.Conn
}

var (
	global = Global{
		OnlineUsers: make(map[int]*net.Conn),
	}
)

func (g *Global) updateOnlineUsers(conn *net.Conn) {
	debug.Println("Updating online users...")
	for k, v := range g.OnlineUsers {
		if *v == *conn {
			g.deleteOnlineUser(k)
		}
	}
}

func (g *Global) getOnlineUser(uid int) (*net.Conn, bool) {
	g.RLock()
	c, ok := global.OnlineUsers[uid]
	g.RUnlock()
	return c, ok
}

func (g *Global) deleteOnlineUser(uid int) {
	g.Lock()
	debug.Println("Delete user " + strconv.Itoa(uid))
	delete(g.OnlineUsers, uid)
	g.Unlock()
}

func (g *Global) addOnlineUser(uid int, conn *net.Conn) {
	g.Lock()
	g.OnlineUsers[uid] = conn
	g.Unlock()
}
