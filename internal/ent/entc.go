//go:build ignore
// +build ignore

package main

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"log"
)

func main() {
	ex, err := entgql.NewExtension()
	err = entc.Generate("./schema",
		&gen.Config{
			Features: []gen.Feature{
				gen.FeaturePrivacy,
				gen.FeatureEntQL,
				gen.FeatureVersionedMigration,
				gen.FeatureUpsert,
				gen.FeatureSnapshot,
			},
			Templates: []*gen.Template{},
		},
		entc.Extensions(ex),
		entc.TemplateDir("./template"),
	)
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
