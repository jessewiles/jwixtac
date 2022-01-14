package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

const feed1 = "https://hacker-news.firebaseio.com/v0/topstories.json"
const feed2 = "https://api.imgflip.com/get_memes"

// FeedWSPayload is a thing
type FeedWSPayload struct {
	Channel string `json:"channel"`
	Message []int  `json:"message"`
}

//FeedWSIFPayload thing
type FeedWSIFPayload struct {
	Channel string `json:"channel"`
	Message []Meme `json:"message"`
}

//ImgFlipPayl payload
type ImgFlipPayl struct {
	Data ImgFlipData `json:"data"`
}

//ImgFlipData data
type ImgFlipData struct {
	Memes []Meme `json:"memes"`
}

//Meme imgflip  meme
type Meme struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	BoxCount int    `json:"box_count"`
}

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

func hnFeed() {
	log.Info("Opening websocket connection...")
	origin := "http://localhost:8088/"
	url := "ws://localhost:8088/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	log.Debug("Opened websocket connection...", ws)
	for {
		var reup = false

		resp, _ := http.Get(feed1)
		var fj []int

		if resp.Body != nil {
			defer resp.Body.Close()
		}

		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		jsonErr := json.Unmarshal(body, &fj)
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		payload := FeedWSPayload{
			Channel: "topstories",
			Message: fj,
		}

		bytes, err := json.Marshal(payload)
		if err != nil {
			log.Warning("Issues marshaling payload...")
			log.Error(err)
		}

		if _, err := ws.Write(bytes); err != nil {
			log.Error(err)
			reup = true
		}

		// discard read bytes
		ws.Read(bytes)
		bytes = []byte{}

		if reup {
			// Close old conneection and open a new one
			ws.Close()

			ws, err = websocket.Dial(url, "", origin)
			if err != nil {
				log.Fatal(err)
			}
			defer ws.Close()
		}
		time.Sleep(15 * time.Second)
	}
}

func imgFlipFeed() {
	log.Info("Opening websocket connection...")
	origin := "http://localhost:8088/"
	url := "ws://localhost:8088/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	log.Debug("Opened websocket connection...", ws)
	for {
		var reup = false

		resp, _ := http.Get(feed2)
		var ifpayl ImgFlipPayl

		if resp.Body != nil {
			defer resp.Body.Close()
		}

		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		jsonErr := json.Unmarshal(body, &ifpayl)
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		payload := FeedWSIFPayload{
			Channel: "imgflip",
			Message: ifpayl.Data.Memes,
		}

		bytes, err := json.Marshal(payload)
		if err != nil {
			log.Warning("Issues marshaling payload...")
			log.Error(err)
		}

		log.Info("HACK: ", string(bytes[0:48]))
		if _, err := ws.Write(bytes); err != nil {
			log.Error(err)
			reup = true
		}

		// discard read bytes
		ws.Read(bytes)
		bytes = []byte{}

		if reup {
			// Close old conneection and open a new one
			ws.Close()

			ws, err = websocket.Dial(url, "", origin)
			if err != nil {
				log.Fatal(err)
			}
			defer ws.Close()
		}
		time.Sleep(15 * time.Second)
	}
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)
	log.Println(path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		log.Println(filepath.Join(h.staticPath, h.indexPath))
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

// SPA is the main function for running the HTTP server
func SPA() {
	hub := newHub()
	go hub.run()

	router := mux.NewRouter()

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	spa := spaHandler{staticPath: "ui/public", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8088",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go hnFeed()
	go imgFlipFeed()

	log.Fatal(srv.ListenAndServe())

}
