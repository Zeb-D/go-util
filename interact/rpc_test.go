package interact

import (
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"testing"
)

var store Store
var (
	listenAddr = flag.String("http", ":8080", "http listen address")
	dataFile   = flag.String("file", "store.gob", "data store file name")
	hostname   = flag.String("host", "localhost:8080", "http host name")
	masterAddr = flag.String("master", "", "RPC master address")
	rpcEnabled = flag.Bool("rpc", true, "enable RPC server")
)

func TestRpcServer(tt *testing.T) {
	flag.Parse()
	if *masterAddr != "" {
		store = NewProxyStore(*masterAddr)
	} else {
		store = NewURLStore(*dataFile)
	}
	if *rpcEnabled { // the master is the rpc server:
		rpc.RegisterName("Store", store)
		rpc.HandleHTTP()
	}
	http.HandleFunc("/", Redirect)
	http.HandleFunc("/add", Add)
	http.ListenAndServe(*listenAddr, nil)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	if key == "" {
		http.NotFound(w, r)
		return
	}
	var url string
	if err := store.Get(&key, &url); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func Add(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		fmt.Fprint(w, AddForm)
		return
	}
	var key string
	if err := store.Put(&url, &key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "http://%s/%s", *hostname, key)
}

const AddForm = `
<form method="POST" action="/add">
URL: <input type="text" name="url">
<input type="submit" value="Add">
</form>
`

const saveQueueLength = 1000

var keyChar = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func genKey(n int) string {
	if n == 0 {
		return string(keyChar[0])
	}
	l := len(keyChar)
	s := make([]byte, 20) // FIXME: will overflow. eventually.
	i := len(s)
	for n > 0 && i >= 0 {
		i--
		j := n % l
		n = (n - j) / l
		s[i] = keyChar[j]
	}
	return string(s[i:])
}

type Store interface {
	Put(url, key *string) error
	Get(key, url *string) error
}

type ProxyStore struct {
	urls   *URLStore // local cache
	client *rpc.Client
}

type URLStore struct {
	urls map[string]string
	mu   sync.RWMutex
	save chan record
}

type record struct {
	Key, URL string
}

func NewURLStore(filename string) *URLStore {
	s := &URLStore{urls: make(map[string]string)}
	if filename != "" {
		s.save = make(chan record, saveQueueLength)
		if err := s.load(filename); err != nil {
			log.Println("Error loading URLStore: ", err)
		}
		go s.saveLoop(filename)
	}
	return s
}

func (s *URLStore) Get(key, url *string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if u, ok := s.urls[*key]; ok {
		*url = u
		return nil
	}
	return errors.New("key not found")
}

func (s *URLStore) Set(key, url *string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, present := s.urls[*key]; present {
		return errors.New("key already exists")
	}
	s.urls[*key] = *url
	return nil
}

func (s *URLStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.urls)
}

func (s *URLStore) Put(url, key *string) error {
	for {
		*key = genKey(s.Count())
		if err := s.Set(key, url); err == nil {
			break
		}
	}
	if s.save != nil {
		s.save <- record{*key, *url}
	}
	return nil
}

func (s *URLStore) load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	d := gob.NewDecoder(f)
	for err == nil {
		var r record
		if err = d.Decode(&r); err == nil {
			s.Set(&r.Key, &r.URL)
		}
	}
	if err == io.EOF {
		return nil
	}
	return err
}

func (s *URLStore) saveLoop(filename string) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error opening URLStore: ", err)
	}
	e := gob.NewEncoder(f)
	for {
		r := <-s.save
		if err := e.Encode(r); err != nil {
			log.Println("Error saving to URLStore: ", err)
		}
	}
}

func NewProxyStore(addr string) *ProxyStore {
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		log.Println("Error constructing ProxyStore: ", err)
		// return // ?
	}
	return &ProxyStore{urls: NewURLStore(""), client: client}
}

func (s *ProxyStore) Get(key, url *string) error {
	if err := s.urls.Get(key, url); err == nil {
		return nil
	}
	// rpc call to master:
	if err := s.client.Call("Store.Get", key, url); err != nil {
		return err
	}
	s.urls.Set(key, url) // update local cache
	return nil
}

func (s *ProxyStore) Put(url, key *string) error {
	// rpc call to master:
	if err := s.client.Call("Store.Put", url, key); err != nil {
		return err
	}
	s.urls.Set(key, url) // update local cache
	return nil
}
