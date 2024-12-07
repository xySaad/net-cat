package modules

import (
	"fmt"
	"io"
	"time"
)

type logger struct {
	writer io.Writer
}

func (s *Server) Info(v ...any) {
	s.log("INFO", v...)
}
func (s *Server) Warn(v ...any) {
	s.log("WARN", v...)
}
func (s *Server) Error(v ...any) {
	s.log("ERROR", v...)
}

func (s *Server) log(logLevel string, v ...any) {
	_, err := fmt.Fprint(s.logger.writer, fmt.Sprintln(time.Now().Format(time.DateTime), logLevel, v))
	if err != nil {
		fmt.Println(err)
	}
}
