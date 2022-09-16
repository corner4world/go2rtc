package streams

import (
	"github.com/AlexxIT/go2rtc/cmd/app"
	"github.com/AlexxIT/go2rtc/cmd/app/store"
	"github.com/rs/zerolog"
	"strings"
)

func Init() {
	var cfg struct {
		Mod map[string]interface{} `yaml:"streams"`
	}

	app.LoadConfig(&cfg)

	log = app.GetLogger("streams")

	for name, item := range cfg.Mod {
		streams[name] = NewStream(item)
	}

	for name, item := range store.GetDict("streams") {
		streams[name] = NewStream(item)
	}
}

func Get(src string) *Stream {
	if stream, ok := streams[src]; ok {
		return stream
	}

	if !HasProducer(src) {
		return nil
	}

	log.Info().Str("url", src).Msg("[streams] create new stream")
	stream := NewStream(src)
	streams[src] = stream
	return stream
}

func Has(src string) bool {
	return streams[src] != nil
}

func New(name string, source interface{}) {
	switch source := source.(type) {
	case string:
		// check if new stream already link on our other stream
		if strings.HasPrefix(source, "rtsp://") {
			if i := strings.IndexByte(source[7:], '/'); i > 0 {
				if stream, ok := streams[source[8+i:]]; ok {
					streams[name] = stream
					return
				}
			}
		}
	}

	streams[name] = NewStream(source)
}

func Delete(name string) {
	delete(streams, name)
}

func All() map[string]interface{} {
	all := map[string]interface{}{}
	for name, stream := range streams {
		all[name] = stream
		//if stream.Active() {
		//	all[name] = stream
		//}
	}
	return all
}

var log zerolog.Logger
var streams = map[string]*Stream{}
