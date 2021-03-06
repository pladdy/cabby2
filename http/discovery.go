package http

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	cabby "github.com/pladdy/cabby2"
	log "github.com/sirupsen/logrus"
)

// DiscoveryHandler holds a cabby DiscoveryService
type DiscoveryHandler struct {
	DiscoveryService cabby.DiscoveryService
	Port             int
}

// Get serves a discovery resource
func (h DiscoveryHandler) Get(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{"handler": "DiscoveryHandler"}).Debug("Handler called")

	discovery, err := h.DiscoveryService.Discovery(r.Context())
	if err != nil {
		internalServerError(w, err)
		return
	}

	discovery.Default = insertPort(discovery.Default, h.Port)

	for i := 0; i < len(discovery.APIRoots); i++ {
		discovery.APIRoots[i] = swapPath(discovery.Default, discovery.APIRoots[i])
	}

	if discovery.Title == "" {
		resourceNotFound(w, errors.New("Discovery not defined"))
	} else {
		writeContent(w, cabby.TaxiiContentType, resourceToJSON(discovery))
	}
}

// Post handles post request
func (h DiscoveryHandler) Post(w http.ResponseWriter, r *http.Request) {
	methodNotAllowed(w, errors.New("HTTP Method "+r.Method+" unrecognized"))
}

/* helpers */

func parseURL(rawurl string) *url.URL {
	u, err := url.Parse(rawurl)
	if err != nil {
		log.WithFields(log.Fields{"URL": rawurl}).Warn("Failed to parse URL and insert port")
	}
	return u
}

func insertPort(rawurl string, port int) string {
	u := parseURL(rawurl)

	if u.Port() == "" {
		return u.Scheme + "://" + u.Host + ":" + strconv.Itoa(port) + u.Path
	}
	return u.Scheme + "://" + u.Host + u.Path
}

func swapPath(rawurl, newPath string) string {
	u := parseURL(rawurl)
	return u.Scheme + "://" + u.Host + "/" + newPath
}
