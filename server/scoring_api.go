package server

import (
	"encoding/json"
	"github.com/clemens97/scion-path-oracle"
	"github.com/clemens97/scion-path-oracle/services"
	"github.com/scionproto/scion/go/lib/addr"
	"net/http"
)

type ScoringQuery struct {
	Queries map[string][]services.ServiceName `json:"scorings"`
}

type FingerprintScores struct {
	Fingerprint oracle.PathFingerprint `json:"fingerprint"`
	Scores      map[string]float64     `json:"scores"`
}

type ScoringResponse map[addr.IA][]FingerprintScores

const (
	availableScoringEndpoint = "/services/"
	scoringEndpoint          = "/scorings/"
)

func (o *OracleServer) RegisterScoringAPI() {
	o.router.HandleFunc(availableScoringEndpoint, func(w http.ResponseWriter, r *http.Request) {
		o.handleGetSupportedScorings(w, r)
	}).Methods(http.MethodGet)

	o.router.HandleFunc(scoringEndpoint, func(w http.ResponseWriter, r *http.Request) {
		o.handleQueryScoring(w, r)
	}).Methods(http.MethodPost)
}

func (o *OracleServer) handleGetSupportedScorings(writer http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	supportedScorings := make([]services.ServiceName, len(o.scoringServices))

	i := 0
	for k := range o.scoringServices {
		supportedScorings[i] = k
		i++
	}

	err := json.NewEncoder(writer).Encode(supportedScorings)
	if err != nil {
		http.Error(writer, "", http.StatusInternalServerError)
		return
	}
	setJsonContentType(writer)
	writer.WriteHeader(http.StatusOK)
}

func (o *OracleServer) handleQueryScoring(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	sr, err := parseScoringRequest(request)
	if err != nil {
		o.logger.Debugw("could not parse scoring query", "error", err)
		http.Error(writer, "", http.StatusBadRequest)
		return
	}
	o.logger.Debugw("handling score query", "query", sr)

	response := o.getScorings(sr)

	setJsonContentType(writer)
	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		http.Error(writer, "", http.StatusInternalServerError)
	}
}

func setJsonContentType(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}

func (o *OracleServer) getScorings(req *ScoringQuery) ScoringResponse {
	response := make(ScoringResponse, len(req.Queries))
	for qDst, qServices := range req.Queries {
		dst, err := addr.IAFromString(qDst)
		if err != nil {
			o.logger.Debugw("couldnt parse destination ia from scoring query",
				"dst_ia", qDst, "error", err)
		}

		tmp := make(map[oracle.PathFingerprint]map[string]float64)

		for _, scoringName := range qServices {
			service := *o.scoringServices[scoringName]
			if service == nil {
				o.logger.Debugw("got request for unsupported scoring", "scoring", scoringName)
				break
			}
			for fp, score := range service.GetScores(dst) {
				if tmp[fp] == nil {
					tmp[fp] = make(map[string]float64)
				}
				tmp[fp][string(scoringName)] = score
			}
		}

		for fp, scores := range tmp {
			response[dst] = append(response[dst], FingerprintScores{
				Fingerprint: fp,
				Scores:      scores,
			})
		}

	}
	return response
}

func parseScoringRequest(request *http.Request) (*ScoringQuery, error) {
	var sr ScoringQuery
	if err := json.NewDecoder(request.Body).Decode(&sr); err != nil {
		return nil, err
	}
	return &sr, nil
}
