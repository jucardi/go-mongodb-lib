package mgo

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/jucardi/go-mongodb-lib/log"
	"gopkg.in/mgo.v2"
)

var (
	// ErrNotFound is the error returned when no results are found in a mongo operation.
	ErrNotFound = mgo.ErrNotFound

	// ErrCursor is the error returned when the cursor used in a mongo operation is not valid.
	ErrCursor = mgo.ErrCursor

	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type ExtraConfig struct {
	// IndependentHosts if set to 'true', will treat the 'Addrs' as a list of db
	// hosts to connect to and use the next(s) as fallback.
	IndependentHosts bool
	// ShuffleHosts indicates if shuffling the list of addresses should be done
	// before dialing.
	ShuffleHosts bool
}

// Dial establishes a new session to the cluster identified by the given seed
// server(s). The session will enable communication with all of the servers in
// the cluster, so the seed servers are used only to find out about the cluster
// topology.
//
//     - See mgo.Dial documentation in `gopkg.in/mgo.v2` for more information.
//
func Dial(url ...string) (ISession, error) {
	if len(url) == 0 {
		return nil, errors.New("expected a 'url'")
	}

	var (
		s   *mgo.Session
		err error
	)

	for _, u := range url {
		s, err = mgo.Dial(u)

		for j := 1; err != nil && j <= DialMaxRetries; j++ {
			log.Get().Error(fmt.Sprintf("Can't connect to mongo on '%s': %v. Retrying in %v", url, err, DialRetrySleep))
			time.Sleep(DialRetrySleep)
			log.Get().Warn(fmt.Sprintf("Retrying to connect to mongo, attempt %d of %d", j, DialMaxRetries))
			s, err = mgo.Dial(u)
		}

		if err == nil {
			break
		}
	}

	return fromSession(s), err
}

// DialWithTimeout works like Dial, but uses timeout as the amount of time to
// wait for a server to respond when first connecting and also on follow up
// operations in the session. If timeout is zero, the call may block
// forever waiting for a connection to be made.
//
// See SetSyncTimeout for customizing the timeout for the session.
func DialWithTimeout(url string, timeout time.Duration) (ISession, error) {
	s, err := mgo.DialWithTimeout(url, timeout)
	return fromSession(s), err
}

// DialWithInfo establishes a new session to the cluster identified by info.
//
//   {extraCfg} - (Optional) additional dialing options.
//
func DialWithInfo(info *mgo.DialInfo, extraCfg ...ExtraConfig) (ISession, error) {
	cfg := ExtraConfig{}

	if len(extraCfg) > 0 {
		cfg = extraCfg[0]
	}

	var (
		s   *mgo.Session
		err error
	)

	if cfg.ShuffleHosts && len(info.Addrs) > 1 {
		for i := len(info.Addrs) - 1; i > 0; i-- { // Fisher–Yates shuffle
			j := rand.Intn(i + 1)
			info.Addrs[i], info.Addrs[j] = info.Addrs[j], info.Addrs[i]
		}
	}

	addrs := info.Addrs

	for _, u := range addrs {
		if cfg.IndependentHosts {
			info.Addrs = []string{u}
		}

		s, err = mgo.DialWithInfo(info)

		for i := 1; err != nil && i <= DialMaxRetries; i++ {
			log.Get().Error(fmt.Sprintf("Can't connect to mongo on '%v': %v. Retrying in %v", info.Addrs, err, DialRetrySleep))
			time.Sleep(DialRetrySleep)
			log.Get().Warn(fmt.Sprintf("Retrying to connect to mongo, attempt %d of %d", i, DialMaxRetries))
			s, err = mgo.DialWithInfo(info)
		}

		if err == nil || !cfg.IndependentHosts {
			break
		}
	}

	return fromSession(s), err
}

// DialWithTls attempts to establish a MongoDB connection using TLS with the provided PEM encoded
// certificate.
//
//    {url}                 - The URL to dial to Mongo
//    {cert}                - The PEM bytes of the CA cert to use for TLS
//    {insecureSkipVerify}  - (optional) controls whether a client verifies the server's certificate
//                            chain andhost name. If 'true', TLS accepts any certificate presented
//                            by the server and any host name in that certificate. In this mode, TLS
//                            is susceptible to man-in-the-middle attacks.  This should be used only
//                            for testing.
//
func DialWithTls(url string, cert []byte, insecureSkipVerify ...bool) (ISession, error) {
	rootCerts := x509.NewCertPool()
	rootCerts.AppendCertsFromPEM(cert)

	var lastErr error

	info, err := mgo.ParseURL(url)
	if err != nil {
		return nil, err
	}
	info.Timeout = DialTimeout
	info.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		tlsCfg := &tls.Config{
			RootCAs: rootCerts,
		}
		if len(insecureSkipVerify) > 0 && insecureSkipVerify[0] {
			tlsCfg.InsecureSkipVerify = insecureSkipVerify[0]
		}
		conn, err := tls.Dial("tcp", addr.String(), tlsCfg)
		if err != nil {
			lastErr = err
		}
		return conn, err
	}

	// Dial with TLS
	s, err := mgo.DialWithInfo(info)

	for i := 1; err != nil && i <= DialMaxRetries; i++ {
		if lastErr != nil {
			err = lastErr
			lastErr = nil
		}
		log.Get().Error(fmt.Sprintf("Can't connect to mongo on '%s': %v. Retrying in %v", url, err, DialRetrySleep))
		time.Sleep(DialRetrySleep)
		log.Get().Warn(fmt.Sprintf("Retrying to connect to mongo, attempt %d of %d", i, DialMaxRetries))
		s, err = mgo.DialWithInfo(info)
	}

	return fromSession(s), err
}

// AddTlsHandler adds the TLS handling logic to a provided `*mgo.DialInfo`
//
//    {dialInfo}            - The DialInfo to use
//    {cert}                - The PEM bytes of the CA cert to use for TLS
//    {insecureSkipVerify}  - Controls whether a client verifies the server's certificate chain and
//                            host name. If 'true', TLS accepts any certificate presented by the
//                            server and any host name in that certificate. In this mode, TLS is
//                            susceptible to man-in-the-middle attacks.  This should be used only
//                            for testing.
//
func AddTlsHandler(dialInfo *mgo.DialInfo, cert []byte, insecureSkipVerify bool) {
	rootCerts := x509.NewCertPool()
	rootCerts.AppendCertsFromPEM(cert)

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), &tls.Config{
			RootCAs:            rootCerts,
			InsecureSkipVerify: insecureSkipVerify,
		})
		if err != nil {
			log.Get().Error("An error has occurred while dialing with TLS, ", err)
		}
		return conn, err
	}
}

// ParseURL parses a MongoDB URL as accepted by the Dial function and returns
// a value suitable for providing into DialWithInfo.
//
// See Dial for more details on the format of url.
func ParseURL(url string) (*mgo.DialInfo, error) {
	return mgo.ParseURL(url)
}

// IsDup returns whether err informs of a duplicate key error because
// a primary key index or a secondary unique index already has an entry
// with the given value.
func IsDup(err error) bool {
	return mgo.IsDup(err)
}
