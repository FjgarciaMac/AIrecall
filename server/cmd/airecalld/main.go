// Command airecalld runs the AIrecall memory server.
//
// A small JSON API over the same SQLite schema the SDK uses in dev
// mode, plus hybrid retrieval and the summarization pass.
package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/FjgarciaMac/AIrecall/server/internal/api"
	"github.com/FjgarciaMac/AIrecall/server/internal/storage"
)

func main() {
	var (
		addr = flag.String("addr", "127.0.0.1:8734", "listen address")
		db   = flag.String("db", "airecall.db", "sqlite store path")
	)
	flag.Parse()

