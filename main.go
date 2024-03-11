package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		html := `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Upload Photo</title>
			</head>
			<body>
				<form action="/upload" method="POST" enctype="multipart/form-data">
					<input type="file" name="photo" />
					<input type="submit" value="Upload" />
				</form>
			</body>
			</html>
		`
		return c.HTML(http.StatusOK, html)
	})

	e.POST("/upload", func(c echo.Context) error {
		// Source
		file, err := c.FormFile("photo")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		fileExtension := strings.ToLower(filepath.Ext(file.Filename))
		dst, err := os.Create("clip" + fileExtension)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		err = copyImageToClipboard()
		if err != nil {
			panic(err)
		}
		return c.HTML(http.StatusOK, "<p>File uploaded successfully.</p>")
	})

	e.Logger.Fatal(e.Start(":80"))
}

func copyImageToClipboard() error {
	// Read the JpG as a png so it can be pasted into obsidian
	cmd := exec.Command("xclip", "-selection", "clipboard", "-target", "image/png", "-i", "clip.jpg")
	return cmd.Run()
}
