// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"allmark/web/header"
	"allmark/web/orchestrator"
	"allmark/web/view/viewmodel"
	"encoding/json"
	"io"
	"net/http"
)

func TagsSearch(headerWriter header.HeaderWriter, tagsSearchOrchestrator *orchestrator.TagsSearchOrchestrator) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// set headers
		headerWriter.Write(w, header.CONTENTTYPE_JSON)

		// get the suggestions
		tags := tagsSearchOrchestrator.GetTags()
		writeTags(w, tags)
	})

}

func writeTags(writer io.Writer, tags []viewmodel.Tag) error { 
	bytes, err := json.MarshalIndent(tags, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}
