package log

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sanbsy/goby/bufferpool"
)

type (
	RotateWriter struct {
		// 日志文件名
		filename string

		// 最大字节数
		maxSize int64

		// 日志备份路径
		backupPath string

		// 当前日志文件大小
		curSize int64

		// 备份文件时间后缀格式
		suffixFormat string

		// 备份文件后缀
		timeSuffix string
		index      int

		// 日志文件句柄
		file *os.File

		mu sync.Mutex
	}

	// Options 日志自动备份配置
	Options struct {
		// 日志文件名称
		FileName string `json:"file_name" mapstructure:"file_name"`

		// 日志备份路径
		BackupPath string `json:"backup_path" mapstructure:"backup_path"`

		// 日志文件最大字节数，单位为MB
		MaxSize int `json:"max_size" mapstructure:"max_size"`

		// 备份文件，日期后缀格式
		SuffixFormat string `json:"suffix_format" mapstructure:"suffix_format"`
	}
)

func (opts *Options) loadDefault() {
	if opts.MaxSize <= 0 {
		opts.MaxSize = 10
	}
	if opts.FileName == "" {
		opts.FileName = filepath.Join(os.TempDir(), filepath.Base(os.Args[0])+".log")
	}
	if opts.BackupPath == "" {
		opts.BackupPath = filepath.Dir(opts.FileName)
	}
	if opts.SuffixFormat == "" {
		opts.SuffixFormat = "20060102T150405"
	}
}

// NewWriter 根据配置信息创建 RotateWriter
func NewWriter(opts *Options) *RotateWriter {
	// 检查配置参数，加载默认值
	if opts == nil {
		opts = &Options{}
	}
	opts.loadDefault()

	return &RotateWriter{
		filename:     opts.FileName,
		backupPath:   opts.BackupPath,
		maxSize:      int64(opts.MaxSize * 1024 * 1024),
		suffixFormat: opts.SuffixFormat,
	}
}

// Write 实现 Writer 接口
func (r *RotateWriter) Write(data []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	dataLen := int64(len(data))
	if dataLen > r.maxSize {
		return 0, errors.New("write length exceeds maximum file size")
	}
	if r.file == nil {
		if err = r.open(); err != nil {
			return 0, err
		}
	}

	if r.curSize+dataLen > r.maxSize {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = r.file.Write(data)

	r.curSize += int64(n)
	return n, err
}

func (r *RotateWriter) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.close()
}

func (r *RotateWriter) close() error {
	if r.file == nil {
		return nil
	}
	err := r.file.Close()
	r.file = nil
	return err
}

func (r *RotateWriter) open() error {
	info, err := os.Stat(r.filename)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error getting log file info: %s", err)
	}

	if err := os.MkdirAll(filepath.Dir(r.filename), 0755); err != nil {
		return fmt.Errorf("can't make directories for logfile: %s", err)
	}

	file, err := os.OpenFile(r.filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	r.curSize = 0
	r.file = file
	if info != nil {
		r.curSize = info.Size()
	}
	return nil
}

func (r *RotateWriter) rotate() error {
	if err := r.close(); err != nil {
		return err
	}
	err := os.MkdirAll(r.backupPath, 0755)
	if err != nil {
		return fmt.Errorf("can't make directories for backup: %s", err)
	}

	if err := os.Rename(r.filename, r.backupName()); err != nil {
		return fmt.Errorf("can't rename log file: %s", err)
	}

	f, err := os.OpenFile(r.filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("can't open new logfile: %s", err)
	}
	r.file = f
	r.curSize = 0
	return nil
}

func (r *RotateWriter) backupName() string {
	buf := bufferpool.Get()
	defer buf.Free()
	ext := filepath.Ext(r.filename)
	buf.WriteString(strings.TrimSuffix(filepath.Base(r.filename), ext))

	timeSuffix := time.Now().Format(r.suffixFormat)
	if timeSuffix != r.timeSuffix {
		r.timeSuffix = timeSuffix
		r.index = 0
	}

	buf.WriteByte('-')
	buf.WriteString(r.timeSuffix)
	if r.index > 0 {
		buf.WriteByte('-')
		buf.WriteInt(r.index)
	}
	r.index++

	buf.WriteString(ext)
	return filepath.Join(r.backupPath, buf.String())
}
