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
									cfg.Association(
										"publications",
										"publication",
										"",
										0,
										nil,
										func(self *data.Object, attrs map[string]interface{}) bool {
											return true
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
									cfg.Association(
										"authors",
										"user",
										"",
										0,
										nil,
										func(self *data.Object, attrs map[string]interface{}) bool {
											switch authorName := attrs["full-name"].(type) {
											case string:
												return self.Attr("author-name") == authorName
											default:
												return false
											}
										})
									cfg.Association(
										"rubric",
										"rubric",
										"",
										1,
										nil,
										func(self *data.Object, attrs map[string]interface{}) bool {
											switch slug := attrs["slug"].(type) {
											case string:
												return self.Attr("rubric-slug") == slug
											default:
												return false
											}
										})
								})
							cfg.Schema(
								"rubric",
								func(cfg *data.SchemaCfgr) {
									cfg.Attribute("title")
									cfg.Association(
										"publications",
										"publication",
										"",
										0,
										nil,
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
					cfg.Root(
						func(cfg *gennode.NodeCfgr) {
							cfg.Schema(
								"publication",
								map[string]interface{}{
									"rubric": nil,
									"author": nil,
								},
							)
							cfg.Layout([]string{})
							cfg.Screen([]string{})
							cfg.Node(
								"rubric",
								func(cfg *gennode.NodeCfgr) {
									cfg.Schema(
										"rubric",
										map[string]interface{}{})
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
										map[string]interface{}{
											"rubric": nil,
											"author": nil,
										})
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
