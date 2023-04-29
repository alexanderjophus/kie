package pkg

import (
	"io"
	"strings"

	"github.com/pachyderm/pachyderm/v2/src/client"
	"github.com/pachyderm/pachyderm/v2/src/pfs"
)

type PachdRepo struct {
	Client *client.APIClient
}

func (c PachdRepo) storeFile(r io.Reader, path string) (err error) {
	// Start a commit in our "livefeed" data repo on the "master" branch.
	commit, err := c.Client.StartProjectCommit(pfs.DefaultProjectName, RepoName, "master")
	if err != nil {
		return err
	}
	defer func() {
		err = c.Client.FinishProjectCommit(pfs.DefaultProjectName, RepoName, "master", commit.ID)
	}()

	path = strings.Replace(path, "?", "/", -1)
	path = strings.Replace(path, "=", "_", -1)
	// Put a file containing the respective project name.
	if err := c.Client.PutFile(commit, path+".raw", r); err != nil {
		return err
	}

	return nil
}
