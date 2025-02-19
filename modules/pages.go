package modules

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/kendfss/but"
	"github.com/labstack/echo/v4"
	qrcode "github.com/skip2/go-qrcode"
)

func QRCode(url string) func(c echo.Context) error {
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return func(c echo.Context) error {
			f, err := os.Open("assets/404.png")
			if err == nil {
				return c.Stream(http.StatusInternalServerError, "image/png", f)
			}
			return c.String(http.StatusInternalServerError, "Sorry, but we've nothing to show.")
		}
	}
	return func(c echo.Context) error {
		buf := bytes.NewBuffer(png)
		return c.Stream(http.StatusOK, "image/png", buf)
	}
}

func IndexPage(c echo.Context) error {
	return c.File("html/index.html")
}

func UploadPage(c echo.Context) error {
	return c.File("html/upload.html")
}

func GetFiles(c echo.Context) error {
	flist := make([]string, 0)
	files, err := os.ReadDir("downloads")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		ahref := fmt.Sprintf("<li><a href=\"/dl/name/%s\">%s</a></li>", f.Name(), f.Name())
		flist = append(flist, ahref)
	}
	// return c.HTML(http.StatusOK, fmt.Sprintf("<h2>Files</h2><p>%s</p>", strings.Join(flist, "<br>")))
	page := fmt.Sprintf(
		`<html>
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<h2>Download</h2>
		<title>Download</title>
	</head>
	<body>
		<ul>
			%s
		</ul>
	</body>
</html>`, strings.Join(flist, "<br>"))
	return c.HTML(http.StatusOK, page)
}

func GetFiles2(c echo.Context) error {
	files, linkList := make(chan file), []string{}
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		for f := range files {
			ahref := fmt.Sprintf("<li><a href=\"/dl/name/%s\">%s</a></li>", f.Name(), f.Name())
			linkList = append(linkList, ahref)
		}
		wg.Done()
	}()
	but.Must(fs.WalkDir(fsys{}, dd, walker(files)))
	wg.Wait()
	page := fmt.Sprintf(
		`<html>
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<h2>Download</h2>
		<title>Download</title>
	</head>
	<body>
		<ol>
			%s
		</ol>
	</body>
</html>`, strings.Join(linkList, "<br>"))
	return c.HTML(http.StatusOK, page)
}

func GetFiles3(c echo.Context) error {
	linkList := []string{}
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		for f := range Walker(pwd, nil, wg) {
			ahref := fmt.Sprintf("<a href=\"/dl/name/%s\">%s</a>", f.Name(), f.Name())
			linkList = append(linkList, ahref)
		}
		// wg.Done()
	}()
	wg.Wait()
	return c.HTML(http.StatusOK, fmt.Sprintf("<h2>Files</h2><p>%s</p>", strings.Join(linkList, "<br>")))
}

func Listdir(root string) []string {
	entries, err := os.ReadDir(root)
	but.Must(err)
	out := make([]string, len(entries))
	for i, file := range entries {
		out[i] = file.Name()
	}
	return out
}

func Walker(root string, ch chan file, wg *sync.WaitGroup) <-chan file {
	but.MustBool(wg != nil, "WaitGroup is nil")
	go func() {
		if ch == nil {
			ch = make(chan file)
			defer close(ch)
			defer wg.Done()
		}
		for _, base := range Listdir(root) {
			path := filepath.Join(root, base)
			stat, err := os.Lstat(path)
			but.Must(err)
			switch mode := stat.Mode(); {
			case mode.IsRegular():
				// if slices.Contains(cli.Args(), base) {
				ch <- file{path: path}
				// }
			case mode.IsDir():
				Walker(path, ch, wg)
			}
		}
	}()
	return ch
}

func walker(files chan file) func(path string, d fs.DirEntry, err error) error {
	return func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			files <- file{path: path}
		}
		return nil
	}
}

type fsys struct{}

func (fsys) Open(name string) (fs.File, error) {
	return os.Open(name)
}

type file struct {
	path string
}

func (f file) Name() string {
	p, err := filepath.Rel(dd, f.path)
	but.Must(err)
	return p
}

var (
	pwd = getwd()
	dd  = getdd()
)

func getwd() string {
	p, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return p
}

func getdd() string {
	p := getwd()
	if !strings.HasSuffix(p, "downloads") {
		return filepath.Join(p, "downloads")
	}
	return p
}
