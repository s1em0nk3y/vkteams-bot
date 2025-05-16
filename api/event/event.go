package event

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type EventService struct {
	cli         Client
	pollSeconds uint
}

func New(cli Client, pollSeconds uint) *EventService { return &EventService{cli, pollSeconds} }

func (e *EventService) UpdatesChannel(ctx context.Context) <-chan Event {
	ch := make(chan Event)
	log := zerolog.Ctx(ctx).With().Str("service", "event").Logger()
	log.Info().Msg("Start listen")
	go func() {
		lastEventId := 0
		defer close(ch)
		events, err := e.pollEvents(ctx, lastEventId, int(e.pollSeconds))
		log.Err(err).Int("event_count", len(events)).Msg("Drop unread messages")
		if length := len(events); length > 0 {
			lastEventId = events[length-1].ID
		}
		for {
			select {
			case <-ctx.Done():
				log.Info().Err(ctx.Err()).Msg("context done; exiting")
				return
			default:
				log.Info().Int("event_id", lastEventId).Msg("Fetching events")
				events, err = e.pollEvents(ctx, lastEventId, int(e.pollSeconds))
				if err != nil {
					log.Err(err).Msg("Error occured; sleeping 5s")
					time.Sleep(time.Second * 5)
					continue
				}
				for _, event := range events {
					select {
					case <-ctx.Done():
						log.Info().Err(ctx.Err()).Msg("context done; exiting")
						return
					case ch <- event:
						log.Info().Int("event", event.ID).Msg("event read")
					}
				}
				if length := len(events); length > 0 {
					lastEventId = events[length-1].ID
				}
			}
		}
	}()
	return ch
}

func (e *EventService) pollEvents(ctx context.Context, lastEventID int, pollTime int) ([]Event, error) {
	params := url.Values{
		"lastEventId": {strconv.Itoa(lastEventID)},
		"pollTime":    {strconv.Itoa(pollTime)},
	}

	req, err := e.cli.PerformRequest(ctx, "GET", "/events/get", params, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %w", err)
	}

	resp, err := e.cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to get response: %w", err)
	}
	defer resp.Body.Close()

	response := struct {
		Ok     bool    `json:"ok"`
		Events []Event `json:"events"`
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response.Events, nil
}
