package ui

import (
	"net/http"

	"github.com/SamuelDBines/kubernetes-manager/pkg/store"
	"github.com/SamuelDBines/kubernetes-manager/pkg/web"
)

type IndexPage struct {
	Title      string
	Route      string
	Heading    string
	Subheading string
	Namespaces []store.NamespaceInfo
}

func Index(r *web.Renderer, outDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		nss, err := store.ListNamespaces(outDir)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		r.Render(w, "pages/index", IndexPage{
			Title:      "Namespaces",
			Route:      "index",
			Heading:    "Namespaces",
			Subheading: "Read from out/ and manage generated Kubernetes configs.",
			Namespaces: nss,
		})
	}
}
