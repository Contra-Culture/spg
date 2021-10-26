package spg_test

import (
	"github.com/Contra-Culture/go2html"
	. "github.com/Contra-Culture/spg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("spg", func() {
	It("works", func() {
		Expect(New("test", "example.com", func(cfg *HostCfgr) {
			cfg.Schema(
				"publication",
				func(cfg *SchemaCfgr) {
					cfg.Attribute("title")
					cfg.Attribute("content")
					cfg.Attribute("published-at")
					cfg.Attribute("author")
					cfg.HasMany(
						"comments",
						func(cfg *ArrowCfgr) {
							cfg.Schema(
								"comment",
								func(attrs map[string]interface{}) bool {
									return true
								})
						})
				},
			)
			cfg.Schema(
				"comment",
				func(cfg *SchemaCfgr) {
					cfg.Attribute("content")
					cfg.Attribute("published-at")
					cfg.Attribute("author")
				},
			)
			cfg.Root(func(cfg *PageGeneratorCfgr) {
				cfg.Schema(
					"publication",
					[]string{"comments"},
					func(cfg *go2html.TemplateConfiguringProxy) {

					})
				cfg.Layout(func(cfg *go2html.TemplateConfiguringProxy) {

				})
				cfg.Screen(func(cfg *go2html.TemplateConfiguringProxy) {

				})
				cfg.Page(func(cfg *PageGeneratorCfgr) {

				})
			})
		})).NotTo(BeNil())
	})
})
