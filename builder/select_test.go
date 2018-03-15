package builder_test

import (
	"fmt"
	"testing"

	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/builder"
	"github.com/ulule/loukoum/stmt"
)

func TestSelect(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name:      "Simple",
			Builder:   loukoum.Select("test"),
			SameQuery: "SELECT test",
		},
		{
			Name:      "Distinct",
			Builder:   loukoum.Select("test").Distinct(),
			SameQuery: "SELECT DISTINCT test",
		},
		{
			Name:      "As",
			Builder:   loukoum.Select(loukoum.Column("test").As("foobar")),
			SameQuery: "SELECT test AS foobar",
		},
		{
			Name:      "Two columns",
			Builder:   loukoum.Select("test", "foobar"),
			SameQuery: "SELECT test, foobar",
		},
		{
			Name:      "Two columns as",
			Builder:   loukoum.Select("test", loukoum.Column("test2").As("foobar")),
			SameQuery: "SELECT test, test2 AS foobar",
		},
		{
			Name: "Three columns as",
			Builders: []builder.Builder{
				loukoum.Select("a", "b", loukoum.Column("c").As("x")),
				loukoum.Select("a", loukoum.Column("b"), loukoum.Column("c").As("x")),
			},
			SameQuery: "SELECT a, b, c AS x",
		},
		{
			Name: "Three columns",
			Builders: []builder.Builder{
				loukoum.Select([]string{"a", "b", "c"}),
				loukoum.Select([]stmt.Column{
					loukoum.Column("a"),
					loukoum.Column("b"),
					loukoum.Column("c"),
				}),
			},
			SameQuery: "SELECT a, b, c",
		},
	})
}

func TestSelect_From(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name:      "Simple",
			Builder:   loukoum.Select("a", "b", "c").From("foobar"),
			SameQuery: "SELECT a, b, c FROM foobar",
		},
		{
			Name:      "As",
			Builder:   loukoum.Select("a").From(loukoum.Table("foobar").As("example")),
			SameQuery: "SELECT a FROM foobar AS example",
		},
	})
}

func TestSelect_Join(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Inner",
			Builders: []builder.Builder{
				loukoum.
					Select("a", "b", "c").
					From("test1").
					Join("test2 ON test1.id = test2.fk_id"),
				loukoum.
					Select("a", "b", "c").
					From("test1").
					Join("test2", "test1.id = test2.fk_id"),
				loukoum.
					Select("a", "b", "c").
					From("test1").
					Join("test2", "test1.id = test2.fk_id", loukoum.InnerJoin),
				loukoum.
					Select("a", "b", "c").
					From("test1").
					Join("test2", "ON test1.id = test2.fk_id"),
				loukoum.
					Select("a", "b", "c").
					From("test1").
					Join("test2", "ON test1.id = test2.fk_id", loukoum.InnerJoin),
			},
			SameQuery: "SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id",
		},
		{
			Name: "Left",
			Builder: loukoum.
				Select("a", "b", "c").
				From("test1").
				Join("test3", "test3.fkey = test1.id", loukoum.LeftJoin),
			SameQuery: "SELECT a, b, c FROM test1 LEFT JOIN test3 ON test3.fkey = test1.id",
		},
		{
			Name: "Right",
			Builder: loukoum.
				Select("a", "b", "c").
				From("test2").
				Join("test4", "test4.gid = test2.id", loukoum.RightJoin),
			SameQuery: "SELECT a, b, c FROM test2 RIGHT JOIN test4 ON test4.gid = test2.id",
		},
		{
			Name: "Two",
			Builders: []builder.Builder{
				loukoum.
					Select("a", "b", "c").
					From("test2").
					Join("test4", "test4.gid = test2.id").
					Join("test3", "test4.uid = test3.id"),
				loukoum.
					Select("a", "b", "c").
					From("test2").
					Join("test4", loukoum.On("test4.gid", "test2.id")).
					Join("test3", loukoum.On("test4.uid", "test3.id")),
				loukoum.
					Select("a", "b", "c").
					From("test2").
					Join(loukoum.Table("test4"), loukoum.On("test4.gid", "test2.id")).
					Join(loukoum.Table("test3"), loukoum.On("test4.uid", "test3.id")),
			},
			SameQuery: fmt.Sprint(
				"SELECT a, b, c FROM test2 INNER JOIN test4 ON test4.gid = test2.id ",
				"INNER JOIN test3 ON test4.uid = test3.id",
			),
		},
	})
}

