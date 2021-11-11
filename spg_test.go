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
									cfg.Arrow(
										"publications",
										"publication",
										0,
										map[string]string{})
								})
							cfg.Schema(
								"publication",
								func(cfg *data.SchemaCfgr) {
									cfg.ID([]string{"slug"})
									cfg.Attribute("title")
									cfg.Attribute("slug")
									cfg.Attribute("rubric-slug")
									cfg.Attribute("content")
									cfg.Attribute("author-name")
									cfg.Attribute("published-at")
									cfg.Arrow(
										"authors",
										"user",
										0,
										map[string]string{})
									cfg.Arrow(
										"rubric",
										"rubric",
										1,
										map[string]string{})
								})
							cfg.Schema(
								"rubric",
								func(cfg *data.SchemaCfgr) {
									cfg.Attribute("title")
									cfg.Arrow(
										"publications",
										"publication",
										0,
										map[string]string{})
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
