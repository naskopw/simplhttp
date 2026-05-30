package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/naskopw/simplhttp/internal/server"
	"github.com/naskopw/simplhttp/internal/util"
	"github.com/spf13/cobra"
)

var (
	port       int
	dir        string
	auth       string
	readOnly   bool
	maxSizeStr string
	certFile   string
	keyFile    string
	https      bool
)

var rootCmd = &cobra.Command{
	Use:   "simplhttp",
	Short: "A modern, bidirectional HTTP file server",
	Long: `simplhttp is a modern replacement for python -m http.server.
It features a React UI, support for file uploads, directory navigation, 
and basic authentication.`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")
	rootCmd.Flags().StringVarP(&dir, "dir", "d", ".", "Root directory to serve")
	rootCmd.Flags().StringVarP(&auth, "auth", "a", "", "Basic auth credentials in user:pass format")
	rootCmd.Flags().BoolVarP(&readOnly, "readonly", "r", false, "Disable uploads and file modifications")
	rootCmd.Flags().StringVarP(&maxSizeStr, "max-size", "m", "100MB", "Max upload size per request (e.g. 100MB, 1GB)")
	rootCmd.Flags().StringVar(&certFile, "cert", "", "Path to TLS certificate file")
	rootCmd.Flags().StringVar(&keyFile, "key", "", "Path to TLS private key file")
	rootCmd.Flags().BoolVar(&https, "https", false, "Enable HTTPS (self-signed if cert/key are missing)")
}

func runServer() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	maxSizeBytes, err := util.ParseHumanSize(maxSizeStr)
	if err != nil {
		slog.Error("invalid max size", "value", maxSizeStr, "error", err)
		os.Exit(1)
	}

	absDir, err := os.Stat(dir)
	if err != nil {
		slog.Error("failed to access directory", "path", dir, "error", err)
		os.Exit(1)
	}
	if !absDir.IsDir() {
		slog.Error("path is not a directory", "path", dir)
		os.Exit(1)
	}

	protocol := "http"
	useTLS := https || (certFile != "" && keyFile != "")
	if useTLS {
		protocol = "https"
	}

	ips, err := util.GetLocalIPs()
	if err != nil {
		slog.Warn("could not determine local IP addresses", "error", err)
	} else {
		fmt.Println("Server starting...")
		for _, ip := range ips {
			fmt.Printf("  %s://%s:%d\n", protocol, ip, port)
		}
		fmt.Printf("  %s://localhost:%d\n", protocol, port)
	}

	slog.Info("server starting", "port", port, "dir", dir, "readonly", readOnly, "maxSize", maxSizeStr, "protocol", protocol)

	srv := server.NewServer(server.Config{
		Port:         port,
		Dir:          dir,
		Auth:         auth,
		ReadOnly:     readOnly,
		MaxSizeBytes: maxSizeBytes,
		CertFile:     certFile,
		KeyFile:      keyFile,
		AutoTLS:      https,
	})

	if err := srv.Start(); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
