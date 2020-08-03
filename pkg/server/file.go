package server

import (
	"compress/flate"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/congcongke/httpfileserver/pkg/common"
	"github.com/congcongke/httpfileserver/pkg/lock"

	"github.com/gin-gonic/gin"
	"github.com/mholt/archiver"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

type LocalFileHandle struct {
	dir string
	l   *lock.Lock
}

func NewLocalFileHandle(dir string) *LocalFileHandle {
	return &LocalFileHandle{
		dir: dir,
		l:   lock.NewLock(),
	}
}

func (l *LocalFileHandle) Get(c *gin.Context) {
	filename := c.Param("filename")
	fullpath := filepath.Join(l.dir, filename)
	_, err := os.Stat(fullpath)
	if err != nil {
		logrus.Errorf("file %s may be not found: %v", fullpath, err)
		common.E(c, common.ENotFound)
		return
	}

	subj := bson.NewObjectId().Hex()

	l.l.Lock(fullpath, subj)
	l.l.Unlock(fullpath, subj, func() error {
		return nil
	})
	c.File(fullpath)
}

func (l *LocalFileHandle) Put(c *gin.Context) {
	filename := c.Param("filename")
	fullpath := filepath.Join(l.dir, filename)

	subj := bson.NewObjectId().Hex()
	fakeName := fullpath + "." + subj

	logrus.Infof(fakeName)

	done := false
	l.l.Lock(fullpath, subj)
	defer func() {
		l.l.Unlock(fullpath, subj, func() error {
			if done {
				e := os.Rename(fakeName, fullpath)
				if e != nil {
					common.E(c, e)
				}
				return e
			}
			common.R(c, nil)
			return nil
		})
	}()

	file, err := os.Create(fakeName)
	if err != nil {
		logrus.Errorf("create file failed: %v", err)
		common.E(c, common.EInternalError.WithPayload(err.Error()))
		return
	}
	defer file.Close()

	written, err := io.Copy(file, c.Request.Body)
	if err != nil {
		logrus.Errorf("write file failed: %v", err)
		common.E(c, common.EInternalError.WithPayload(err.Error))
		return
	}

	done = true

	logrus.Infof("%v bytes written to %v", written, fullpath)
}

func (l *LocalFileHandle) List(c *gin.Context) {
	infos, err := ioutil.ReadDir(l.dir)
	if err != nil {
		common.E(c, err)
		return
	}

	files := []string{}
	for _, v := range infos {
		files = append(files, v.Name())
	}

	common.R(c, files)
}

func ArchiveDir(dir string, writer io.Writer) error {
	z := archiver.TarGz{
		CompressionLevel: flate.DefaultCompression,
		Tar: &archiver.Tar{
			MkdirAll:          true,
			OverwriteExisting: true,
		},
	}

	info, err := os.Stat(dir)
	if err != nil {
		logrus.Errorf("stat file %s failed: %v", dir, err)
		return err
	}
	if !info.IsDir() {
		logrus.Errorf("%v is not a dir", dir)
		return fmt.Errorf("%v is not a dir", dir)
	}

	err = z.Create(writer)
	if err != nil {
		logrus.Errorf("create writer failed: %v", err)
		return err
	}
	defer z.Close()

	err = WalkDir(dir, func(path string, info os.FileInfo, e error) error {
		filename := filepath.Join(path, info.Name())
		internalName, walkErr := archiver.NameInArchive(info, filename, filename)
		if walkErr != nil {
			logrus.Errorf("name in archive failed: %v", walkErr)
			return walkErr
		}

		file, walkErr := os.Open(filename)
		if walkErr != nil {
			logrus.Errorf("open file %v failed: %v", filename, walkErr)
			return walkErr
		}
		defer file.Close()

		walkErr = z.Write(archiver.File{
			FileInfo: archiver.FileInfo{
				FileInfo:   info,
				CustomName: internalName,
			},
			ReadCloser: file,
		})
		if walkErr != nil {
			logrus.Errorf("write file to archive %v failed: %v", filename, walkErr)
		}
		return walkErr
	})

	return nil
}

func WalkDir(dir string, walkFunc filepath.WalkFunc) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		logrus.Errorf("read dir failed: %v", err)
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			err = WalkDir(filepath.Join(dir, file.Name()), walkFunc)
			if err != nil {
				logrus.Errorf("walk dir failed: %v", err)
			}
			continue
		}

		if err = walkFunc(dir, file, err); err != nil {
			logrus.Errorf("walk failed for file %v: %v", file.Name(), err)
			return err
		}
	}

	return nil
}
