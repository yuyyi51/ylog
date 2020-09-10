package ylog

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"
)

type logWriter struct {
	writerMux   *sync.Mutex
	writer      io.Writer
	buffer      *bytes.Buffer
	objectQueue chan *logObject
}

func newLogWriter(path, prefix string) (*logWriter, error) {
	fileWriter, err := newFileWriter(path, prefix)
	if err != nil {
		return nil, err
	}
	return &logWriter{
		writer:      fileWriter,
		buffer:      &bytes.Buffer{},
		objectQueue: make(chan *logObject, 10000),
		writerMux:   new(sync.Mutex),
	}, nil
}

func (l *logWriter) startRun() {
	go l.guardRun()
}

func (l *logWriter) guardRun() {
	defer func() {
		err := recover()
		if err != nil {
			// exit abnormal, restart
			go l.guardRun()
		}
	}()
	l.cycle()
}

func (l *logWriter) cycle() {
	ticker := time.NewTicker(time.Millisecond * 100)
	for {
		select {
		case obj := <-l.objectQueue:
			l.writerMux.Lock()
			_, err := l.writeObjectToBuffer(obj)
			if err != nil {
				fmt.Printf("logWriter write to buffer error: %v", err)
			}
			if l.buffer.Len() > 4*1024*1024 /* 4MB */ {
				_, err := l.writeBufferToWriter()
				if err != nil {
					fmt.Printf("logWriter write to writer error: %v", err)
				}
			}
			l.writerMux.Unlock()
		case <-ticker.C:
			l.writerMux.Lock()
			if l.buffer.Len() != 0 {
				_, err := l.writeBufferToWriter()
				if err != nil {
					fmt.Printf("logWriter write to writer error: %v", err)
				}
			}
			l.writerMux.Unlock()
		}

	}
}

func formatLog(obj *logObject) string {
	return fmt.Sprintf(fmt.Sprintf("%s %s %s:%d %s\n", obj.logTime.Format(time.RFC3339Nano), LogLevelToString(obj.level), obj.file, obj.line, obj.format), obj.args...)
}

func (l *logWriter) writeObjectToBuffer(obj *logObject) (int, error) {
	return l.buffer.WriteString(formatLog(obj))
}

func (l *logWriter) writeBufferToWriter() (int, error) {
	defer l.buffer.Reset()
	return l.writer.Write(l.buffer.Bytes())
}

func (l *logWriter) forceWriteObject(obj *logObject) {
	// To make sure panic and fatal log was written before process exit
	l.writerMux.Lock()
	_, _ = l.writeObjectToBuffer(obj)
	_, _ = l.writeBufferToWriter()
	l.writerMux.Unlock()
}
