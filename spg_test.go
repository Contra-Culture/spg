package spg_test

import (
	"github.com/Contra-Culture/go2html"
	. "github.com/Contra-Culture/spg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("spg", func() {
	It("works", func() {
		Expect(
			New(
				"test",
				"example.com",
				func(cfg *HostCfgr) {
					cfg.Repo(
						func(cfg *RepoCfgr) {
							cfg.Schema(
								"publication",
								func(cfg *SchemaCfgr) {
									cfg.PK([]string{"slug"})
									cfg.Attribute("title")
									cfg.Attribute("slug")
									cfg.Attribute("rubric-slug")
									cfg.Attribute("content")
									cfg.Attribute("published-at")
									cfg.Attribute("author")
									cfg.BelongsTo(
										"rubric",
										func(cfg *ArrowCfgr) {
											cfg.Schema(
												"rubric",
												func(self *Object, attrs map[string]interface{}) bool {
													return true
												})
										})
								})
							cfg.Schema(
								"rubric",
								func(cfg *SchemaCfgr) {
									cfg.Attribute("title")
									cfg.HasMany(
										"publications",
										func(cfg *ArrowCfgr) {
											cfg.Schema(
												"publication",
												func(self *Object, attrs map[string]interface{}) bool {
													switch slug := attrs["rubric-slug"].(type) {
													case string:
														return self.Attr("slug") == slug
													default:
														return false
													}
												})
										})
								})
						})
					cfg.Root(
						func(cfg *PageGeneratorCfgr) {
							cfg.Schema(
								"publication",
								[]string{"comments"},
								func(cfg *go2html.TemplateConfiguringProxy) {

								})
							cfg.Layout(func(cfg *go2html.TemplateConfiguringProxy) {

							})
							cfg.Screen(func(cfg *go2html.TemplateConfiguringProxy) {

							})
							cfg.PageGenerator(
								"rubric",
								func(cfg *PageGeneratorCfgr) {
									cfg.Schema(
										"rubric",
										[]string{},
										func(*go2html.TemplateConfiguringProxy) {

										},
									)
								})
						})
				})).NotTo(BeNil())
	})
})
