package gomeassistant

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/saml-dev/gome-assistant/internal"
	"github.com/saml-dev/gome-assistant/internal/http"
	pq "github.com/saml-dev/gome-assistant/internal/priorityqueue"
	ws "github.com/saml-dev/gome-assistant/internal/websocket"
)

type app struct {
	ctx        context.Context
	ctxCancel  context.CancelFunc
	conn       *websocket.Conn
	httpClient *http.HttpClient

	service *Service
	state   *State

	schedules         pq.PriorityQueue
	entityListeners   map[string][]entityListener
	entityListenerIds map[int64]entityListenerCallback
}

/*
NewApp establishes the websocket connection and returns an object
you can use to register schedules and listeners.
*/
func NewApp(connString string) app {
	token := os.Getenv("HA_AUTH_TOKEN")
	conn, ctx, ctxCancel := ws.SetupConnection(connString, token)

	httpClient := http.NewHttpClient(connString, token)

	service := NewService(conn, ctx, httpClient)
	state := NewState(httpClient)

	return app{
		conn:              conn,
		ctx:               ctx,
		ctxCancel:         ctxCancel,
		httpClient:        httpClient,
		service:           service,
		state:             state,
		schedules:         pq.New(),
		entityListeners:   map[string][]entityListener{},
		entityListenerIds: map[int64]entityListenerCallback{},
	}
}

func (a *app) Cleanup() {
	if a.ctxCancel != nil {
		a.ctxCancel()
	}
}

func (a *app) RegisterSchedule(s schedule) {
	if s.err != nil {
		log.Fatalln(s.err) // something wasn't configured properly when the schedule was built
	}

	if s.frequency == 0 {
		log.Fatalln("A schedule must call either Daily() or Every() when built.")
	}

	// TODO: consider moving all time stuff to carbon?
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()) // start at midnight today

	// apply offset if set
	if s.offset.Minutes() > 0 {
		startTime.Add(s.offset)
	}

	// advance first scheduled time by frequency until it is in the future
	for startTime.Before(now) {
		startTime = startTime.Add(s.frequency)
	}

	s.realStartTime = startTime
	a.schedules.Insert(s, float64(startTime.Unix())) // TODO: this blows up because schedule can't be used as key for map in prio queue lib. Just copy/paste and tweak as needed
}

func (a *app) RegisterEntityListener(el entityListener) {
	for _, entity := range el.entityIds {
		id := internal.GetId()
		subscribeTriggerMsg := subscribeMsg{
			Id:   id,
			Type: "subscribe_trigger",
			Trigger: subscribeMsgTrigger{
				Platform: "state",
				EntityId: entity,
			},
		}
		if el.fromState != "" {
			subscribeTriggerMsg.Trigger.From = el.fromState
		}
		if el.toState != "" {
			subscribeTriggerMsg.Trigger.To = el.toState
		}
		log.Default().Println(subscribeTriggerMsg)
		ws.WriteMessage(subscribeTriggerMsg, a.conn, a.ctx)
		msg, _ := ws.ReadMessage(a.conn, a.ctx)
		log.Default().Println(string(msg))
		a.entityListenerIds[id] = el.callback
	}

}

func (a *app) Start() {
	// schedules
	go RunSchedules(a)

	// entity listeners
	elChan := make(chan ws.ChanMsg)
	go ws.ListenWebsocket(a.conn, a.ctx, elChan)

	log.Default().Println(a.entityListenerIds)
	var msg ws.ChanMsg
	for {
		msg = <-elChan
		log.Default().Println(string(msg.Raw))
		if callback, ok := a.entityListenerIds[msg.Id]; ok {
			log.Default().Println(msg, callback)
		}
	}

	// NOTE:should the prio queue and websocket listener both write to a channel or something?
	// then select from that and spawn new goroutine to call callback?

	// TODO: loop through schedules and create heap priority queue

	// TODO: figure out looping listening to messages for
	// listeners
}

const (
	FrequencyMissing time.Duration = 0

	Daily    time.Duration = time.Hour * 24
	Hourly   time.Duration = time.Hour
	Minutely time.Duration = time.Minute
)
