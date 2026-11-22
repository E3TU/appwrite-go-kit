package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/e3tu/appwrite-base"
	"github.com/e3tu/appwrite-base/auth"
	"github.com/gin-gonic/gin"
)

// serve runs a Gin server with example routes that exercise the enabled
// capabilities. Routes for disabled capabilities are not registered.
func serve() error {
	cfg, err := loadEnv()
	if err != nil {
		return err
	}

	b, err := base.New(cfg)
	if err != nil {
		return err
	}

	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	if svc, ok := b.Auth(); ok {
		guarded := r.Group("/api")
		guarded.Use(svc.Middleware(cfg.Auth.SessionHeader))
		guarded.GET("/me", meHandler)
	}

	if db, ok := b.DB(); ok {
		r.GET("/api/tables/:table", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"database": db.DatabaseID(),
				"table":    c.Param("table"),
			})
		})
	}

	if st, ok := b.Storage(); ok {
		r.GET("/api/files/:id/url", func(c *gin.Context) {
			u, err := st.PublicFileURL("", c.Param("id"))
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"url": u})
		})
	}

	if fn, ok := b.Functions(); ok {
		r.GET("/api/functions", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"enabled": fn.Functions() != nil})
		})
	}

	addr := listenAddr()

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errc := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errc <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errc:
		return err
	case sig := <-stop:
		fmt.Fprintf(os.Stderr, "received %s, shutting down\n", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}

func meHandler(c *gin.Context) {
	u, ok := auth.UserFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":    u.ID,
		"email": u.Email,
		"name":  u.Name,
	})
}
