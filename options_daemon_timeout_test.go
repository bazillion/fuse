// Test for adjustable timeout between a FUSE request and the daemon's response.
//
// +build darwin freebsd

package fuse_test

import (
	"os"
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/presotto/fuse"
	"github.com/presotto/fuse/fs"
	"github.com/presotto/fuse/fs/fstestutil"
	"golang.org/x/net/context"
)

type slowCreaterDir struct {
	fstestutil.Dir
}

var _ fs.NodeCreater = slowCreaterDir{}

func (c slowCreaterDir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	time.Sleep(10*time.Second)
	// pick a really distinct error, to identify it later
	return nil, nil, fuse.Errno(syscall.ENAMETOOLONG)
}

func TestDaemonTimeout(t *testing.T) {
	if runtime.GOOS != "darwin" && runtime.GOOS != "freebsd" {
		return
	}
	t.Parallel()
	mnt, err := fstestutil.MountedT(t,
		fstestutil.SimpleFS{slowCreaterDir{}},
		nil,
		fuse.DaemonTimeout("2"),
	)

	// This should fail by the kernel timing out the request.
	f, err := os.Create(mnt.Dir + "/child")
	if err == nil {
		f.Close()
		t.Fatal("expected an error")
	}
	perr, ok := err.(*os.PathError)
	if !ok {
		t.Fatalf("expected PathError, got %T: %v", err, err)
	}
	if perr.Err == syscall.ENAMETOOLONG {
		t.Fatalf("expected other than ENAMETOOLONG, got %T: %v", err, err)
	}
}