func TestSelect_WhereOperatorOrder(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Simple",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("id").Equal(1)),
			String:     `SELECT id FROM table WHERE (id = 1)`,
			Query:      `SELECT id FROM table WHERE (id = $1)`,
			NamedQuery: `SELECT id FROM table WHERE (id = :arg_1)`,
			Args:       []interface{}{1},
		},
		{
			Name: "And",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("id").Equal(1)).
				And(loukoum.Condition("slug").Equal("foo")),
			String:     "SELECT id FROM table WHERE ((id = 1) AND (slug = 'foo'))",
			Query:      "SELECT id FROM table WHERE ((id = $1) AND (slug = $2))",
			NamedQuery: "SELECT id FROM table WHERE ((id = :arg_1) AND (slug = :arg_2))",
			Args:       []interface{}{1, "foo"},
		},
		{
			Name: "And with three expressions",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("id").Equal(1)).
				And(loukoum.Condition("slug").Equal("foo")).
				And(loukoum.Condition("title").Equal("hello")),
			String:     "SELECT id FROM table WHERE (((id = 1) AND (slug = 'foo')) AND (title = 'hello'))",
			Query:      "SELECT id FROM table WHERE (((id = $1) AND (slug = $2)) AND (title = $3))",
			NamedQuery: "SELECT id FROM table WHERE (((id = :arg_1) AND (slug = :arg_2)) AND (title = :arg_3))",
			Args:       []interface{}{1, "foo", "hello"},
		},
		{
			Name: "Or",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("id").Equal(1)).
				Or(loukoum.Condition("slug").Equal("foo")).
				Or(loukoum.Condition("title").Equal("hello")),
			String:     "SELECT id FROM table WHERE (((id = 1) OR (slug = 'foo')) OR (title = 'hello'))",
			Query:      "SELECT id FROM table WHERE (((id = $1) OR (slug = $2)) OR (title = $3))",
			NamedQuery: "SELECT id FROM table WHERE (((id = :arg_1) OR (slug = :arg_2)) OR (title = :arg_3))",
			Args:       []interface{}{1, "foo", "hello"},
		},
		{
			Name: "Or with And",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("id").Equal(1)).
				And(loukoum.Condition("slug").Equal("foo")).
				Or(loukoum.Condition("title").Equal("hello")),
			String:     `SELECT id FROM table WHERE (((id = 1) AND (slug = 'foo')) OR (title = 'hello'))`,
			Query:      `SELECT id FROM table WHERE (((id = $1) AND (slug = $2)) OR (title = $3))`,
			NamedQuery: `SELECT id FROM table WHERE (((id = :arg_1) AND (slug = :arg_2)) OR (title = :arg_3))`,
			Args:       []interface{}{1, "foo", "hello"},
		},
		{
			Name: "Or with And nested with 4 subexpressions",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(
					loukoum.Or(loukoum.Condition("id").Equal(1), loukoum.Condition("title").Equal("hello")),
				).
				Or(
					loukoum.And(loukoum.Condition("slug").Equal("foo"), loukoum.Condition("active").Equal(true)),
				),
			String: fmt.Sprint(
				"SELECT id FROM table WHERE (((id = 1) OR (title = 'hello')) OR ",
				"((slug = 'foo') AND (active = true)))",
			),
			Query: fmt.Sprint(
				"SELECT id FROM table WHERE (((id = $1) OR (title = $2)) OR ",
				"((slug = $3) AND (active = $4)))",
			),
			NamedQuery: fmt.Sprint(
				"SELECT id FROM table WHERE (((id = :arg_1) OR (title = :arg_2)) OR ",
				"((slug = :arg_3) AND (active = :arg_4)))",
			),
			Args: []interface{}{1, "hello", "foo", true},
		},
		{
			Name: "Simple And",
			Builders: []builder.Builder{
				loukoum.
					Select("id").
					From("table").
					Where(loukoum.And(
						loukoum.Condition("id").Equal(1),
						loukoum.Condition("title").Equal("hello"),
					)),
				loukoum.
					Select("id").
					From("table").
					Where(loukoum.Condition("id").Equal(1)).
					Where(loukoum.Condition("title").Equal("hello")),
			},
			String:     "SELECT id FROM table WHERE ((id = 1) AND (title = 'hello'))",
			Query:      "SELECT id FROM table WHERE ((id = $1) AND (title = $2))",
			NamedQuery: "SELECT id FROM table WHERE ((id = :arg_1) AND (title = :arg_2))",
			Args:       []interface{}{1, "hello"},
		},
		{
			Name: "Three wheres",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("id").Equal(1)).
				Where(loukoum.Condition("title").Equal("hello")).
				Where(loukoum.Condition("disable").Equal(false)),
			String:     "SELECT id FROM table WHERE (((id = 1) AND (title = 'hello')) AND (disable = false))",
			Query:      "SELECT id FROM table WHERE (((id = $1) AND (title = $2)) AND (disable = $3))",
			NamedQuery: "SELECT id FROM table WHERE (((id = :arg_1) AND (title = :arg_2)) AND (disable = :arg_3))",
			Args:       []interface{}{1, "hello", false},
		},
		{
			Name: "Where Or",
			Builders: []builder.Builder{
				loukoum.
					Select("id").
					From("table").
					Where(loukoum.Condition("id").Equal(1)).
					Or(
						loukoum.Condition("slug").Equal("foo").And(loukoum.Condition("active").Equal(true)),
					),
			},
			String:     "SELECT id FROM table WHERE ((id = 1) OR ((slug = 'foo') AND (active = true)))",
			Query:      "SELECT id FROM table WHERE ((id = $1) OR ((slug = $2) AND (active = $3)))",
			NamedQuery: "SELECT id FROM table WHERE ((id = :arg_1) OR ((slug = :arg_2) AND (active = :arg_3)))",
			Args:       []interface{}{1, "foo", true},
		},
		{
			Name: "Or in Where And",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("id").Equal(1).And(loukoum.Condition("slug").Equal("foo"))).
				Or(loukoum.Condition("active").Equal(true)),
			String:     "SELECT id FROM table WHERE (((id = 1) AND (slug = 'foo')) OR (active = true))",
			Query:      "SELECT id FROM table WHERE (((id = $1) AND (slug = $2)) OR (active = $3))",
			NamedQuery: "SELECT id FROM table WHERE (((id = :arg_1) AND (slug = :arg_2)) OR (active = :arg_3))",
			Args:       []interface{}{1, "foo", true},
		},
	})
}

