// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"allmark/common/route"
	"allmark/web/view/viewmodel"
	"net/url"
)

type TagsSearchOrchestrator struct {
	*Orchestrator

	tags     []viewmodel.Tag
}

func (orchestrator *TagsSearchOrchestrator) GetTags() []viewmodel.Tag {

	rootItem := orchestrator.rootItem()
	if rootItem == nil {
		orchestrator.logger.Fatal("No root item found")
	}

	if orchestrator.tags == nil {
		orchestrator.SetTagsCache()
	}

	// tagModels := make([]viewmodel.Tag, 0)
	// for tag, item := range orchestrator.tags {

	// 	tagModels = append(tagModels, viewmodel.Tag{
	// 		Name:  item.Name,
	// 		Anchor:   url.QueryEscape(tag),
	// 		Route:    orchestrator.tagPather().Path(url.QueryEscape(tag)),
	// 	})
	// }

	// return tagModels

	return orchestrator.tags
}

func (orchestrator *TagsSearchOrchestrator) SetTagsCache() []viewmodel.Tag {

	// updateTags creates a tags list and assigns it to the orchestrator cache.
	updateTags := func(route route.Route) {

		rootItem := orchestrator.rootItem()
		if rootItem == nil {
			orchestrator.logger.Fatal("No root item found")
		}

		// items by tag
		itemsByTag := make(map[string][]viewmodel.Model)
		for _, item := range orchestrator.getAllItems() {

			itemViewModel := viewmodel.Model{
				Base: getBaseModel(rootItem, item, orchestrator.config),
			}

			for _, tag := range item.MetaData.Tags {
				if items, exists := itemsByTag[tag]; exists {
					itemsByTag[tag] = append(items, itemViewModel)
				} else {
					itemsByTag[tag] = []viewmodel.Model{itemViewModel}
				}
			}

		}

		// create tag models
		tags := make([]viewmodel.Tag, 0)
		for tag, items := range itemsByTag {

			// create view model
			tagModel := viewmodel.Tag{
				Name:     tag,
				Anchor:   url.QueryEscape(tag),
				Route:    orchestrator.tagPather().Path(url.QueryEscape(tag)),
				Children: items,
			}

			// append to list
			tags = append(tags, tagModel)
		}

		// sort the tags
		viewmodel.SortTagBy(tagsByName).Sort(tags)

		orchestrator.tags = tags
	}

	asyncUpdate := func(route route.Route) {
		go updateTags(route)
	}

	// register update callbacks
	orchestrator.registerUpdateCallback("update tags", UpdateTypeNew, asyncUpdate)
	orchestrator.registerUpdateCallback("update tags", UpdateTypeModified, asyncUpdate)
	orchestrator.registerUpdateCallback("update tags", UpdateTypeDeleted, asyncUpdate)

	// build the cache
	updateTags(route.New())

	return orchestrator.tags
}
