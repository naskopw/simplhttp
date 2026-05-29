package server

import (
	"archive/zip"
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/naskopw/simphttp/internal/util"
	"github.com/naskopw/simphttp/ui"
)

type Config struct {
	Port         int
	Dir          string
	Auth         string
	ReadOnly     bool
	MaxSizeBytes int64
	CertFile     string
	KeyFile      string
	AutoTLS      bool
}

type Server struct {
	echo   *echo.Echo
	config Config
}

type FileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	IsDir   bool      `json:"isDir"`
	ModTime time.Time `json:"modTime"`
}

type ServerConfig struct {
	MaxSizeBytes int64 `json:"maxSizeBytes"`
	ReadOnly     bool  `json:"readOnly"`
}

func NewServer(config Config) *Server {
	e := echo.New()
	e.Logger = slog.Default()

	s := &Server{
		echo:   e,
		config: config,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

func (s *Server) setupMiddleware() {
	// Request logging using slog
	s.echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogMethod: true,
		HandleError: true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			slog.Info("request",
				"method", v.Method,
				"uri", v.URI,
				"status", v.Status,
			)
			return nil
		},
	}))

	s.echo.Use(middleware.Recover())

	// Basic Auth if configured
	if s.config.Auth != "" {
		parts := strings.SplitN(s.config.Auth, ":", 2)
		if len(parts) == 2 {
			s.echo.Use(middleware.BasicAuth(func(c *echo.Context, username, password string) (bool, error) {
				// Use constant-time comparison to prevent timing attacks
				userMatch := subtle.ConstantTimeCompare([]byte(username), []byte(parts[0])) == 1
				passMatch := subtle.ConstantTimeCompare([]byte(password), []byte(parts[1])) == 1
				return userMatch && passMatch, nil
			}))
		} else {
			slog.Warn("invalid auth format, expected user:pass. Basic auth disabled.")
		}
	}
}

func (s *Server) setupRoutes() {
	// API routes
	api := s.echo.Group("/api")
	api.GET("/health", func(c *echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	api.GET("/config", s.handleConfig)
	api.GET("/fs/*", s.handleFSGet)
	api.GET("/zip/*", s.handleZipDownload)
	api.POST("/upload/*", s.handleUpload)

	// Static assets from embedded FS
	subFS := echo.MustSubFS(ui.Assets, "dist")

	// Serve UI assets
	s.echo.StaticFS("/", subFS)

	// SPA Fallback: for any route not found, serve index.html
	s.echo.HTTPErrorHandler = func(c *echo.Context, err error) {
		if err == echo.ErrNotFound {
			// If not an API route, serve index.html
			if !strings.HasPrefix(c.Request().URL.Path, "/api") {
				indexFile, err := subFS.Open("index.html")
				if err != nil {
					slog.Error("failed to open index.html", "error", err)
					return
				}
				defer indexFile.Close()
				http.ServeContent(c.Response(), c.Request(), "index.html", time.Now(), indexFile.(io.ReadSeeker))
				return
			}
		}
		echo.DefaultHTTPErrorHandler(true)(c, err)
	}
}

func (s *Server) Start() error {
	address := fmt.Sprintf(":%d", s.config.Port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc := echo.StartConfig{
		Address:    address,
		HideBanner: true,
	}

	useTLS := s.config.AutoTLS || (s.config.CertFile != "" && s.config.KeyFile != "")
	if useTLS {
		var cert, key any
		if s.config.AutoTLS && (s.config.CertFile == "" || s.config.KeyFile == "") {
			certPEM, keyPEM, err := util.GenerateSelfSignedCert()
			if err != nil {
				return fmt.Errorf("failed to generate self-signed cert: %w", err)
			}
			cert = certPEM
			key = keyPEM
		} else {
			cert = s.config.CertFile
			key = s.config.KeyFile
		}

		err := sc.StartTLS(ctx, s.echo, cert, key)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}

	err := sc.Start(ctx, s.echo)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}


func (s *Server) handleConfig(c *echo.Context) error {
	return c.JSON(http.StatusOK, ServerConfig{
		MaxSizeBytes: s.config.MaxSizeBytes,
		ReadOnly:     s.config.ReadOnly,
	})
}

func (s *Server) resolvePath(p string) (string, error) {
	// Clean the path to prevent traversal
	cleanPath := filepath.Clean(p)
	if strings.HasPrefix(cleanPath, "..") {
		return "", fmt.Errorf("invalid path")
	}

	absRoot, err := filepath.Abs(s.config.Dir)
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(absRoot, cleanPath)

	// Ensure the fullPath is still within the absRoot
	if !strings.HasPrefix(fullPath, absRoot) {
		return "", fmt.Errorf("invalid path")
	}

	return fullPath, nil
}

func (s *Server) handleFSGet(c *echo.Context) error {
	relPath := c.Param("*")
	fullPath, err := s.resolvePath(relPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid path")
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return echo.NewHTTPError(http.StatusNotFound, "File not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if info.IsDir() {
		entries, err := os.ReadDir(fullPath)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		files := make([]FileInfo, 0, len(entries))
		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			files = append(files, FileInfo{
				Name:    entry.Name(),
				Size:    info.Size(),
				IsDir:   entry.IsDir(),
				ModTime: info.ModTime(),
			})
		}
		return c.JSON(http.StatusOK, files)
	}

	return c.File(fullPath)
}

func (s *Server) handleUpload(c *echo.Context) error {
	if s.config.ReadOnly {
		return echo.NewHTTPError(http.StatusForbidden, "Server is in read-only mode")
	}

	// Enforce body limit
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, s.config.MaxSizeBytes)

	relPath := c.Param("*")
	fullPath, err := s.resolvePath(relPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid path")
	}

	// Ensure target directory exists
	targetDir := fullPath
	if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
		// If fullPath is a file, we probably want to upload to its directory?
		// Or maybe the relPath should always be a directory.
		// Let's assume relPath is the directory to upload into.
	} else if err != nil && !os.IsNotExist(err) {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dstPath := filepath.Join(targetDir, file.Filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusCreated)
}

func (s *Server) handleZipDownload(c *echo.Context) error {
	relPath := c.Param("*")
	fullPath, err := s.resolvePath(relPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid path")
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return echo.NewHTTPError(http.StatusNotFound, "Directory not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !info.IsDir() {
		return echo.NewHTTPError(http.StatusBadRequest, "Path is not a directory")
	}

	zipName := filepath.Base(fullPath) + ".zip"
	if zipName == "." + ".zip" || zipName == "/" + ".zip" || zipName == ".zip" {
		zipName = "root.zip"
	}

	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%q", zipName))
	c.Response().Header().Set(echo.HeaderContentType, "application/zip")
	c.Response().WriteHeader(http.StatusOK)

	zw := zip.NewWriter(c.Response())
	defer zw.Close()

	err = filepath.WalkDir(fullPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(fullPath, path)
		if err != nil {
			return err
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		w, err := zw.Create(rel)
		if err != nil {
			return err
		}

		_, err = io.Copy(w, f)
		return err
	})

	return err
}
