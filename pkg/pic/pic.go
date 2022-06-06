package pic

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"time"

	"github.com/blackjack/webcam"
)

type CamConfig struct {
	PicDev string `arg:"env:PIC_DEV,--pic-dev" default:"/dev/video2" help:"camera video device file path" placeholder:"DEV"`
	// v4l2-ctl --list-formats-ext --device /dev/video2
	PicFormat     string        `arg:"env:PIC_FORMAT,--pic-format" default:"Motion-JPEG" help:"camera preferred image format" placeholder:"STR"`
	PicTimeout    time.Duration `arg:"env:PIC_TIMEOUT,--pic-timeout" default:"2s" help:"how long to give camera time to start" placeholder:"DUR"`
	PicSkipFrames int           `arg:"env:PIC_SKIP_FRAMES,--pic-skip-frames" default:"15" help:"how many frames to skip until picture snap" placeholder:"N"`
}

func chooseFormat(cam *webcam.Webcam, preferred string) (webcam.PixelFormat, string, *webcam.FrameSize, error) {
	fmap := cam.GetSupportedFormats()
	var format webcam.PixelFormat = 0
	var formatStr string
	for f, s := range fmap {
		format = f
		formatStr = s
		if s == preferred {
			break
		}
	}

	if format == 0 {
		return 0, "", nil, errors.New("no format found")
	}

	frameSizes := cam.GetSupportedFrameSizes(format)
	if len(frameSizes) == 0 {
		return 0, "", nil, errors.New("no frame size found")
	}

	return format, formatStr, &frameSizes[0], nil
}

type YCbCr struct {
	rect image.Rectangle
	buf  []byte
}

// compile time interface check
var _ image.Image = (*YCbCr)(nil)

func (i *YCbCr) ColorModel() color.Model {
	return color.YCbCrModel
}
func (i *YCbCr) Bounds() image.Rectangle {
	return i.rect
}
func (i *YCbCr) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(i.rect)) {
		return color.YCbCr{}
	}
	pixIx := (y*i.rect.Max.X + x)
	// two pixels = four bytes
	cellIx := (pixIx / 2) * 4

	Y0, Cb, Y1, Cr := i.buf[cellIx+0], i.buf[cellIx+1], i.buf[cellIx+2], i.buf[cellIx+3]

	Y := Y0
	if pixIx%2 == 1 {
		Y = Y1
	}

	return color.YCbCr{
		Y, Cb, Cr,
	}
}

func convertFrame(frame []byte, formatStr string, w, h int) (image.Image, error) {
	rect := image.Rect(0, 0, int(w), int(h))

	switch formatStr {
	case "YUYV 4:2:2":
		buf := make([]byte, len(frame))
		copy(buf, frame)
		return &YCbCr{
			rect: rect,
			buf:  buf,
		}, nil
	case "Motion-JPEG":
		b := bytes.NewBuffer(frame)
		return jpeg.Decode(b)
	default:
		return nil, fmt.Errorf(`unhandled format string "%s"`, formatStr)
	}
}

func Snap(config CamConfig) (image.Image, error) {
	cam, err := webcam.Open(config.PicDev)
	if err != nil {
		return nil, err
	}
	defer cam.Close()

	format, formatStr, size, err := chooseFormat(cam, config.PicFormat)
	if err != nil {
		return nil, err
	}

	f, w, h, err := cam.SetImageFormat(format, uint32(size.MaxWidth), uint32(size.MaxHeight))
	if err != nil {
		return nil, err
	}
	if f != format {
		return nil, errors.New("unable to choose format")
	}

	err = cam.StartStreaming()
	if err != nil {
		return nil, err
	}

	timeout := uint32(config.PicTimeout.Seconds())
	var frame []byte
	for i := 0; i < config.PicSkipFrames; i++ {
		err = cam.WaitForFrame(timeout)

		switch err.(type) {
		case nil:
		case *webcam.Timeout:
			return nil, errors.New("camera timed out")
		default:
			return nil, err
		}

		frame, err = cam.ReadFrame()
		if err != nil {
			return nil, err
		}
		if len(frame) == 0 {
			return nil, errors.New("received empty frame")
		}
	}

	return convertFrame(frame, formatStr, int(w), int(h))
}
