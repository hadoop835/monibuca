package record

import (
	. "github.com/langhuihui/monibuca/monica"
	"github.com/langhuihui/monibuca/monica/avformat"
	"github.com/langhuihui/monibuca/monica/pool"
	"github.com/langhuihui/monibuca/monica/util"
	"os"
	"syscall"
)

func SaveFlv(streamPath string, append bool) error {
	flag := os.O_CREATE
	if append {
		flag = flag | os.O_RDWR | os.O_APPEND
	} else {
		flag = flag | os.O_TRUNC | os.O_WRONLY
	}
	filePath := config.Path + streamPath + ".flv"
	file, err := os.OpenFile(filePath, flag, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	p := OutputStream{SendHandler: func(packet *pool.SendPacket) error {
		return avformat.WriteFLVTag(file, packet)
	}}
	p.ID = filePath
	p.Type = "FlvRecord"
	if append {
		_, err = file.Seek(4, syscall.FILE_END)
		if err == nil {
			var tagSize uint32
			if tagSize, err = util.ReadByteToUint32(file, true); err == nil {
				_, err = file.Seek(int64(tagSize+4), syscall.FILE_END)
				if err == nil {
					var tag *pool.AVPacket
					tag, err = avformat.ReadFLVTag(file)
					if err == nil {
						p.OffsetTime = tag.Timestamp
					}
				}
			}
		}
	} else {
		_, err = file.Write(avformat.FLVHeader)
	}
	if err == nil {
		recordings.Store(filePath, &p)
		go p.Play(streamPath)
	}
	return err
}