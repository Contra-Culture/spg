package spg_test

import (
	"fmt"
	"os"
	"strings"

	. "github.com/Contra-Culture/spg"
	"github.com/Contra-Culture/spg/data"
	"github.com/Contra-Culture/spg/node"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("spg", func() {
	It("works", func() {
		path, err := os.Executable()
		if err != nil {
			panic(err)
		}
		chunks := strings.Split(path, "/")
		path = strings.Join(chunks[:len(chunks)-1], "/")
		Expect(
			New(
				"test",
				"example.com",
				func(cfg *HostCfgr) {
					cfg.DataGraph(
						fmt.Sprintf("%s/test/data", path),
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
										func(cfgr *data.ArrowCfgr) {

										})
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
										func(cfgr *data.ArrowCfgr) {

										})
									cfg.Arrow(
										"rubric",
										"rubric",
										func(cfgr *data.ArrowCfgr) {

										})
								})
							cfg.Schema(
								"rubric",
								func(cfg *data.SchemaCfgr) {
									cfg.Attribute("title")
									cfg.Arrow(
										"publications",
										"publication",
										func(cfgr *data.ArrowCfgr) {

										})
								})
						})
					cfg.Root(
						func(cfg *node.NodeCfgr) {
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
								func(cfg *node.NodeCfgr) {
									cfg.Schema(
										"rubric",
										map[string]interface{}{})
									cfg.Path(
										func(o *data.Object) []string {
											title, err := o.Prop("title")
											if err != nil {
												panic(err)
											}
											title = strings.NewReplacer(" ", "-", "&", "-and-", "?", "").Replace(title)
											return []string{"rubric", title}
										})
								})
							cfg.Node(
								"publication",
								func(cfg *node.NodeCfgr) {
									cfg.Schema(
										"publication",
										map[string]interface{}{
											"rubric": nil,
											"author": nil,
										})
									cfg.Path(
										func(o *data.Object) []string {
											title, err := o.Prop("title")
											if err != nil {
												panic(err)
											}
											title = strings.NewReplacer(" ", "-", "&", "-and-", "?", "").Replace(title)
											return []string{"publication", title}
										})
								})
						})
				})).NotTo(BeNil())
	})
})
