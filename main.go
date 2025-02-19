package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"FileServer/modules"

	"github.com/kendfss/but"
	"github.com/labstack/echo/v4"
)

func main() {
	// Initalize
	e := echo.New()
	but.Must(os.MkdirAll("downloads", os.ModePerm))

	dom := fmt.Sprintf("%s:%s", IP(), PORT())
	link := fmt.Sprintf("http://%s", dom)

	// Routes
	e.GET("/", modules.IndexPage)
	e.GET("/upload", modules.UploadPage)
	e.POST("/api/upload", modules.HandleUpload)
	e.GET("/dl/id/:id", modules.DownloadFile)
	e.GET("/dl/name/:name", modules.DownloadFile)
	e.GET("/download", modules.GetFiles)
	// e.GET("/downloads", modules.GetFiles)
	// e.GET("/downloads", modules.GetFiles2)
	// e.GET("/downloads", modules.GetFiles3)
	e.GET("assets/qrcode.png", modules.QRCode(link))

	fmt.Printf("serving on %s\n", link)

	// Start server
	// e.Logger.Fatal(e.Start(":" + port))
	e.Logger.Fatal(e.Start(dom))
}

func IP() string {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("hostname", "-I")
	}
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(bytes))
}

func PORT() string {
	out := os.Getenv("PORT")
	if out == "" {
		return strconv.Itoa(rand.Intn(math.MaxInt16 + 1))
	}
	return out
}
