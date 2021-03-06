/*
Copyright (C) 2018 Expedia Group.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package execution

import (
	encodingjson "encoding/json"
	"io"
	"net/http"
	"github.com/HotelsDotCom/flyte/flytepath"
	"github.com/HotelsDotCom/flyte/httputil"
	"github.com/HotelsDotCom/flyte/json"
)

func toEvent(pack Pack, reader io.Reader) (*Event, error) {

	event := &Event{Pack: pack}
	if err := encodingjson.NewDecoder(reader).Decode(event); err != nil {
		return nil, err
	}
	return event, nil
}

type actionResponse struct {
	Command string          `json:"command"`
	Input   json.Json       `json:"input"`
	Links   []httputil.Link `json:"links,omitempty"`
}

func toActionResponse(r *http.Request, packId string, action Action) actionResponse {
	return actionResponse{
		Command: action.Name,
		Input:   action.Input,
		Links:   getTakeActionLinks(r, packId, action),
	}
}

func getTakeActionLinks(r *http.Request, packId string, action Action) []httputil.Link {
	link := httputil.Link{Href: httputil.UriBuilder(r).
		Path(flytepath.TakeActionResultPath).
		Replace(":packId", packId).
		Replace(":actionId", action.Id).
		Build(),
		Rel: flytepath.GetUriDocPathFor(flytepath.TakeActionResultDoc)}

	return []httputil.Link{link}
}
