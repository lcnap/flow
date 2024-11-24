package main

import (
	"log/slog"
	"os"
)

func NewLogger(path string) *slog.Logger {

	out, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{
		//AddSource: true,
		Level: slog.LevelDebug,
	}))

}
