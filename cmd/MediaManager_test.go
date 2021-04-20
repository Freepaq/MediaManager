package main

import (
	"testing"

	"github.com/Freepaq/MediaManagement/pkg/MediaUtils"
)

func TestPhoto(t *testing.T) {
	fname := "testfiles/phototest.jpg"
	f, err := MediaUtils.GetMeta(fname)
	if nil != err {
		t.Failed()
		t.Errorf(err.Error())
	}
	if f.CreationDate.Format("2006-01-02 15:04") != "2021-04-17 17:11" && f.MetaOrigin != MediaUtils.METAORIGINFILE {
		t.Failed()
		t.Errorf("Failed to read Meta from photo (" + f.CreationDate.Format("2006-01-02 15:04") + ") " + fname)
	}

}

func TestVideo(t *testing.T) {
	/*fname := "testfiles/videotest.mp4"
	f, err := MediaUtils.GetMeta(fname)
	if nil != err {
		t.Failed()
		t.Errorf(err.Error())
	}
	if f.CreationDate.Format("2006-01-02 15:04") != "2017-07-07 10:16" {
		t.Failed()
		t.Errorf("Failed to read Meta from video (" + f.CreationDate.Format("2006-01-02 15:04") + ") " + fname)
	}*/

}
