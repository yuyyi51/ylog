package ylog

import (
	"fmt"
	"os"
	"path"
	"time"
)

type DateInfo struct {
	Year  int
	Month time.Month
	Day   int
	Hour  int
}

func (info *DateInfo) Match(now time.Time) bool {
	return now.Hour() == info.Hour &&
		now.Day() == info.Day &&
		now.Month() == info.Month &&
		now.Year() == info.Year
}

func (info *DateInfo) MatchNow() bool {
	return info.Match(time.Now())
}

func (info *DateInfo) Update(new time.Time) {
	info.Year = new.Year()
	info.Month = new.Month()
	info.Day = new.Day()
	info.Hour = new.Hour()
}

func (info *DateInfo) UpdateNow() {
	info.Update(time.Now())
}

type fileWriter struct {
	path            string
	prefix          string
	currentFile     *os.File
	currentDateInfo DateInfo
}

func newFileWriter(path, prefix string) (*fileWriter, error) {
	return &fileWriter{
		path:   path,
		prefix: prefix,
	}, nil
}

func (w *fileWriter) Write(b []byte) (int, error) {
	if !w.currentDateInfo.MatchNow() {
		err := os.MkdirAll(w.path, 0755)
		if err != nil {
			return 0, err
		}
		w.currentDateInfo.UpdateNow()
		filePath := path.Join(w.path, w.createFileName())
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
		if err != nil {
			return 0, nil
		}
		w.currentFile = file
	}
	return w.currentFile.Write(b)
}

func (w *fileWriter) createFileName() string {
	return fmt.Sprintf("%s.%d-%02d-%02d_%02d",
		w.prefix, w.currentDateInfo.Year, w.currentDateInfo.Month, w.currentDateInfo.Day, w.currentDateInfo.Hour)
}
