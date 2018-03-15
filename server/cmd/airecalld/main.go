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
