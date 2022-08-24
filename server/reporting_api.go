package server

import (
	"encoding/json"
	"github.com/clemens97/scion-path-oracle"
	"github.com/gorilla/mux"
	"github.com/scionproto/scion/go/lib/snet"
	"net/http"
)

const reportEndpoint = "/reports/{" + dstIsdKey + "}/{" + dstAsKey + "}/{" + pathKey + "}/"

func (o *OracleServer) RegisterStatsAPI() {
	o.router.HandleFunc(reportEndpoint, func(w http.ResponseWriter, r *http.Request) {
		o.handleReport(w, r)
	}).Methods("POST")
}

func (o *OracleServer) handleReport(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	report, err := parseReport(request)
	if err != nil {
		o.logger.Infow("error parsing Report from request", "error", err)
		http.Error(writer, "", http.StatusBadRequest)
		return
	}

	if local := o.querier.LocalIA(); local != report.SrcIA {
		badReq(writer)
		o.logger.Infow("received report from AS other than the local oracle AS - dropping report",
			"local_ia", local, "report_origin_ia", report.SrcIA)
		return
	}

	path, err := o.querier.GetPath(report.DstIA, report.PathFp)

	if err != nil || path == nil {
		writer.WriteHeader(http.StatusForbidden)
		o.logger.Infow("received Report for non existing path - dropping report",
			"fingerprint", report.PathFp)
		return
	}

	if o.scoringServices == nil {
		writer.WriteHeader(http.StatusInternalServerError)
		o.logger.Infow("received Report but no scoring services are registered")
		return
	}

	for _, s := range o.scoringServices {
		o.logger.Debugw("forwarding stats report to scoring service",
			"report", *report, "service", (*s).Name())
		(*(s)).OnStatsReceived(*report, path)
	}
	o.logger.Infow("passed report to all registered scoring services")
	writer.WriteHeader(http.StatusCreated)
}

func badReq(writer http.ResponseWriter) {
	http.Error(writer, "", http.StatusBadRequest)
}

func parseReport(request *http.Request) (*oracle.Report, error) {
	var report oracle.Report
	// parse stats from HTTP body
	if err := json.NewDecoder(request.Body).Decode(&report); err != nil {
		return nil, err
	}

	// path parameters
	vars := mux.Vars(request)
	report.PathFp = oracle.PathFingerprint(vars[pathKey])
	if ia, err := parseIAFromNumeric(vars[dstIsdKey], vars[dstAsKey]); err != nil {
		return nil, err
	} else {
		report.DstIA = ia
	}

	if hostAddr, err := snet.ParseUDPAddr(request.RemoteAddr); err != nil {
		return nil, err
	} else {
		report.SrcIA = hostAddr.IA
	}

	return &report, nil
}
