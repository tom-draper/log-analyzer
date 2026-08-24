package analyze

import (
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/tom-draper/log-analyzer/internal/location"
	"github.com/tom-draper/log-analyzer/internal/server"
	"github.com/tom-draper/log-analyzer/pkg/parse"
)

// addLocations resolves the country for any new IP-typed parameter values and
// stores them in locations, keyed by IP string. Existing entries are left
// untouched so repeated IPs are only looked up once.
func addLocations(locations map[string]string, extraction []parse.Extraction) {
	for _, e := range extraction {
		for _, p := range e.Params {
			ipAddress, ok := p.Value.(net.IP)
			if !ok {
				continue
			}
			ipAddressStr := ipAddress.String()
			if _, ok := locations[ipAddressStr]; ok {
				continue
			}
			loc, err := location.GetCountryCode(ipAddress)
			if err != nil {
				continue
			}
			locations[ipAddressStr] = loc
		}
	}
}

func NewData(extraction []parse.Extraction, config *parse.Config) *server.Data {
	locations := make(map[string]string)
	addLocations(locations, extraction)
	return &server.Data{Extraction: extraction, Locations: locations, Config: config}
}

func Run(extraction []parse.Extraction, config *parse.Config, host, port string, indexHTML []byte) {
	data := NewData(extraction, config)
	store := server.NewStore(data)
	server.Start(store, host, port, indexHTML)
}

// RunFollow serves the dashboard and tails the given source (a file path, or
// "-" for stdin), streaming newly extracted records into the dashboard so it
// updates live.
func RunFollow(path string, config *parse.Config, host, port string, maxRecords int, indexHTML []byte) {
	data := &server.Data{
		Extraction: []parse.Extraction{},
		Locations:  map[string]string{},
		Config:     config,
	}
	store := server.NewStore(data)

	go tail(path, config, maxRecords, store)

	server.Start(store, host, port, indexHTML)
}

// tail streams records from the source into the store. Incoming records are
// buffered and flushed on a fixed interval so that a burst of new lines results
// in at most a couple of dashboard refreshes per second rather than one per line.
func tail(path string, config *parse.Config, maxRecords int, store *server.Store) {
	var reader = os.Stdin
	follow := true
	if path != "-" {
		file, err := os.Open(path)
		if err != nil {
			log.Printf("follow: %s", err)
			return
		}
		defer file.Close()
		reader = file
	} else {
		// A pipe streams line-by-line and ends at EOF.
		follow = false
	}

	var (
		mu      sync.Mutex
		pending []parse.Extraction
	)
	emit := func(e parse.Extraction) error {
		mu.Lock()
		pending = append(pending, e)
		mu.Unlock()
		return nil
	}

	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				flush(store, &mu, &pending, maxRecords)
				return
			case <-ticker.C:
				flush(store, &mu, &pending, maxRecords)
			}
		}
	}()

	if err := parse.ParseStream(reader, config, follow, emit); err != nil {
		log.Printf("follow: %s", err)
	}
	close(done)
}

func flush(store *server.Store, mu *sync.Mutex, pending *[]parse.Extraction, maxRecords int) {
	mu.Lock()
	batch := *pending
	*pending = nil
	mu.Unlock()
	if len(batch) == 0 {
		return
	}

	// Resolve locations for the batch outside the store lock, then merge.
	batchLocations := make(map[string]string)
	addLocations(batchLocations, batch)

	store.Update(func(d *server.Data) {
		d.Extraction = append(d.Extraction, batch...)
		d.Extraction = retainRecent(d.Extraction, maxRecords)
		for ip, loc := range batchLocations {
			d.Locations[ip] = loc
		}
		d.Locations = retainedLocations(d.Extraction, d.Locations)
	})
}

// retainRecent bounds the live dashboard's working set. A limit of zero keeps
// all records, which is useful when the caller has its own resource controls.
func retainRecent(extraction []parse.Extraction, maxRecords int) []parse.Extraction {
	if maxRecords <= 0 || len(extraction) <= maxRecords {
		return extraction
	}
	return append([]parse.Extraction(nil), extraction[len(extraction)-maxRecords:]...)
}

// retainedLocations drops country lookups for records that have aged out of
// the dashboard window, keeping location memory proportional to the window.
func retainedLocations(extraction []parse.Extraction, locations map[string]string) map[string]string {
	retained := make(map[string]string)
	for _, e := range extraction {
		for _, p := range e.Params {
			ip, ok := p.Value.(net.IP)
			if !ok {
				continue
			}
			if country, ok := locations[ip.String()]; ok {
				retained[ip.String()] = country
			}
		}
	}
	return retained
}
