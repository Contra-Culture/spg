package data_test

import (
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
								cfg.String("login")
								cfg.String("firstName")
								cfg.String("lastName")
								cfg.String("bio")
								cfg.String("publicationsCount")
								cfg.Link("publications", []string{"publication"}, func(cfg *LinkCfgr) {
									cfg.Map("login", "authors")
									cfg.CounterCache("publicationsCount")
								})
							})
							cfg.Node("rubric", func(cfg *NodeCfgr) {
								cfg.PK([]string{"slug"})
								cfg.String("slug")
								cfg.String("title")
								cfg.String("description")
								cfg.String("publicationsCount")
								cfg.Link("publications", []string{"publication"}, func(cfg *LinkCfgr) {
									cfg.Map("slug", "rubricSlug")
									cfg.CounterCache("publicationsCount")
								})
							})
							cfg.Node("publication", func(cfg *NodeCfgr) {
								cfg.PK([]string{"slug"})
								cfg.String("rubricSlug")
								cfg.String("slug")
								cfg.String("title")
								cfg.String("content")
								cfg.String("createdAt")
								cfg.String("updatedAt")
								cfg.String("accountLogin")
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
						Expect(report.ToString(r)).To(Equal(""))
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
								cfg.String("login")
								cfg.String("firstName")
								cfg.String("lastName")
								cfg.String("bio")
								cfg.String("publicationsCount")
								cfg.Link("publications", []string{"publication"}, func(cfg *LinkCfgr) {
									cfg.Map("login", "authors")
									cfg.CounterCache("publicationsCount")
								})
							})
							cfg.Node("rubric", func(cfg *NodeCfgr) {
								cfg.PK([]string{"slug"})
								cfg.String("slug")
								cfg.String("title")
								cfg.String("description")
								cfg.String("publicationsCount")
								cfg.Link("publications", []string{"publication"}, func(cfg *LinkCfgr) {
									cfg.Map("slug", "rubricSlug")
									cfg.CounterCache("publicationsCount")
								})
							})
							cfg.Node("publication", func(cfg *NodeCfgr) {
								cfg.PK([]string{"slug"})
								cfg.String("rubricSlug")
								cfg.String("slug")
								cfg.String("title")
								cfg.String("content")
								cfg.String("createdAt")
								cfg.String("updatedAt")
								cfg.String("accountLogin")
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
						cfg.String("login")
						cfg.String("firstName")
						cfg.String("lastName")
						cfg.String("bio")
						cfg.String("publicationsCount")
						cfg.Link("publications", []string{"publication"}, func(cfg *LinkCfgr) {
							cfg.Map("login", "authors")
							cfg.CounterCache("publicationsCount")
						})
					})
					cfg.Node("rubric", func(cfg *NodeCfgr) {
						cfg.PK([]string{"slug"})
						cfg.String("slug")
						cfg.String("title")
						cfg.String("description")
						cfg.String("publicationsCount")
						cfg.Link("publications", []string{"publication"}, func(cfg *LinkCfgr) {
							cfg.Map("slug", "rubricSlug")
							cfg.CounterCache("publicationsCount")
						})
					})
					cfg.Node("publication", func(cfg *NodeCfgr) {
						cfg.PK([]string{"slug"})
						cfg.String("rubricSlug")
						cfg.String("slug")
						cfg.String("title")
						cfg.String("content")
						cfg.String("createdAt")
						cfg.String("updatedAt")
						cfg.String("accountLogin")
						cfg.Link("rubric", []string{"rubric"}, func(cfg *LinkCfgr) {
							cfg.Map("rubricSlug", "slug")
							cfg.Limit(1)
						})
						cfg.Link("authors", []string{"account"}, func(cfg *LinkCfgr) {
							cfg.Map("authors", "login")
						})
					})
				})
				Expect(g.JSONString()).To(Equal("{\"dataNodes\":{\"account\":{\"pk\":[\"login\",\"login\"],\"nodes\":{\"login\":{\"pk\":[],\"nodes\":[],\"links\":[]},\"firstName\":{\"pk\":[],\"nodes\":[],\"links\":[]},\"lastName\":{\"pk\":[],\"nodes\":[],\"links\":[]},\"bio\":{\"pk\":[],\"nodes\":[],\"links\":[]},\"publicationsCount\":{\"pk\":[],\"nodes\":[],\"links\":[]}},\"links\":[{\"limit\":\"0\",\"name\":\"publications\",\"hostNode\":[\"account\"],\"remoteNode\":[\"publication\"],\"remoteLink\":\"\",\"counterCache\":\"publicationsCount\",\"mapping\":[[\"login\",\"authors\"]]}]},\"rubric\":{\"pk\":[\"slug\",\"slug\"],\"nodes\":{\"slug\":{\"pk\":[],\"nodes\":[],\"links\":[]},\"title\":{\"pk\":[],\"nodes\":[],\"links\":[]},\"description\":{\"pk\":[],\"nodes\":[],\"links\":[]},\"publicationsCount\":{\"pk\":[],\"nodes\":[],\"links\":[]}},\"links\":[{\"limit\":\"0\",\"name\":\"publications\",\"hostNode\":[\"rubric\"],\"remoteNode\":[\"publication\"],\"remoteLink\":\"\",\"counterCache\":\"publicationsCount\",\"mapping\":[[\"slug\",\"rubricSlug\"]]}]},\"publication\":{\"pk\":[\"slug\",\"slug\"],\"nodes\":{\"rubricSlug\":{\"pk\":[],\"nodes\":[],\"links\":[]},\"slug\":{\"pk\":[],\"nodes\":[],\"links\":[]},\"title\":{\"pk\":[],\"nodes\":[],\"links\":[]},\"content\":{\"pk\":[],\"nodes\":[],\"links\":[]},\"createdAt\":{\"pk\":[],\"nodes\":[],\"links\":[]},\"updatedAt\":{\"pk\":[],\"nodes\":[],\"links\":[]},\"accountLogin\":{\"pk\":[],\"nodes\":[],\"links\":[]}},\"links\":[{\"limit\":\"1\",\"name\":\"rubric\",\"hostNode\":[\"publication\"],\"remoteNode\":[\"rubric\"],\"remoteLink\":\"\",\"counterCache\":\"\",\"mapping\":[[\"rubricSlug\",\"slug\"]]},{\"limit\":\"0\",\"name\":\"authors\",\"hostNode\":[\"publication\"],\"remoteNode\":[\"account\"],\"remoteLink\":\"\",\"counterCache\":\"\",\"mapping\":[[\"authors\",\"login\"]]}]}}}"))
			})
		})
	})
})
