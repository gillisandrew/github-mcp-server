package cmd

import (
	"context"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/github/github-mcp-server/pkg/github"
	iolog "github.com/github/github-mcp-server/pkg/log"
	"github.com/github/github-mcp-server/pkg/translations"
	gogithub "github.com/google/go-github/v69/github"
	"github.com/mark3labs/mcp-go/server"
)

func RunStdIO(cfg runConfig) error {
	// Create app context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create GitHub API client

	ghClient, err := getGitHubClient(cfg.Host, cfg.AuthToken)

	if err != nil {
		cfg.Logger.Fatalf("failed to create HTTP client: %v", err)
	}

	ghClient.UserAgent = fmt.Sprintf("github-mcp-server/%s", cfg.Version)

	// Create MCP server
	t, dumpTranslations := translations.TranslationHelper()

	ghServer := github.NewServer(func(_ context.Context) (*gogithub.Client, error) {
		return ghClient, nil // closing over client
	}, cfg.Version, cfg.ReadOnly, t)
	stdioServer := server.NewStdioServer(ghServer)

	stdLogger := stdlog.New(cfg.Logger.Writer(), "stdioserver", 0)
	stdioServer.SetErrorLogger(stdLogger)

	if cfg.ExportTranslations {
		// Once server is initialized, all translations are loaded
		dumpTranslations()
	}

	// Start listening for messages
	errC := make(chan error, 1)
	go func() {
		in, out := io.Reader(os.Stdin), io.Writer(os.Stdout)

		if cfg.LogCommands {
			loggedIO := iolog.NewIOLogger(in, out, cfg.Logger)
			in, out = loggedIO, loggedIO
		}

		errC <- stdioServer.Listen(ctx, in, out)
	}()

	// Output github-mcp-server string
	_, _ = fmt.Fprintf(os.Stderr, "GitHub MCP Server running on stdio\n")

	// Wait for shutdown signal
	select {
	case <-ctx.Done():
		cfg.Logger.Infof("shutting down server...")
	case err := <-errC:
		if err != nil {
			return fmt.Errorf("error running server: %w", err)
		}
	}

	return nil
}
