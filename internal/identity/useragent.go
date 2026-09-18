package identity

import "math/rand"

var userAgents = []string{
	"TrafficLab/0.1 (Go)",
	"Mozilla/5.0 (compatible; TrafficLab/0.1)",
	"TrafficLab-load-test/0.1",
}

func RandomUserAgent(r *rand.Rand) string {
	return userAgents[r.Intn(len(userAgents))]
}
