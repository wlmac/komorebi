// Komorebi, middleware that edits images to save bandwidth.
// Copyright (C) 2023  nyiyui
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nyiyui/komorebi/notify"
	"github.com/nyiyui/komorebi/server"
	"gopkg.in/gographics/imagick.v3/imagick"
)

var lisNet string
var lisAddr string
var sourcePath string
var configRaw string

func main() {
	fmt.Print(`Komorebi  Copyright (C) 2023  nyiyui
This program comes with ABSOLUTELY NO WARRANTY; for details see the license.
This is free software, and you are welcome to redistribute it
under certain conditions; see the license for details.

You should have received a copy of the license used,
the GNU General Public License along with this program.
If not, see <https://www.gnu.org/licenses/>.
`)
	flag.StringVar(&lisNet, "net", "", "listen network (e.g. unix)")
	flag.StringVar(&lisAddr, "addr", "", "listen address (e.g. /tmp/komorebi.sock)")
	flag.StringVar(&sourcePath, "src", "", "source media path")
	flag.StringVar(&configRaw, "cfg", "", "config in JSON (not path)")
	flag.Parse()

	imagick.Initialize()
	defer imagick.Terminate()

	config := server.Config{}
	err := json.Unmarshal([]byte(configRaw), &config)
	if err != nil {
		log.Fatalf("parse config: %s", err)
	}
	config.SourcePath = sourcePath
	config.CachePath = filepath.Join(os.Getenv("CACHE_DIRECTORY"), "media")

	lis, err := net.Listen(lisNet, lisAddr)
	if err != nil {
		log.Fatalf("listen: %s", err)
	}
	handler, err := server.New(config)
	if err != nil {
		log.Fatalf("setup: %s", err)
	}
	notify.Notify("READY=1")
	log.Fatalf("serve: %s", http.Serve(lis, handler))
}
