package server

import (
	"fmt"
	"strings"
)

const unknownVersionKey = "ver"

func parseVersionHeader(versionHdr []string) (map[string]string, error) {
	vmap := make(map[string]string)

	for _, ver := range versionHdr {
		// version header doesn't have key:value pairs
		if !strings.Contains(ver, ":") {
			vmap[unknownVersionKey] = ver
			continue
		}

		for _, kv := range strings.Split(ver, " ") {
			n := strings.SplitN(kv, ":", 2)

			if len(n) != 2 {
				return nil, fmt.Errorf("failed to parse version: %s", kv)
			}

			vmap[n[0]] = n[1]
		}
	}

	return vmap, nil
}
