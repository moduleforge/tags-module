package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// handleSubjectTags handles GET /entities/{uuid}/tags.
// Returns all tags whose subject is the given entity UUID, filtered by authz.
func (h *handlers) handleSubjectTags(w http.ResponseWriter, r *http.Request) {
	p, ok := h.d.Principal.FromContext(r.Context())
	if !ok {
		jsonErr(w, http.StatusUnauthorized, "unauthorized", "authentication required")
		return
	}

	subjectUUID, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		jsonErr(w, http.StatusBadRequest, "bad_request", "invalid uuid")
		return
	}

	q := r.URL.Query()

	var purposeFilter *string
	if pStr := q.Get("purpose"); pStr != "" {
		s := pStr
		purposeFilter = &s
	}

	pag := parsePagination(q)

	tags, err := h.d.Services.Tag.ListBySubject(
		r.Context(),
		h.d.CoreQuerier,
		h.d.Services.Querier(),
		*p,
		subjectUUID,
		purposeFilter,
		pag,
	)
	if err != nil {
		writeServiceErr(w, err)
		return
	}

	resp := make([]tagResponse, 0, len(tags))
	for _, t := range tags {
		resp = append(resp, toTagResponse(t))
	}
	jsonOK(w, http.StatusOK, map[string]any{"tags": resp})
}
