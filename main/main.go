package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"justchess/pkg/auth"
	"justchess/pkg/db"
	"justchess/pkg/middleware"
	"justchess/pkg/ws"

	"github.com/joho/godotenv"
)

func main() {
	setupLogger()
	// Load environment variables.
	err := godotenv.Load()
	if err != nil {
		slog.Error(".env file cannot be loaded", "err", err)
		return
	}
	// Setup db.
	err = db.OpenDB()
	defer db.CloseDB()
	if err != nil {
		slog.Error("Database cannot be opened.", "err", err)
		return
	}
	err = db.CreateTables("./pkg/db/schema.sql")
	if err != nil {
		slog.Error("Tables cannot be created.", "err", err)
		return
	}
	// Setup routes.
	router := setupRouter()
	http.ListenAndServe(":3502", router)
}

func setupRouter() *http.ServeMux {
	// Setup the chain of middlewares.
	authStack := middleware.CreateStack(
		middleware.LogRequest,
		middleware.AllowCors,
	)
	router := http.NewServeMux()
	router.Handle("/auth/", http.StripPrefix(
		"/auth",
		authStack(auth.AuthRouter()),
	))
	// Instantiate manager to handle ws connections.
	m := ws.NewManager()
	router.HandleFunc("/ws", m.HandleConnection)
	return router
}

func setupLogger() {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(time.Now().Format("01/02/2006 15:04:05"))
			}
			return a
		},
	})
	slog.SetDefault(slog.New(h))
}