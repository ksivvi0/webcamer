package webcamer

import (
	"fmt"
	"gocv.io/x/gocv"
	"time"
)

type WebCamerer interface {
	DoSnapshot() (string, error)
	DoVideo(sec int) (string, error)
}

type Webcamer struct {
	deviceID int
	basePath string
}

func NewWebcamer(devID int) WebCamerer {
	return &Webcamer{deviceID: devID, basePath: "/tmp"}
}

func (w *Webcamer) DoSnapshot() (string, error) {
	webcam, err := gocv.OpenVideoCapture(w.deviceID)
	if err != nil {
		return "", err
	}
	defer webcam.Close()

	gocv.NewMat()
	img := gocv.NewMat()
	defer img.Close()

	if ok := webcam.Read(&img); !ok {
		return "", err
	}
	if img.Empty() {
		return "", err
	}
	outFile := fmt.Sprintf("%s/%d.jpg", w.basePath, time.Now().UnixNano())
	gocv.IMWrite(outFile, img)
	return outFile, err
}

func (w *Webcamer) DoVideo(sec int) (string, error) {
	webcam, err := gocv.OpenVideoCapture(w.deviceID)
	if err != nil {
		return "", err
	}
	defer webcam.Close()

	gocv.NewMat()
	img := gocv.NewMat()
	defer img.Close()

	if ok := webcam.Read(&img); !ok {
		return "", err
	}

	outFile := fmt.Sprintf("%s/%d.mp4", w.basePath, time.Now().UnixNano())
	writer, err := gocv.VideoWriterFile(outFile, "mp4v", 25, img.Cols(), img.Rows(), true)
	if err != nil {
		return "", err
	}
	defer writer.Close()

	for i := 0; i < 25*sec; i++ {
		if ok := webcam.Read(&img); !ok {
			return "", err
		}
		if img.Empty() {
			continue
		}

		_ = writer.Write(img)
	}

	return outFile, nil
}
