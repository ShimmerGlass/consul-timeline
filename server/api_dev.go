// +build !release

package server

import (
	"net/http"
	"os"
	"path"
)

func (s *Server) serveStatic() {
	cwd, _ := os.Getwd()
	path := path.Join(cwd, path.Dir(os.Args[0]), "public")
	http.FileServer(http.Dir(path))

	s.router.ServeFiles("/web/*filepath", http.Dir(path))
}
