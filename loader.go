package conflate

import (
	"net"
	"net/http"
	"os"
	"time"
)

var getwd = os.Getwd

type loader struct {
	newFiledata func([]byte, include) (filedata, error)
}

func (l *loader) loadURLsRecursive(parentIncs includes, incs ...include) (filedatas, error) {
	var allData filedatas
	for _, inc := range incs {
		data, err := l.loadURLRecursive(parentIncs, inc)
		if err != nil {
			return nil, err
		}
		allData = append(allData, data...)
	}
	return allData, nil
}

func (l *loader) loadURLRecursive(parentIncs includes, inc include) (filedatas, error) {
	data, err := inc.load()
	if err != nil {
		return nil, err
	}
	fd, err := l.newFiledata(data, inc)
	if err != nil {
		return nil, err
	}
	return l.loadDatumRecursive(parentIncs, inc, fd)
}

func (l *loader) loadDataRecursive(parentIncs includes, data ...filedata) (filedatas, error) {
	var allData filedatas
	for _, datum := range data {
		childData, err := l.loadDatumRecursive(parentIncs, emptyInclude, datum)
		if err != nil {
			return nil, err
		}
		allData = append(allData, childData...)
	}
	return allData, nil
}

func (l *loader) loadDatumRecursive(parentIncs includes, inc include, fd filedata) (filedatas, error) {
	if fd.isEmpty() {
		return nil, nil
	}
	if parentIncs.contains(inc) {
		return nil, makeError("The url recursively includes itself (%v)", inc.URL)
	}
	var newParentIncs includes
	newParentIncs = append(newParentIncs, parentIncs...)
	if !inc.isEmpty() {
		newParentIncs = append(newParentIncs, inc)
	}
	childFds, err := l.loadURLsRecursive(newParentIncs, fd.includes...)
	if err != nil {
		return nil, err
	}
	var fds filedatas
	fds = append(fds, childFds...)
	fds = append(fds, fd)
	return fds, nil
}

func (l *loader) wrapFiledata(bytes []byte) (filedata, error) {
	return l.newFiledata(bytes, emptyInclude)
}

func (l *loader) wrapFiledatas(bytes ...[]byte) (filedatas, error) {
	var fds []filedata
	for _, b := range bytes {
		fd, err := l.wrapFiledata(b)
		if err != nil {
			return nil, err
		}
		fds = append(fds, fd)
	}
	return fds, nil
}

func newTransport() *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	return transport
}
