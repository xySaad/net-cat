package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/jroimartin/gocui"
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func writetoconn(g *gocui.Gui, v *gocui.View) error {
	in, err := g.View("v1")
	defer in.Clear()
	if err != nil {
		return err
	}
	rr := make([]byte, 1024)
	for {
		n, err := in.Read(rr)
		if n == 0 {
			return nil
		}
		if err != nil {
			log.Fatal(err)
		}
		if rr[n-1] == '\n' {
			in.Write([]byte("\r"))
			if err := in.SetCursor(0, 0); err != nil {
				return err
			}
			fmt.Fprint((*c), string(rr[:n-1])+"\n")
			if !strings.Contains(vv.Title, "[ENTER GROUP NAME]:") && !strings.Contains(vv.Title, "[ENTER YOUR NAME]:") {
				g.Update(func(g *gocui.Gui) error {
					ch.Lock()
					fmt.Fprint(ch.cc, "\n"+vv.Title+string(rr[:n-1]))
					ch.Unlock()
					return nil
				})
			} else if strings.Contains(vv.Title, "[ENTER GROUP NAME]:") {
				hh.Title = string(rr[:n-1]) + " chat:"
			}

			return nil
		}
	}
}

var (
	c  *net.Conn
	vv *gocui.View
	ch = typ{}
	hh *gocui.View
)

type typ struct {
	sync.Mutex
	cc *gocui.View
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("v1", 0, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		vv = v
		v.Editable = true
		v.Wrap = true
		if _, err = setCurrentViewOnTop(g, "v1"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("v2", 0, 0, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		hh = v
		v.Title = "chat"
		v.Wrap = true
		v.Autoscroll = true
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func readconn(conn *net.Conn, g *gocui.Gui) {
	buffer := make([]byte, 4096)
	chat, err := g.View("v2")
	if err != nil {
		log.Fatal(err)
	}
	ch.cc = chat
	n, err := (*conn).Read(buffer)
	if err != nil {
		if err == io.EOF {
			log.Println("Connection closed by the server")
			os.Exit(1)
		}
		log.Printf("Error reading from connection: %v", err)
		os.Exit(1)
	}
	dd := replace(Split(string(buffer[:n])))
	ee := strings.Split(dd, "\n")
	if len(ee[len(ee)-1]) >= 18 {
		vv.Title = ee[len(ee)-1]
		ee = ee[:len(ee)-1]
	}

	g.Update(func(g *gocui.Gui) error {
		ch.Lock()
		fmt.Fprint(ch.cc, strings.Join(ee, "\n"))
		ch.Unlock()
		return nil
	})
	readconn(conn, g)
}

func Sort(s []int) []int {
	for i := 0; i < len(s)-1; i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] > s[j] {
				s[i], s[j] = s[j], s[i]
			}
		}
	}
	return s
}

func Split(s string) []string {
	e := []string{}
	mapp := make(map[int]int)
	i := 0
	for i < len(s) {
		slice := []int{}

		slice = append(slice, strings.Index(s[i:], "\x1b[2J\x1b[3J\x1b[H"))
		mapp[strings.Index(s[i:], "\x1b[2J\x1b[3J\x1b[H")] = len("\x1b[2J\x1b[3J\x1b[H")
		slice = append(slice, strings.Index(s[i:], "\x1b[F\x1b[2K"))
		mapp[strings.Index(s[i:], "\x1b[F\x1b[2K")] = len("\x1b[F\x1b[2K")
		slice = append(slice, strings.Index(s[i:], "\x1b[s\n\x1b[F\x1b[2K"))
		mapp[strings.Index(s[i:], "\x1b[s\n\x1b[F\x1b[2K")] = len("\x1b[s\n\x1b[F\x1b[2K")
		slice = append(slice, strings.Index(s[i:], "\x1b[G\x1b[2K"))
		mapp[strings.Index(s[i:], "\x1b[G\x1b[2K")] = len("\x1b[G\x1b[2K")
		slice = append(slice, strings.Index(s[i:], "\x1b[u\x1b[B"))
		mapp[strings.Index(s[i:], "\x1b[u\x1b[B")] = len("\x1b[u\x1b[B")
		slice = append(slice, strings.Index(s[i:], "\x1b[38;2;0;184;30m"))
		mapp[strings.Index(s[i:], "\x1b[38;2;0;184;30m")] = len("\x1b[38;2;0;184;30m")
		slice = append(slice, strings.Index(s[i:], "\x1b[38;2;255;0;0m"))
		mapp[strings.Index(s[i:], "\x1b[38;2;255;0;0m")] = len("\x1b[38;2;255;0;0m")
		slice = append(slice, strings.Index(s[i:], "\x1b]0;"))
		mapp[strings.Index(s[i:], "\x1b]0;")] = len("\x1b]0;")
		slice = append(slice, strings.Index(s[i:], "\x1b[0m"))
		mapp[strings.Index(s[i:], "\x1b[0m")] = len("\x1b[0m")
		slice = Sort(slice)
		max := slice[len(slice)-1]
		if max < 0 {
			e = append(e, s[i:])
			break
		}
		for ii := 0; ii < len(slice); ii++ {
			slice[ii] += i
		}
		ccc := i
		for _, v := range slice {
			if i == len(s)-1 {
				break
			}
			if v == i-1 || v < i {
				continue
			}
			if !(i < 0 || v > len(s)) {
				e = append(e, s[i:v])
			}
			if v+mapp[v-i] < len(s) {
				e = append(e, s[v:v+mapp[v-ccc]])
			}
			i = v + mapp[v-ccc]
		}

	}
	return e
}

func replace(s []string) string {
	t := ""
	for _, v := range s {
		switch v {
		case "\x1b[2J\x1b[3J\x1b[H", "\x1b[u\x1b[B", "\x1b[38;2;0;184;30m", "\x1b[38;2;255;0;0m", "\x1b[0m":
		case "\x1b[F\x1b[2K", "\x1b[G\x1b[2K":
			t = "\r"
		case "\x1b[s\n\x1b[F\x1b[2K":
			t += "\n"
		case "\x1b]0;":
			t += "[ENTER GROUP NAME]:"
		default:
			t += v
		}
	}
	return t
}

func Client(protocol, adress string) {
	conn, err := net.Dial(protocol, adress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c = &conn
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorBlue
	g.BgColor = gocui.AttrUnderline
	g.SetManagerFunc(layout)
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, writetoconn); err != nil {
		log.Panicln(err)
	}
	go readconn(&conn, g)
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