func TestSelect_WhereEqual(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Equal",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("disabled").Equal(false)),
			String:     "SELECT id FROM table WHERE (disabled = false)",
			Query:      "SELECT id FROM table WHERE (disabled = $1)",
			NamedQuery: "SELECT id FROM table WHERE (disabled = :arg_1)",
			Args:       []interface{}{false},
		},
		{
			Name: "Not Equal",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("disabled").NotEqual(false)),
			String:     "SELECT id FROM table WHERE (disabled != false)",
			Query:      "SELECT id FROM table WHERE (disabled != $1)",
			NamedQuery: "SELECT id FROM table WHERE (disabled != :arg_1)",
			Args:       []interface{}{false},
		},
	})
}

func TestSelect_WhereIs(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Null",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("disabled").Is(nil)),
			SameQuery: "SELECT id FROM table WHERE (disabled IS NULL)",
		},
		{
			Name: "Not true",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("active").IsNot(true)),
			String:     "SELECT id FROM table WHERE (active IS NOT true)",
			Query:      "SELECT id FROM table WHERE (active IS NOT $1)",
			NamedQuery: "SELECT id FROM table WHERE (active IS NOT :arg_1)",
			Args:       []interface{}{true},
		},
	})
}

func TestSelect_WhereGreaterThan(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Greater than 4",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("count").GreaterThan(2)),
			String:     "SELECT id FROM table WHERE (count > 2)",
			Query:      "SELECT id FROM table WHERE (count > $1)",
			NamedQuery: "SELECT id FROM table WHERE (count > :arg_1)",
			Args:       []interface{}{2},
		},
		{
			Name: "Greater than or equal to 4",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("count").GreaterThanOrEqual(4)),
			String:     "SELECT id FROM table WHERE (count >= 4)",
			Query:      "SELECT id FROM table WHERE (count >= $1)",
			NamedQuery: "SELECT id FROM table WHERE (count >= :arg_1)",
			Args:       []interface{}{4},
		},
		{
			Name: "Greater than raw now",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("updated_at").GreaterThanOrEqual(loukoum.Raw("NOW()"))),
			SameQuery: "SELECT id FROM table WHERE (updated_at >= NOW())",
		},
	})
}

