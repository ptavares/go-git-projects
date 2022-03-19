package git

import (
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type AuthType int64

const (
	Basic AuthType = iota
	SSH
)

// Define Git Authentication
type Auth struct {
	AuthType  AuthType
	BasicAuth *http.BasicAuth
	SSHAuth   *ssh.PublicKeys
}

// Define all Gitlab group needed information
type Group struct {
	Path     string
	GID      int
	Projects []*Project
}

// String for Group
func (g *Group) String() string {
	s := new(strings.Builder)
	s.WriteString("{ \"Path\": \"")
	s.WriteString(g.Path)
	s.WriteString("\", \"GID\" : \"")
	s.WriteString(strconv.Itoa(g.GID))
	s.WriteString("\", \"Projects\" : \"[")
	for _, p := range g.Projects {
		s.WriteString(p.String())
	}
	s.WriteString("]\"}")
	return s.String()
}

type Project struct {
	Name          string
	PID           int
	SSHURLToRepo  string
	HTTPURLToRepo string
}

// String for Project
func (p *Project) String() string {
	s := new(strings.Builder)
	s.WriteString("{ \"name\": \"")
	s.WriteString(p.Name)
	s.WriteString("\", \"pid\" : \"")
	s.WriteString(strconv.Itoa(p.PID))
	s.WriteString("\", \"SSHURLToRepo\" : \"")
	s.WriteString(p.SSHURLToRepo)
	s.WriteString("\", \"HTTPURLToRepo\" : \"")
	s.WriteString(p.HTTPURLToRepo)
	s.WriteString("\"}")
	return s.String()
}
