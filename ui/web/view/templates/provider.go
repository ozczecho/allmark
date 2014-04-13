// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"fmt"
	"github.com/andreaskoch/allmark2/model"
	"strings"
	"text/template"
)

const (
	// Template placholders
	ChildTemplatePlaceholder = "@childtemplate"

	// Template names
	MasterTemplateName = "master"
	ErrorTemplateName  = "error"

	SitemapTemplateName        = "sitemap"
	SitemapContentTemplateName = "sitemapcontent"

	XmlSitemapTemplateName        = "xmlsitemap"
	XmlSitemapContentTemplateName = "xmlsitemapcontent"

	RssFeedTemplateName        = "rssfeed"
	RssFeedContentTemplateName = "rssfeedcontent"

	TagmapTemplateName        = "tagmap"
	TagmapContentTemplateName = "tagmapcontent"

	SearchTemplateName        = "search"
	SearchContentTemplateName = "searchcontent"
)

type Provider struct {
	Modified chan bool

	folder    string
	templates map[string]*Template
	cache     map[string]*template.Template
}

func NewProvider(templateFolder string) *Provider {

	// intialize the templates
	templateModified := make(chan bool)
	templates := make(map[string]*Template)

	templates[MasterTemplateName] = NewTemplate(templateFolder, MasterTemplateName, masterTemplate, templateModified)
	templates[ErrorTemplateName] = NewTemplate(templateFolder, ErrorTemplateName, errorTemplate, templateModified)

	templates[TagmapTemplateName] = NewTemplate(templateFolder, TagmapTemplateName, tagmapTemplate, templateModified)
	templates[TagmapContentTemplateName] = NewTemplate(templateFolder, TagmapContentTemplateName, tagmapContentTemplate, templateModified)

	templates[SitemapTemplateName] = NewTemplate(templateFolder, SitemapTemplateName, sitemapTemplate, templateModified)
	templates[SitemapContentTemplateName] = NewTemplate(templateFolder, SitemapContentTemplateName, sitemapContentTemplate, templateModified)

	templates[XmlSitemapTemplateName] = NewTemplate(templateFolder, XmlSitemapTemplateName, xmlSitemapTemplate, templateModified)
	templates[XmlSitemapContentTemplateName] = NewTemplate(templateFolder, XmlSitemapContentTemplateName, xmlSitemapContentTemplate, templateModified)

	templates[RssFeedTemplateName] = NewTemplate(templateFolder, RssFeedTemplateName, rssFeedTemplate, templateModified)
	templates[RssFeedContentTemplateName] = NewTemplate(templateFolder, RssFeedContentTemplateName, rssFeedContentTemplate, templateModified)

	templates[SearchTemplateName] = NewTemplate(templateFolder, SearchTemplateName, searchTemplate, templateModified)
	templates[SearchContentTemplateName] = NewTemplate(templateFolder, SearchContentTemplateName, searchContentTemplate, templateModified)

	templates[model.TypeDocument.String()] = NewTemplate(templateFolder, model.TypeDocument.String(), documentTemplate, templateModified)
	templates[model.TypeLocation.String()] = NewTemplate(templateFolder, model.TypeLocation.String(), locationTemplate, templateModified)
	templates[model.TypeMessage.String()] = NewTemplate(templateFolder, model.TypeMessage.String(), messageTemplate, templateModified)
	templates[model.TypeRepository.String()] = NewTemplate(templateFolder, model.TypeRepository.String(), repositoryTemplate, templateModified)
	templates[model.TypePresentation.String()] = NewTemplate(templateFolder, model.TypePresentation.String(), presentationTemplate, templateModified)

	// create the provider
	provider := &Provider{
		Modified: make(chan bool),

		folder:    templateFolder,
		templates: templates,
		cache:     make(map[string]*template.Template),
	}

	// watch for changes
	go func() {
		for {
			select {
			case <-templateModified:
				provider.ClearCache()
				go func() {
					provider.Modified <- true
				}()
			}
		}
	}()

	return provider
}

func (provider *Provider) GetFullTemplate(templateName string) (*template.Template, error) {
	return provider.getParsedTemplate(templateName, true)
}

func (provider *Provider) GetSubTemplate(templateName string) (*template.Template, error) {
	return provider.getParsedTemplate(templateName, false)
}

func (provider *Provider) StoreTemplatesOnDisc() (success bool, err error) {
	for _, template := range provider.templates {
		if savedToDisc, err := template.StoreOnDisc(); !savedToDisc {
			return false, err
		}
	}

	return true, nil
}

func (provider *Provider) ClearCache() {
	fmt.Println("Clearing the template cache.")

	provider.cache = make(map[string]*template.Template)
}

func (provider *Provider) getParsedTemplate(templateName string, includeMaster bool) (*template.Template, error) {

	// get template from cache
	if template, ok := provider.cache[templateName]; ok {
		return template, nil
	}

	// assemble to the template
	childTemplate := provider.getTemplate(templateName)
	if childTemplate == nil {
		return nil, fmt.Errorf("Child template %q not found.", templateName)
	}

	templateText := childTemplate.Text()
	if includeMaster {

		masterTemplate := provider.getTemplate(MasterTemplateName)
		if masterTemplate == nil {
			return nil, fmt.Errorf("Master template not found.")
		}

		// merge master and child template
		templateText = strings.Replace(masterTemplate.Text(), ChildTemplatePlaceholder, templateText, 1)

	}

	// parse the template
	template, err := template.New(templateName).Parse(templateText)
	if err != nil {
		return nil, err
	}

	// add template to cache
	provider.cache[templateName] = template

	return template, nil
}

func (provider *Provider) getTemplate(templateName string) *Template {

	if template, exists := provider.templates[templateName]; exists {
		return template
	}

	return nil
}