func TestSelect_WhereLessThan(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Less than 3",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("count").LessThan(3)),
			String:     "SELECT id FROM table WHERE (count < 3)",
			Query:      "SELECT id FROM table WHERE (count < $1)",
			NamedQuery: "SELECT id FROM table WHERE (count < :arg_1)",
			Args:       []interface{}{3},
		},
		{
			Name: "Less than or equal to 6",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("count").LessThanOrEqual(6)),
			String:     "SELECT id FROM table WHERE (count <= 6)",
			Query:      "SELECT id FROM table WHERE (count <= $1)",
			NamedQuery: "SELECT id FROM table WHERE (count <= :arg_1)",
			Args:       []interface{}{6},
		},
	})
}

func TestSelect_WhereLike(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Like",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("title").Like("foo%")),
			String:     "SELECT id FROM table WHERE (title LIKE 'foo%')",
			Query:      "SELECT id FROM table WHERE (title LIKE $1)",
			NamedQuery: "SELECT id FROM table WHERE (title LIKE :arg_1)",
			Args:       []interface{}{"foo%"},
		},
		{
			Name: "Not like",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("title").NotLike("foo%")),
			String:     "SELECT id FROM table WHERE (title NOT LIKE 'foo%')",
			Query:      "SELECT id FROM table WHERE (title NOT LIKE $1)",
			NamedQuery: "SELECT id FROM table WHERE (title NOT LIKE :arg_1)",
			Args:       []interface{}{"foo%"},
		},
		{
			Name: "Ilike",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("title").ILike("foo%")),
			String:     "SELECT id FROM table WHERE (title ILIKE 'foo%')",
			Query:      "SELECT id FROM table WHERE (title ILIKE $1)",
			NamedQuery: "SELECT id FROM table WHERE (title ILIKE :arg_1)",
			Args:       []interface{}{"foo%"},
		},
		{
			Name: "Not ilike",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("title").NotILike("foo%")),
			String:     "SELECT id FROM table WHERE (title NOT ILIKE 'foo%')",
			Query:      "SELECT id FROM table WHERE (title NOT ILIKE $1)",
			NamedQuery: "SELECT id FROM table WHERE (title NOT ILIKE :arg_1)",
			Args:       []interface{}{"foo%"},
		},
	})
}

func TestSelect_WhereBetween(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "10 and 20",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("count").Between(10, 20)),
			String:     "SELECT id FROM table WHERE (count BETWEEN 10 AND 20)",
			Query:      "SELECT id FROM table WHERE (count BETWEEN $1 AND $2)",
			NamedQuery: "SELECT id FROM table WHERE (count BETWEEN :arg_1 AND :arg_2)",
			Args:       []interface{}{10, 20},
		},
		{
			Name: "50 and 70",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("count").NotBetween(50, 70)),
			String:     "SELECT id FROM table WHERE (count NOT BETWEEN 50 AND 70)",
			Query:      "SELECT id FROM table WHERE (count NOT BETWEEN $1 AND $2)",
			NamedQuery: "SELECT id FROM table WHERE (count NOT BETWEEN :arg_1 AND :arg_2)",
			Args:       []interface{}{50, 70},
		},
	})
}

