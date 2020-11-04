package schema

import (
	"time"

	"github.com/facebook/ent"
	"github.com/facebook/ent/dialect"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int8("age").SchemaType(map[string]string{
			dialect.MySQL:    "tinyint(3) unsigned", // Override MySQL.
			dialect.Postgres: "int2",                // Override Postgres.
		}).Positive().StructTag(`json:"age" gqlgen:"age"`),
		field.String("name").Unique(),
		field.Time("created_at").SchemaType(map[string]string{
			dialect.MySQL:    "datetime",       // Override MySQL.
			dialect.Postgres: "timestamptz(6)", // Override Postgres.
		}).Default(time.Now).StructTag(`json:"created_at" gqlgen:"created_at"`),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("cars", Car.Type).StorageKey(edge.Column("owner_id")),
		edge.To("pets", Pet.Type).StorageKey(edge.Column("owner_id")),
		edge.To("friends", User.Type),
		// Create an inverse-edge called "groups" of type `Group`
		// and reference it to the "users" edge (in Group schema)
		// explicitly using the `Ref` method.
		edge.From("groups", Group.Type).Ref("users"),
		edge.From("manage", Group.Type).Ref("admin"),
	}
}
