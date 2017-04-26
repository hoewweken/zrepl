package zfs

import (
	"github.com/zrepl/zrepl/model"
	"os/exec"
	"bufio"
	"strings"
	"errors"
	"io"
	"fmt"
	"io/ioutil"
)

func InitialSend(snapshot string) (io.Reader, error) {
	return nil, nil
}

func IncrementalSend(from, to string) (io.Reader, error) {
	return nil, nil
}

func FilesystemsAtRoot(root string) (fs model.Filesystem, err error) {

	_, _ = zfsList("zroot", func(path DatasetPath) bool {
		return true
	})

	return

}

type DatasetPath []string

func (p DatasetPath) ToString() string {
	return strings.Join(p, "/")
}

func NewDatasetPath(s string) (p DatasetPath, err error) {
	// TODO validation
	return toDatasetPath(s), nil
}

func toDatasetPath(s string) DatasetPath {
	return strings.Split(s, "/")
}

type DatasetFilter func(path DatasetPath) bool

type ZFSError struct {
	Stderr  []byte
	WaitErr error
}

func (e ZFSError) Error() string {
	return fmt.Sprintf("zfs exited with error: %s", e.WaitErr.Error())
}

var ZFS_BINARY string = "zfs"

func zfsList(root string, filter DatasetFilter) (datasets []DatasetPath, err error) {

	const ZFS_LIST_FIELD_COUNT = 1

	cmd := exec.Command(ZFS_BINARY, "list", "-H", "-r",
						"-t", "filesystem,volume",
						"-o", "name",
						root)

	var stdout io.Reader
	var stderr io.Reader

	if stdout, err = cmd.StdoutPipe(); err != nil {
		return
	}

	if stderr, err = cmd.StderrPipe(); err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	s := bufio.NewScanner(stdout)
	buf := make([]byte, 1024)
	s.Buffer(buf, 0)

	datasets = make([]DatasetPath, 0)

	for s.Scan() {
		fields := strings.SplitN(s.Text(), "\t", ZFS_LIST_FIELD_COUNT)
		if len(fields) != ZFS_LIST_FIELD_COUNT {
			err = errors.New("unexpected output")
			return
		}

		dp := toDatasetPath(fields[0])

		if filter(dp) {
			datasets = append(datasets, dp)
		}
	}

	stderrOutput, err := ioutil.ReadAll(stderr)

	if waitErr:= cmd.Wait(); waitErr != nil {
		err := ZFSError{
			Stderr: stderrOutput,
			WaitErr: waitErr,
		}
		return nil, err
	}

	return

}