func TestSelect_WhereIn(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "In integers",
			Builders: []builder.Builder{
				loukoum.
					Select("id").
					From("table").
					Where(loukoum.Condition("id").In([]int64{1, 2, 3})),
				loukoum.
					Select("id").
					From("table").
					Where(loukoum.Condition("id").In(int64(1), int64(2), int64(3))),
			},
			String:     "SELECT id FROM table WHERE (id IN (1, 2, 3))",
			Query:      "SELECT id FROM table WHERE (id IN ($1, $2, $3))",
			NamedQuery: "SELECT id FROM table WHERE (id IN (:arg_1, :arg_2, :arg_3))",
			Args:       []interface{}{int64(1), int64(2), int64(3)},
		},
		{
			Name: "Not in integers",
			Builders: []builder.Builder{
				loukoum.
					Select("id").
					From("table").
					Where(loukoum.Condition("id").NotIn([]int{1, 2, 3})),
				loukoum.
					Select("id").
					From("table").
					Where(loukoum.Condition("id").NotIn(1, 2, 3)),
			},
			String:     "SELECT id FROM table WHERE (id NOT IN (1, 2, 3))",
			Query:      "SELECT id FROM table WHERE (id NOT IN ($1, $2, $3))",
			NamedQuery: "SELECT id FROM table WHERE (id NOT IN (:arg_1, :arg_2, :arg_3))",
			Args:       []interface{}{1, 2, 3},
		},
		{
			Name: "In strings",
			Builders: []builder.Builder{
				loukoum.
					Select("id").
					From("table").
					Where(loukoum.Condition("status").In([]string{"read", "unread"})),
				loukoum.
					Select("id").
					From("table").
					Where(loukoum.Condition("status").In("read", "unread")),
			},
			String:     "SELECT id FROM table WHERE (status IN ('read', 'unread'))",
			Query:      "SELECT id FROM table WHERE (status IN ($1, $2))",
			NamedQuery: "SELECT id FROM table WHERE (status IN (:arg_1, :arg_2))",
			Args:       []interface{}{"read", "unread"},
		},
		{
			Name: "Not in strings",
			Builders: []builder.Builder{
				loukoum.
					Select("id").
					From("table").
					Where(loukoum.Condition("status").NotIn([]string{"read", "unread"})),
				loukoum.
					Select("id").
					From("table").
					Where(loukoum.Condition("status").NotIn("read", "unread")),
			},
			String:     "SELECT id FROM table WHERE (status NOT IN ('read', 'unread'))",
			Query:      "SELECT id FROM table WHERE (status NOT IN ($1, $2))",
			NamedQuery: "SELECT id FROM table WHERE (status NOT IN (:arg_1, :arg_2))",
			Args:       []interface{}{"read", "unread"},
		},
		{
			Name: "In single string",
			Builders: []builder.Builder{
				loukoum.
					Select("id").
					From("table").
					Where(loukoum.Condition("status").In("read")),
				loukoum.
					Select("id").
					From("table").
					Where(loukoum.Condition("status").In([]string{"read"})),
			},
			String:     "SELECT id FROM table WHERE (status IN ('read'))",
			Query:      "SELECT id FROM table WHERE (status IN ($1))",
			NamedQuery: "SELECT id FROM table WHERE (status IN (:arg_1))",
			Args:       []interface{}{"read"},
		},
		{
			Name: "Not in single string",
			Builders: []builder.Builder{
				loukoum.
					Select("id").
					From("table").
					Where(loukoum.Condition("status").NotIn("read")),
				loukoum.
					Select("id").
					From("table").
					Where(loukoum.Condition("status").NotIn([]string{"read"})),
			},
			String:     "SELECT id FROM table WHERE (status NOT IN ('read'))",
			Query:      "SELECT id FROM table WHERE (status NOT IN ($1))",
			NamedQuery: "SELECT id FROM table WHERE (status NOT IN (:arg_1))",
			Args:       []interface{}{"read"},
		},
		{
			Name: "In subquery",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("id").In(
					loukoum.Select("id").
						From("table").
						Where(loukoum.Condition("id").Equal(1)),
				)),
			String:     "SELECT id FROM table WHERE (id IN (SELECT id FROM table WHERE (id = 1)))",
			Query:      "SELECT id FROM table WHERE (id IN (SELECT id FROM table WHERE (id = $1)))",
			NamedQuery: "SELECT id FROM table WHERE (id IN (SELECT id FROM table WHERE (id = :arg_1)))",
			Args:       []interface{}{1},
		},
		{
			Name: "Not in subquery",
			Builder: loukoum.
				Select("id").
				From("table").
				Where(loukoum.Condition("id").NotIn(
					loukoum.Select("id").
						From("table").
						Where(loukoum.Condition("id").Equal(1)),
				)),
			String:     "SELECT id FROM table WHERE (id NOT IN (SELECT id FROM table WHERE (id = 1)))",
			Query:      "SELECT id FROM table WHERE (id NOT IN (SELECT id FROM table WHERE (id = $1)))",
			NamedQuery: "SELECT id FROM table WHERE (id NOT IN (SELECT id FROM table WHERE (id = :arg_1)))",
			Args:       []interface{}{1},
		},
	})
}

