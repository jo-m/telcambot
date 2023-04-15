package libcamera

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os/exec"
	"sync"

	"github.com/rs/zerolog/log"
)

type Config struct {
	CameraIx int `arg:"env:LIBCAMERA_CAMERA_IX,--libcamera-camera-ix" default:"0" help:"libcamera camera number, see 'libcamera-jpeg --list-cameras'" placeholder:"N"`
	Width    int `arg:"env:LIBCAMERA_W,--libcamera-w" default:"1920" help:"width" placeholder:"W"`
	Height   int `arg:"env:LIBCAMERA_H,--libcamera-h" default:"1080" help:"height" placeholder:"H"`
}

func processErr(e io.Reader) {
	scanner := bufio.NewScanner(e)
	for scanner.Scan() {
		line := scanner.Text()
		log.Info().Str("src", "stderr").Msg(line)
	}
}

func Snap(c Config) (image.Image, error) {
	args := []string{
		"--output", "-",
		"--camera", fmt.Sprint(c.CameraIx),
		"--timeout", "100",
		"--quality", "100",
		"--nopreview",
		"--immediate",
		"--width", fmt.Sprint(c.Width),
		"--height", fmt.Sprint(c.Height),
	}

	log.Info().Strs("args", args).Msg("libcamera-jpeg args")
	cmd := exec.Command("libcamera-jpeg", args...)

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	errPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	go processErr(errPipe)

	out := bytes.Buffer{}
	var outErr error
	wait := sync.WaitGroup{}
	wait.Add(1)
	go func() {
		defer wait.Done()
		_, outErr = io.Copy(&out, outPipe)
	}()

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("running libcamera-jpeg failed: %w", err)
	}

	wait.Wait()
	if outErr != nil {
		return nil, fmt.Errorf("reading libcamera-jpeg stdout failed: %w", outErr)
	}

	img, err := jpeg.Decode(&out)
	if err != nil {
		return nil, fmt.Errorf("decoding libcamera-jpeg output failed: %w", outErr)
	}

	return img, nil
}
