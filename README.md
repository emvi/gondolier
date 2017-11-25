<p align="center">
    <img src="gondoliergopher.svg" width="300px" />
</p>

# Gondolier

[![GoDoc](https://godoc.org/github.com/emvicom/gondolier?status.svg)](https://godoc.org/github.com/emvicom/gondolier)
[![CircleCI](https://circleci.com/gh/emvicom/gondolier.svg?style=svg)](https://circleci.com/gh/emvicom/gondolier)

## Description

Gondolier is a library to auto migrate database schemas in Go (golang) using structs. Quick demo:

```
type Customer struct {
    Id   uint64 `gondolier:"type:bigint;id"`
    Name string `gondolier:"type:varchar(255);notnull"`
    Age  int    `gondolier:"type:integer;notnull"`
}

type Order struct {
    Id    uint64 `gondolier:"type:bigint;id"`
    Buyer uint64 `gondolier:"type:bigint;fk:customer.id;notnull"`
}

type OrderPosition struct {
    Id       uint64 `gondolier:"type:bigint;id"`
    Order    uint64 `gondolier:"type:bigint;fk:order.id;notnull"`
    Quantity int    `gondolier:"type:integer;notnull"`
    Cost     int    `gondolier:"type:integer;notnull"`
}

type Obsolete struct{}

func main() {
    // connect to database
    db, _ := sql.Open("postgres", dbString())
    defer db.Close()

    // migrate your schema
    gondolier.Use(db, &gondolier.Postgres{Schema: "public",
        DropColumns: true,
        Log:         true})
    gondolier.Model(Customer{}, Order{}, OrderPosition{})
    gondolier.Drop(Obsolete{})
    gondolier.Migrate()
}
```

[View the full demo](https://github.com/emvicom/gondolier-example)

### Features

* create the initial schema just from your data model defined in Go
* update the schema just from your data model defined in Go
* drop columns when they're no longer needed (removed in struct)
* drop tables by passing a struct (which can be empty)

#### Supported databases

* Postgres

### Limits

* no multi primary key support yet
* there is no way Gondolier can check the data model is valid, so it might fail to execute the migration (with a panic)

## Installation

To install Gondolier, go get all dependencies and Gondolier:

```
go get github.com/lib/pq # for Postgres
go get github.com/emvicom/gondolier
```

## Usage

Gondolier consists just out of a few methods. First, you setup Gondolier by passing the database connection and the migrator to *Use*:

```
gondolier.Use(dbconn, migrator)
```

The migrator is an interface which is used by Gondolier to migrate the data model. Example:

```
gondolier.Postgres{Schema: "public", DropColums: true, Log: true}
```

This will configure the Postgres migrator to use the schema "public", drop columns when the field is missing in the data model and output executed SQL statements to log (using the standard log library).

Now you can define a naming schema used to name tables and columns:

```
gondolier.Naming(&gondolier.SnakeCase{})
```

You can define your own naming schema by implementing the NameSchema interface. Currently SnakeCase is the default. You don't need to call *Naming* to set it.

Now call *Model* and pass the models which define your database schema:

```
gondolier.Model(MyModel{}, &AnotherModel{})
```

*Model* accepts objects and pointers. The models must define a decorator with meta information. For defails take a look at the Postgres Migrator or the example implementation. Here is a short example:

```
type MyModel struct {
    Id       uint64 `gondolier:"type:bigint;id"` // id is a shortcut
    SomeAttr string `gondolier:"type:text;notnull"`
}

type AnotherModel struct {
    Id         uint64   `gondolier:"type:bigint;pk;seq:1,1,-,-,1;default:nextval(seq);notnull"` // long version for "id"
    UniqueAttr int      `gondolier:"type:integer;notnull;unique;default:42"`
    AnArray    []string `gondolier:"type:varchar(100)[]"`
    ForeignKey uint64   `gondolier:"type:bigint;fk:MyModel.Id;notnull"`
}
```

Afterwards, call *Migrate* to start the migration:

```
gondolier.Migrate()
```

To drop a table that is no longer needed, call *Drop*. You can remove all attributes from the struct, just the name must match the old struct:

```
type DropMe struct {}
gondolier.Drop(DropMe{})
```

## Contribute

[See CONTRIBUTING.md](CONTRIBUTING.md)

## License

MIT