func TestSelect_GroupBy(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "One column",
			Builders: []builder.Builder{
				loukoum.
					Select("name", "COUNT(*)").
					From("user").
					Where(loukoum.Condition("disabled").IsNull(false)).
					GroupBy("name"),
				loukoum.
					Select("name", "COUNT(*)").
					From("user").
					Where(loukoum.Condition("disabled").IsNull(false)).
					GroupBy(loukoum.Column("name")),
			},
			SameQuery: "SELECT name, COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name",
		},
		{
			Name: "Two columns",
			Builders: []builder.Builder{
				loukoum.
					Select("name", "locale", "COUNT(*)").
					From("user").
					Where(loukoum.Condition("disabled").IsNull(false)).
					GroupBy("name", "locale"),
				loukoum.
					Select("name", "locale", "COUNT(*)").
					From("user").
					Where(loukoum.Condition("disabled").IsNull(false)).
					GroupBy(loukoum.Column("name"), loukoum.Column("locale")),
			},
			SameQuery: fmt.Sprint(
				"SELECT name, locale, COUNT(*) FROM user ",
				"WHERE (disabled IS NOT NULL) GROUP BY name, locale",
			),
		},
		{
			Name: "Three columns",
			Builders: []builder.Builder{
				loukoum.
					Select("name", "locale", "country", "COUNT(*)").
					From("user").
					Where(loukoum.Condition("disabled").IsNull(false)).
					GroupBy("name", "locale", "country"),
				loukoum.
					Select("name", "locale", "country", "COUNT(*)").
					From("user").
					Where(loukoum.Condition("disabled").IsNull(false)).
					GroupBy(loukoum.Column("name"), loukoum.Column("locale"), loukoum.Column("country")),
			},
			SameQuery: fmt.Sprint(
				"SELECT name, locale, country, COUNT(*) FROM user ",
				"WHERE (disabled IS NOT NULL) GROUP BY name, locale, country",
			),
		},
	})
}

func TestSelect_Having(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "One condition",
			Builder: loukoum.
				Select("name", "COUNT(*)").
				From("user").
				Where(loukoum.Condition("disabled").IsNull(false)).
				GroupBy("name").
				Having(loukoum.Condition("COUNT(*)").GreaterThan(10)),
			String: fmt.Sprint(
				"SELECT name, COUNT(*) FROM user ",
				"WHERE (disabled IS NOT NULL) GROUP BY name HAVING (COUNT(*) > 10)",
			),
			Query: fmt.Sprint(
				"SELECT name, COUNT(*) FROM user ",
				"WHERE (disabled IS NOT NULL) GROUP BY name HAVING (COUNT(*) > $1)",
			),
			NamedQuery: fmt.Sprint(
				"SELECT name, COUNT(*) FROM user ",
				"WHERE (disabled IS NOT NULL) GROUP BY name HAVING (COUNT(*) > :arg_1)",
			),
			Args: []interface{}{10},
		},
		{
			Name: "Two conditions",
			Builder: loukoum.
				Select("name", "COUNT(*)").
				From("user").
				Where(loukoum.Condition("disabled").IsNull(false)).
				GroupBy("name").
				Having(
					loukoum.Condition("COUNT(*)").GreaterThan(10).And(loukoum.Condition("COUNT(*)").LessThan(500)),
				),
			String: fmt.Sprint(
				"SELECT name, COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name ",
				"HAVING ((COUNT(*) > 10) AND (COUNT(*) < 500))",
			),
			Query: fmt.Sprint(
				"SELECT name, COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name ",
				"HAVING ((COUNT(*) > $1) AND (COUNT(*) < $2))",
			),
			NamedQuery: fmt.Sprint(
				"SELECT name, COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name ",
				"HAVING ((COUNT(*) > :arg_1) AND (COUNT(*) < :arg_2))",
			),
			Args: []interface{}{10, 500},
		},
	})
}

