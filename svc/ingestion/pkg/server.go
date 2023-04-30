package pkg

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	ingestionv1 "github.com/alexanderjophus/kie/gen/ingestion/v1"
	"github.com/bufbuild/connect-go"
)

const (
	RepoName = "statsapi"
	baseURL  = "https://statsapi.web.nhl.com"
)

// A Server is a server
type Server struct {
	store FileStorer
	c     *http.Client
}

// FileStorer represents the ability to store files
type FileStorer interface {
	storeFile(r io.Reader, path string) (err error)
}

// NewServer creates a new server struct, initialised with the routing set
func NewServer(fs FileStorer) *Server {
	return &Server{
		store: fs,
		c: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// addFile handles /addFile endpoint
func (s *Server) AddFile(ctx context.Context, req *connect.Request[ingestionv1.AddFileRequest]) (*connect.Response[ingestionv1.AddFileResponse], error) {
	err := s.linkStorer(req.Msg.Link)
	if err != nil {
		return nil, fmt.Errorf("unable to store file: %w", err)
	}
	return &connect.Response[ingestionv1.AddFileResponse]{}, nil
}

func (s *Server) linkStorer(link string) error {
	r, err := s.c.Get(baseURL + link)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", r.Status)
	}

	err = s.store.storeFile(r.Body, link)
	if err != nil {
		return err
	}

	return nil
}
