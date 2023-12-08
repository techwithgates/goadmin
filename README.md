# Easy Admin Panel
## Introduction
Easy Admin Panel (EAP) is an open source project that allows Go developers (Gophers) to easily perform CRUD operations on database through a web browser. The purpose of this application is to facilitate development process and it's not made for production use.
<br><br>

## How does it work?
On "tables" page, EAP will display the available table names. When the user selects a table, it will redirect the user to "objects" page where EAP will display the associated table data. When the user clicks the add button or any of the existing data record, EAP will generate a form with certain input fields base on the user's table field definitions.
<br><br>

## Limitations
- compatible with only PostgreSQL
- EAP will only retrieve the tables whose table schema is `public` (default table schema is public)
- tables with the absent of primary key field is not supported
- tables with custom primary key constraint names are not supported (A default pk constraint name looks like `<table_name>_pkey`)
- composite primary key fields may not work properly in update operation
- pagination and search bar features are not currently available
<br><br>

## Explicitly Supported Fields / Data Types
- character and character varying
- integer, smallint and bigint
- numeric
- text
- boolean
- date
- time without time zone
- timestamp without time zone
- json and jsonb
- ARRAY (1 dimensional array)
<br><br>

## Field Specification
- `character` and `character varying` maps to `<input type="text">`

- `integer`, `smallint` and `bigint` maps to `<input type="number">`

- `numeric` maps to `<input type="number" step="any">`

- `text` maps to `<textarea>`

- `boolean` maps to `<input type="checkbox">`

- `date` maps to `<input type="date">`

- `time without time zone` maps to `<input type="time">`

- `timestamp without time zone` maps to `<input type="datetime-local">`

- `json` and `jsonb` maps to `<textarea>`

- `ARRAY` maps to `<textarea>`

Note that EAP will generate the form field as `<input type="text">` if the table has a field or a data type that is not present in the above.

EAP utilizes the comment feature of PostgreSQL to determine what HTML form field should be rendered on the UI. For example, if your table has an email field, you must specify the comment as "email" for that field.

- The comment `email` maps to `<input type="email">`

- The comment `password` maps to `<input type="password">`

- The comment `url` maps to `<input type="url">`

- The comment `file` maps to `<input type="file">`

Note that the comment feature will only work on `character varying` field and if `character varying` field does not have any comment or has an unsupported comment, EAP will select `<input type="text">` by default. Comments are not case sensitive as EAP will handle the conversion implicitly.

EAP also checks for nullable constraint for the table fields and will construct the appropriate HTML input field. For example, if the table field is not nullable, EAP will apply the `required` keyword into the HTML input field. 

EAP will include the primary key fields as part of the HTML input fields on the UI unless the primary key type is `SERIAL`.
<br><br>

# Getting Started
Run the following command to install EAP into your GOPATH.

```
go get github.com/techwithgates/goadmin
```

Next, in your project root directory, create a folder called "media" 

```
mkdir media
```

This is to keep the uploaded files in the media directory and if you don't create the folder with the following name, EAP will create one for you.

After that, specify your database credentials and the EAP server port number to run.

```
goadmin start -d=postgres://<dbuser>:<dbpasswd>@localhost:<dbport>/<dbname> -p=<port-number>
```

# License

This software is released under [MIT LICENSE](./LICENSE)