func TestSelect_OrderBy(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "With order asc",
			Builders: []builder.Builder{
				loukoum.
					Select("name").
					From("user").
					OrderBy(loukoum.Order("id")),
				loukoum.
					Select("name").
					From("user").
					OrderBy(loukoum.Order("id", loukoum.Asc)),
				loukoum.
					Select("name").
					From("user").
					OrderBy(loukoum.Column("id").Asc()),
			},
			SameQuery: "SELECT name FROM user ORDER BY id ASC",
		},
		{
			Name: "With order desc",
			Builders: []builder.Builder{
				loukoum.
					Select("name").
					From("user").
					OrderBy(loukoum.Order("id", loukoum.Desc)),
				loukoum.
					Select("name").
					From("user").
					OrderBy(loukoum.Column("id").Desc()),
			},
			SameQuery: "SELECT name FROM user ORDER BY id DESC",
		},
		{
			Name: "With order desc and asc",
			Builders: []builder.Builder{
				loukoum.
					Select("name").
					From("user").
					OrderBy(loukoum.Order("locale"), loukoum.Order("id", loukoum.Desc)),
				loukoum.
					Select("name").
					From("user").
					OrderBy(loukoum.Order("locale")).
					OrderBy(loukoum.Order("id", loukoum.Desc)),
				loukoum.
					Select("name").
					From("user").
					OrderBy(loukoum.Column("locale").Asc(), loukoum.Column("id").Desc()),
				loukoum.
					Select("name").
					From("user").
					OrderBy(loukoum.Column("locale").Asc()).
					OrderBy(loukoum.Column("id").Desc()),
			},
			SameQuery: "SELECT name FROM user ORDER BY locale ASC, id DESC",
		},
	})
}

func TestSelect_Limit(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "int 10",
			Builder: loukoum.
				Select("name").
				From("user").
				Limit(10),
			SameQuery: "SELECT name FROM user LIMIT 10",
		},
		{
			Name: "string 50",
			Builder: loukoum.
				Select("name").
				From("user").
				Limit("50"),
			SameQuery: "SELECT name FROM user LIMIT 50",
		},
		{
			Name: "uint64 700",
			Builder: loukoum.
				Select("name").
				From("user").
				Limit(uint64(700)),
			SameQuery: "SELECT name FROM user LIMIT 700",
		},
		{
			Name: "Corner case 0",
			Failure: func() builder.Builder {
				return loukoum.Select("name").From("user").Limit(700.2)
			},
		},
		{
			Name: "Corner case 1",
			Failure: func() builder.Builder {
				return loukoum.Select("name").From("user").Limit(float32(700.2))
			},
		},
		{
			Name: "Corner case 2",
			Failure: func() builder.Builder {
				return loukoum.Select("name").From("user").Limit(-700)
			},
		},
	})
}

func TestSelect_Offset(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "int 10",
			Builder: loukoum.
				Select("name").
				From("user").
				Offset(10),
			SameQuery: "SELECT name FROM user OFFSET 10",
		},
		{
			Name: "string 50",
			Builder: loukoum.
				Select("name").
				From("user").
				Offset("50"),
			SameQuery: "SELECT name FROM user OFFSET 50",
		},
		{
			Name: "uint64 700",
			Builder: loukoum.
				Select("name").
				From("user").
				Offset(uint64(700)),
			SameQuery: "SELECT name FROM user OFFSET 700",
		},
		{
			Name: "Corner case 0",
			Failure: func() builder.Builder {
				return loukoum.Select("name").From("user").Offset(700.2)
			},
		},
		{
			Name: "Corner case 1",
			Failure: func() builder.Builder {
				return loukoum.Select("name").From("user").Offset(float32(700.2))
			},
		},
		{
			Name: "Corner case 2",
			Failure: func() builder.Builder {
				return loukoum.Select("name").From("user").Offset(-700)
			},
		},
	})
}

func TestSelect_Extra(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Prefix",
			Builder: loukoum.
				Select("name").
				From("user").
				Prefix("EXPLAIN ANALYZE"),
			SameQuery: "EXPLAIN ANALYZE SELECT name FROM user",
		},
		{
			Name: "Suffix",
			Builder: loukoum.
				Select("name").
				From("user").
				Suffix("FOR UPDATE"),
			SameQuery: "SELECT name FROM user FOR UPDATE",
		},
	})
}
