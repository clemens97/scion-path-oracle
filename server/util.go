package server

import (
	"github.com/scionproto/scion/go/lib/addr"
	"strconv"
)

func parseIAFromNumeric(isd, as string) (addr.IA, error) {
	i, err := strconv.ParseInt(isd, 10, 16)
	if err != nil {
		return addr.IA{}, err
	}

	a, err := strconv.ParseInt(as, 10, 64)
	if err != nil {
		return addr.IA{}, err
	}
	return addr.IA{
		I: addr.ISD(i),
		A: addr.AS(a),
	}, nil
}
