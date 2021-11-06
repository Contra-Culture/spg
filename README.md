# spg - static pages generator

Think about **spg** like about Ruby on Rails, but for static sites and written in Golang.

![model](doc/model.png "model")

## Purpose
Generate static sites using templates.
## Glossary
- **Host** is a top level entity and an entry point.
- **PageGenerator** is an object responsible for page generation.
- **Schema** - is a specification of tuples and associations (arrows) that _spg_ uses to populate templates with exact data by PageGenerators.
- **Arrow** - is a mapping rule between tuples. Arrow can be of type: "has one", "has one through", "has many", "has many through" or "belongs to". Arrows allows to filter tuples of other data related to the particular one tuple.

## DSL hierarchy

- Host
- - Repo
- - - Schema
- - - - Attribute
- - - - FullView
- - - - ListItemView
- - - - CardView
- - - - LinkView
- - - - PageGenerator
- - - - Arrow
- - - - - CollectionView
- - - - - ItemView
- - [ / ] RootPageGenerator
- - - Schema
- - - - View (redefines Host.Repo.Schema.FullView)
- - - - Arrow(<Arrow>) ItemView (redefines Host.Repo.Schema.Arrow.ItemView)
- - - - Arrow(<Arrow>) CollectionView (redefines Host.Repo.Schema.Arrow.CollectionView)
- - - Layout
- - - Template
- - - [ / .?. ] PageGenerator
- - - - Schema
- - - - Layout (redefines higher level PageGenerator's layout)
- - - - Template (redefines higher level PageGenerator's template)
- - - - [ / ... / .?. ] PageGenerator
## Model


## Usage

## Conventions

## Best Practices

## Limitations
