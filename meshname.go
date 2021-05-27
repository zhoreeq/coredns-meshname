// Package meshname implements a plugin that returns details about the resolving
// querying it.
package meshname

import (
	"context"

	_meshname "github.com/zhoreeq/meshname/pkg/meshname"

	"github.com/miekg/dns"
)

const name = "meshname"

// Meshname is a plugin that resolves .meshname domains 
type Meshname struct{
	dnsClient *dns.Client
}

// ServeDNS implements the plugin.Handler interface.
func (mn Meshname) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	var remoteLookups = make(map[string][]dns.Question)

	a := new(dns.Msg)
	a.SetReply(r)
	a.Authoritative = true

	for _, q := range r.Question {
		labels := dns.SplitDomainName(q.Name)
		if len(labels) < 2 {
			// s.log.Debugln("Error: invalid domain requested")
			continue
		}
		subDomain := labels[len(labels)-2]


		resolvedAddr, err := _meshname.IPFromDomain(&subDomain)
		if err != nil {
			// s.log.Debugln(err)
			continue
		}

		remoteLookups[resolvedAddr.String()] = append(remoteLookups[resolvedAddr.String()], q)
	}

	for remoteServer, questions := range remoteLookups {
		rm := new(dns.Msg)
		rm.RecursionDesired = true
		rm.Question = questions
		resp, _, err := mn.dnsClient.Exchange(rm, "["+remoteServer+"]:53") // no retries
		if err != nil {
			// s.log.Debugln(err)
			continue
		}
		// s.log.Debugln(resp.String())
		a.Answer = append(a.Answer, resp.Answer...)
		a.Ns = append(a.Ns, resp.Ns...)
		a.Extra = append(a.Extra, resp.Extra...)
	}

	w.WriteMsg(a)

	return 0, nil
}

// Name implements the Handler interface.
func (mn Meshname) Name() string { return name }
