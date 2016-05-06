package cologappengine

import (
	"io/ioutil"
	"net/http"

	"github.com/kmtr/colog"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// NewCologAppEngine creates new colog logger for AppEngine.
func NewCologAppEngine(w http.ResponseWriter, r *http.Request, prefix string, flag int, lvMap LevelMap) *colog.CoLog {
	l := colog.NewCoLog(w, prefix, flag)
	l.SetOutput(ioutil.Discard)
	h := &cologAppEngineHook{
		ctx:       appengine.NewContext(r),
		formatter: &colog.StdFormatter{},
	}
	if lvMap != nil {
		h.lvMap = lvMap
	} else {
		h.lvMap = defaultCologAppEngineLevelMap
	}
	h.lvs = levelMapKeys(h.lvMap)
	l.AddHook(h)

	return l
}

func levelMapKeys(m LevelMap) []colog.Level {
	keys := make([]colog.Level, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

// AppEngineLogLevel represens severity level in AppEngine
type AppEngineLogLevel uint8

const (
	// AppEngineLDebug represents debug severity level in AppEngine
	AppEngineLDebug = iota
	// AppEngineLInfo represents info severity level in AppEngine
	AppEngineLInfo
	// AppEngineLWarning represents warning severity level in AppEngine
	AppEngineLWarning
	// AppEngineLError represents error severity level in AppEngine
	AppEngineLError
	// AppEngineLCritical represents critical severity level in AppEngine
	AppEngineLCritical
)

// LevelMap is convert map CoLog -> AppEngine log
type LevelMap map[colog.Level]AppEngineLogLevel

var defaultCologAppEngineLevelMap = map[colog.Level]AppEngineLogLevel{
	colog.LTrace:   AppEngineLDebug,
	colog.LDebug:   AppEngineLDebug,
	colog.LInfo:    AppEngineLInfo,
	colog.LWarning: AppEngineLWarning,
	colog.LError:   AppEngineLError,
	colog.LAlert:   AppEngineLCritical,
}

type cologAppEngineHook struct {
	lvs       []colog.Level
	lvMap     LevelMap
	ctx       context.Context
	formatter colog.Formatter
}

// Levels returns the set of levels for which the hook should be triggered
func (h *cologAppEngineHook) Levels() []colog.Level {
	return h.lvs
}

// Fire method converts log level from colog to appengine and puts log.
func (h *cologAppEngineHook) Fire(e *colog.Entry) error {
	lv := defaultCologAppEngineLevelMap[e.Level]
	b, err := h.formatter.Format(e)
	if err != nil {
		return err
	}
	msg := string(b)
	switch lv {
	case AppEngineLDebug:
		log.Debugf(h.ctx, msg)
	case AppEngineLInfo:
		log.Infof(h.ctx, msg)
	case AppEngineLWarning:
		log.Warningf(h.ctx, msg)
	case AppEngineLError:
		log.Errorf(h.ctx, msg)
	case AppEngineLCritical:
		log.Criticalf(h.ctx, msg)
	}
	return nil
}
