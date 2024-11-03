//go:generate ent generate .
package schema

import (
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Sensor struct {
	ent.Schema
}

func (Sensor) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			NotEmpty().
			Unique().
			Comment("sensor id"),
		field.Int64("type_id").
			Comment("sensor type"),
		field.String("desc").
			Optional().
			Default("").
			Comment("description on sensor"),
		field.Int64("rule_id").
			Optional().
			Default(-1).
			Comment("sensor rule iD"),
		field.Float("threshold").
			Optional().
			Default(0).
			Comment("sensor threshold"),
		field.Float("previous_value").
			Optional().
			Default(0).
			Comment("old value of sensor"),
		field.Bool("deleted").
			Default(false).
			Comment("whether is deleted"),
		field.Time("create_time").
			Default(time.Now).
			Immutable().
			Comment("create time"),
		field.Time("last_update").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("last update time"),
	}
}

// Indexes
func (Sensor) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type_id"),
	}
}
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{
			Table:     "sys_sensors",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_unicode_ci",
			Options:   "ENGINE = InnoDB",
		},
		schema.Comment("Sensor information sheet"),
	}
}
