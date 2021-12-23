package main

import (
	"context"
	"encoding/json"
	"github.com/vineet-garg/assignment-1/config"
	"github.com/vineet-garg/assignment-1/stats"
	"github.com/vineet-garg/assignment-1/store"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

var hashPost = regexp.MustCompile("^/hash$")
var hashGet = regexp.MustCompile("^/hash/[0-9]+$")
var statsGet = regexp.MustCompile("^/stats$")
var shutDownPost = regexp.MustCompile("^/shutdown$")

var srv = http.Server{
	// TODO SHIFT to HTTPS in production
	Addr:    config.Addr,
	Handler: hashAPIHandler{},
}

type hashAPIHandler struct {
}

var stop int32 = 0

var wg = sync.WaitGroup{}

func (h hashAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&stop) != 0 {
		http.Error(w, "503 Service Unavailable.", http.StatusServiceUnavailable)
		return
	}
	switch {
	case hashGet.MatchString(r.URL.Path):
		switch r.Method {
		case "GET":
			segments := strings.Split(r.URL.Path, "/")
			id, _ := strconv.ParseInt(segments[len(segments)-1], 10, 64)
			if hash, ok := store.GetStore().GetHash(id); ok {
				w.Write([]byte(hash))
				return
			}
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		default:
			http.Error(w, "405 Method Not Allowed.", http.StatusMethodNotAllowed)
			return
		}
	case hashPost.MatchString(r.URL.Path):
		switch r.Method {
		case "POST":
			if err := r.ParseForm(); err != nil {
				log.Printf("Error : %v", err)
				http.Error(w, "400 Bad Request.", http.StatusBadRequest)
				return
			}
			pwd := r.FormValue("password")
			if pwd == "" {
				http.Error(w, "400 Bad Request.", http.StatusBadRequest)
				return
			}

			// TODO MUST validate the input size to mitigate DOS and bufferoverflow attacks.

			var id int64
			start := time.Now()
			wg.Add(1)
			id = store.GetStore().AddPwd([]byte(pwd), &wg)
			stats.GetStatsStore().Update(time.Since(start))
			w.Header().Set("Content-Type", "application/json")
			jsonResp, err := json.Marshal(id)
			if err != nil {
				log.Printf("Error : %v", err)
				http.Error(w, "500 Server Error.", http.StatusInternalServerError)
				return
			}
			w.Write(jsonResp)
			return
		default:
			http.Error(w, "405 Method Not Allowed.", http.StatusMethodNotAllowed)
			return
		}
	case statsGet.MatchString(r.URL.Path):
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			resp := make(map[string]int64)
			count, avg := stats.GetStatsStore().Get()
			resp["total"] = count
			resp["average"] = avg
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Printf("Error : %v", err)
				http.Error(w, "500 Server Error.", http.StatusInternalServerError)
				return
			}
			w.Write(jsonResp)
			return
		default:
			http.Error(w, "405 Method Not Allowed.", http.StatusMethodNotAllowed)
			return
		}
	case shutDownPost.MatchString(r.URL.Path):
		switch r.Method {
		case "POST":
			atomic.StoreInt32(&stop, 1)
			// wait for all delayed hashes to be added
			wg.Wait()
			go func() {
				log.Printf("shutting down server gracefully")
				if err := srv.Shutdown(context.Background()); err != nil {
					log.Printf("HTTP server Shutdown: %v", err)
				}
				os.Exit(1)
			}()
			return
		default:
			http.Error(w, "405 Method Not Allowed.", http.StatusMethodNotAllowed)
			return
		}
	default:
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
}

func main() {
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)

		<-sigint
		log.Printf("shutting down server gracefully")
		// received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()
	log.Printf("starting server at localhost:8080")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnsClosed
}