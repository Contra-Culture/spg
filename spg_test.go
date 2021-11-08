package spg_test

import (
	"strings"

	. "github.com/Contra-Culture/spg"
	"github.com/Contra-Culture/spg/data"
	"github.com/Contra-Culture/spg/gennode"
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
					cfg.DataGraph(
						func(cfg *data.GraphCfgr) {
							cfg.Schema(
								"author",
								func(cfg *data.SchemaCfgr) {
									cfg.Attribute("full-name")
									cfg.Attribute("slug")
									cfg.Attribute("bio")
									cfg.HasMany(
										"publications",
										func(cfg *data.ArrowCfgr) {
											cfg.Schema(
												"publication",
												func(self *data.Object, attrs map[string]interface{}) bool {
													return true
												})
										})
								})
							cfg.Schema(
								"publication",
								func(cfg *data.SchemaCfgr) {
									cfg.PK([]string{"slug"})
									cfg.Attribute("title")
									cfg.Attribute("slug")
									cfg.Attribute("rubric-slug")
									cfg.Attribute("content")
									cfg.Attribute("author-name")
									cfg.Attribute("published-at")
									cfg.BelongsTo(
										"author",
										func(cfg *data.ArrowCfgr) {
											cfg.Schema(
												"author",
												func(self *data.Object, attrs map[string]interface{}) bool {
													switch authorName := attrs["full-name"].(type) {
													case string:
														return self.Attr("author-name") == authorName
													default:
														return false
													}
												})
										})
									cfg.BelongsTo(
										"rubric",
										func(cfg *data.ArrowCfgr) {
											cfg.Schema(
												"rubric",
												func(self *data.Object, attrs map[string]interface{}) bool {
													switch slug := attrs["slug"].(type) {
													case string:
														return self.Attr("rubric-slug") == slug
													default:
														return false
													}
												})
										})
								})
							cfg.Schema(
								"rubric",
								func(cfg *data.SchemaCfgr) {
									cfg.Attribute("title")
									cfg.HasMany(
										"publications",
										func(cfg *data.ArrowCfgr) {
											cfg.Schema(
												"publication",
												func(self *data.Object, attrs map[string]interface{}) bool {
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
						func(cfg *gennode.NodeCfgr) {
							cfg.Schema(
								"publication",
								[]string{"rubric", "author"},
							)
							cfg.Layout([]string{})
							cfg.Screen([]string{})
							cfg.Node(
								"rubric",
								func(cfg *gennode.NodeCfgr) {
									cfg.Schema(
										"rubric",
										[]string{})
									cfg.Path(
										func(o *data.Object) []string {
											title := o.Attr("title")
											title = strings.NewReplacer(" ", "-", "&", "-and-", "?", "").Replace(title)
											return []string{"rubric", title}
										})
								})
							cfg.Node(
								"publication",
								func(cfg *gennode.NodeCfgr) {
									cfg.Schema(
										"publication",
										[]string{"rubric", "author"})
									cfg.Path(
										func(o *data.Object) []string {
											title := o.Attr("title")
											title = strings.NewReplacer(" ", "-", "&", "-and-", "?", "").Replace(title)
											return []string{"publication", title}
										})
								})
						})
				})).NotTo(BeNil())
	})
})
