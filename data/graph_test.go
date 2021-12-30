package data_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/Contra-Culture/report"
	. "github.com/Contra-Culture/spg/data"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("data", func() {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	chunks := strings.Split(path, "/")
	path = strings.Join(append(chunks[:len(chunks)-2], "test/data"), "/")
	Describe("*Graph", func() {
		Describe("creation | New()", func() {
			Context("when valid specification", func() {
				It("returns *Graph", func() {
					r := report.New("graph")
					g := New(r, path, func(cfg *GraphCfgr) {})
					Expect(g).NotTo(BeNil())
					Expect(report.ToString(r)).To(Equal("| graph\n"))
				})
			})
			Describe(".Node() specification", func() {
				Context("when valid specification", func() {
					It("adds Schema specification to the graph", func() {
						r := report.New("graph")
						g := New(r, path, func(cfg *GraphCfgr) {
							cfg.Node("account", func(cfg *NodeCfgr) {
								cfg.PK([]string{"login"})
								cfg.Attribute("login")
								cfg.Attribute("firstName")
								cfg.Attribute("lastName")
								cfg.Attribute("bio")
								cfg.Attribute("publicationsCount")
								cfg.Link("publications", []string{"publication"}, func(cfg *LinkCfgr) {
									cfg.Map("login", "authors")
									cfg.CounterCache("publicationsCount")
								})
							})
							cfg.Node("rubric", func(cfg *NodeCfgr) {
								cfg.PK([]string{"slug"})
								cfg.Attribute("slug")
								cfg.Attribute("title")
								cfg.Attribute("description")
								cfg.Attribute("publicationsCount")
								cfg.Link("publications", []string{"publication"}, func(cfg *LinkCfgr) {
									cfg.Map("slug", "rubricSlug")
									cfg.CounterCache("publicationsCount")
								})
							})
							cfg.Node("publication", func(cfg *NodeCfgr) {
								cfg.PK([]string{"slug"})
								cfg.Attribute("rubricSlug")
								cfg.Attribute("slug")
								cfg.Attribute("title")
								cfg.Attribute("content")
								cfg.Attribute("createdAt")
								cfg.Attribute("updatedAt")
								cfg.Attribute("accountLogin")
								cfg.Link("rubric", []string{"rubric"}, func(cfg *LinkCfgr) {
									cfg.Map("rubricSlug", "slug")
									cfg.Limit(1)
								})
								cfg.Link("authors", []string{"account"}, func(cfg *LinkCfgr) {
									cfg.Map("authors", "login")
								})
							})
						})
						Expect(g).NotTo(BeNil())
						Expect(report.ToString(r)).To(Equal("| graph\n\t| schema: account\n\t| schema: rubric\n\t| schema: publication\n"))
					})
				})
			})
		})
		Describe("objects updating", func() {
			Describe(".Update()", func() {
				Context("when successfull", func() {
					It("returns object id", func() {
						r := report.New("graph")
						g := New(r, path, func(cfg *GraphCfgr) {
							cfg.Node("account", func(cfg *NodeCfgr) {
								cfg.PK([]string{"login"})
								cfg.Attribute("login")
								cfg.Attribute("firstName")
								cfg.Attribute("lastName")
								cfg.Attribute("bio")
								cfg.Attribute("publicationsCount")
								cfg.Link("publications", []string{"publication"}, func(cfg *LinkCfgr) {
									cfg.Map("login", "authors")
									cfg.CounterCache("publicationsCount")
								})
							})
							cfg.Node("rubric", func(cfg *NodeCfgr) {
								cfg.PK([]string{"slug"})
								cfg.Attribute("slug")
								cfg.Attribute("title")
								cfg.Attribute("description")
								cfg.Attribute("publicationsCount")
								cfg.Link("publications", []string{"publication"}, func(cfg *LinkCfgr) {
									cfg.Map("slug", "rubricSlug")
									cfg.CounterCache("publicationsCount")
								})
							})
							cfg.Node("publication", func(cfg *NodeCfgr) {
								cfg.PK([]string{"slug"})
								cfg.Attribute("rubricSlug")
								cfg.Attribute("slug")
								cfg.Attribute("title")
								cfg.Attribute("content")
								cfg.Attribute("createdAt")
								cfg.Attribute("updatedAt")
								cfg.Attribute("accountLogin")
								cfg.Link("rubric", []string{"rubric"}, func(cfg *LinkCfgr) {
									cfg.Map("rubricSlug", "slug")
									cfg.Limit(1)
								})
								cfg.Link("authors", []string{"account"}, func(cfg *LinkCfgr) {
									cfg.Map("authors", "login")
								})
							})
						})
						id, err := g.Update(
							"rubric",
							map[string]string{
								"slug":        "interviews",
								"title":       "Interviews",
								"description": "Interviews with passionate people.",
							})
						Expect(err).NotTo(HaveOccurred())
						Expect(id).To(Equal("interviews"))
						o, err := g.Get("rubric", "interviews")
						Expect(err).NotTo(HaveOccurred())
						Expect(o).NotTo(BeNil())
					})
				})
			})
		})
		Describe("JSON representation", func() {
			It("returns JSON string", func() {
				r := report.New("graph")
				g := New(r, path, func(cfg *GraphCfgr) {
					cfg.Node("account", func(cfg *NodeCfgr) {
						cfg.PK([]string{"login"})
						cfg.Attribute("login")
						cfg.Attribute("firstName")
						cfg.Attribute("lastName")
						cfg.Attribute("bio")
						cfg.Attribute("publicationsCount")
						cfg.Link("publications", []string{"publication"}, func(cfg *LinkCfgr) {
							cfg.Map("login", "authors")
							cfg.CounterCache("publicationsCount")
						})
					})
					cfg.Node("rubric", func(cfg *NodeCfgr) {
						cfg.PK([]string{"slug"})
						cfg.Attribute("slug")
						cfg.Attribute("title")
						cfg.Attribute("description")
						cfg.Attribute("publicationsCount")
						cfg.Link("publications", []string{"publication"}, func(cfg *LinkCfgr) {
							cfg.Map("slug", "rubricSlug")
							cfg.CounterCache("publicationsCount")
						})
					})
					cfg.Node("publication", func(cfg *NodeCfgr) {
						cfg.PK([]string{"slug"})
						cfg.Attribute("rubricSlug")
						cfg.Attribute("slug")
						cfg.Attribute("title")
						cfg.Attribute("content")
						cfg.Attribute("createdAt")
						cfg.Attribute("updatedAt")
						cfg.Attribute("accountLogin")
						cfg.Link("rubric", []string{"rubric"}, func(cfg *LinkCfgr) {
							cfg.Map("rubricSlug", "slug")
							cfg.Limit(1)
						})
						cfg.Link("authors", []string{"account"}, func(cfg *LinkCfgr) {
							cfg.Map("authors", "login")
						})
					})
				})
				fmt.Printf("\n\nDEBUG: %s\n\n", g.JSONString())
				Expect(g.JSONString()).To(Equal("{\"schemas\":{\"account\":{\"name\":\"account\",\"id\":[\"login\",\"login\"],\"attributes\":[\"login\",\"firstName\",\"lastName\",\"bio\"],\"arrows\":[]},\"rubric\":{\"name\":\"rubric\",\"id\":[\"slug\",\"slug\"],\"attributes\":[\"slug\",\"title\",\"description\"],\"arrows\":[]},\"publication\":{\"name\":\"publication\",\"id\":[\"slug\",\"slug\"],\"attributes\":[\"slug\",\"title\",\"content\",\"createdAt\",\"updatedAt\",\"accountLogin\"],\"arrows\":[]}}}"))
			})
		})
	})
})
