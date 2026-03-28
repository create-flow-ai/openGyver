# format

Code formatting and prettifying utilities -- format, beautify, or minify HTML, XML, CSS, and SQL.

## Usage

```bash
openGyver format [subcommand] [input] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Wrap output in `{"input","output","format"}` JSON envelope |
| `--file` | `-f` | string | | Read input from a file instead of an argument |
| `--output` | `-o` | string | | Write output to a file instead of stdout |
| `--help` | `-h` | bool | | Show help for the command |

## Input Methods

All subcommands accept input in three ways (in order of precedence):

1. **`--file` / `-f` flag** -- read from a file
2. **Positional argument** -- pass the string directly
3. **Stdin** -- pipe content from another command

## Subcommands

### html

Format and beautify HTML with proper indentation. Parses the input as HTML, then re-emits it with correct nesting and indentation. Handles self-closing (void) tags, attributes, text nodes, comments, and doctype declarations.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` or stdin) | HTML string to format |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--indent` | `-i` | int | `2` | Number of spaces per indentation level |
| `--minify` | `-m` | bool | `false` | Strip whitespace instead of prettifying |

#### Examples

```bash
# Beautify inline HTML
openGyver format html '<div><p>Hello</p></div>'

# Beautify with 4-space indent
openGyver format html --indent 4 '<html><body><h1>Title</h1></body></html>'

# Minify HTML
openGyver format html --minify '<div>  <p> Hello </p>  </div>'

# Format a file, print to stdout
openGyver format html --file index.html

# Format a file and write to a new file
openGyver format html --file index.html --output pretty.html

# Pipe HTML through the formatter
echo '<ul><li>A</li><li>B</li></ul>' | openGyver format html

# JSON envelope output
openGyver format html '<div><span>test</span></div>' --json

# Minify a file in-place (write to same file)
openGyver format html --minify --file page.html -o page.html
```

#### JSON Output Format

```json
{
  "input": "<div><p>Hello</p></div>",
  "output": "<div>\n  <p>\n    Hello\n  </p>\n</div>",
  "format": "html"
}
```

---

### xml

Format and beautify XML with proper indentation. Parses the input with Go's `encoding/xml` and re-emits it with correct nesting and indentation. Handles elements, attributes, text nodes, comments, processing instructions, and CDATA sections.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` or stdin) | XML string to format |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--indent` | `-i` | int | `2` | Number of spaces per indentation level |
| `--minify` | `-m` | bool | `false` | Strip whitespace instead of prettifying |

#### Examples

```bash
# Beautify inline XML
openGyver format xml '<root><item id="1"><name>Foo</name></item></root>'

# Beautify with 4-space indent from a file
openGyver format xml --indent 4 --file data.xml

# Minify XML
openGyver format xml --minify '<root>  <item />  </root>'

# Format and write to a new file
openGyver format xml --file input.xml --output pretty.xml

# Pipe XML through the formatter
echo '<a><b>text</b></a>' | openGyver format xml

# JSON envelope output
openGyver format xml '<config><key>value</key></config>' --json

# Minify a large XML file
openGyver format xml --minify --file large.xml -o large.min.xml

# Format a Maven POM file
openGyver format xml --file pom.xml --indent 4
```

#### JSON Output Format

```json
{
  "input": "<root><item id=\"1\"><name>Foo</name></item></root>",
  "output": "<root>\n  <item id=\"1\">\n    <name>Foo</name>\n  </item>\n</root>",
  "format": "xml"
}
```

---

### css

Format and beautify CSS with proper indentation. Uses a rule-based approach: puts each property on its own line, indents inside braces, and normalises spacing around colons, semicolons, and braces. Handles selectors, media queries, nested at-rules, and comments.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` or stdin) | CSS string to format |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--indent` | `-i` | int | `2` | Number of spaces per indentation level |
| `--minify` | `-m` | bool | `false` | Strip all unnecessary whitespace |

#### Examples

```bash
# Beautify inline CSS
openGyver format css 'body { color: red; margin: 0; }'

# Beautify with 4-space indent
openGyver format css --indent 4 'h1{font-size:2em;color:blue}'

# Minify CSS
openGyver format css --minify 'body { color: red; margin: 0; }'

# Format a CSS file
openGyver format css --file styles.css

# Format and save to a new file
openGyver format css --file styles.css --output pretty.css

# Pipe CSS through the formatter
echo '.box { padding: 10px; }' | openGyver format css

# Minify a CSS file for production
openGyver format css --minify --file styles.css -o styles.min.css

# JSON envelope output
openGyver format css '.btn { color: blue; }' --json
```

#### JSON Output Format

```json
{
  "input": "body { color: red; margin: 0; }",
  "output": "body {\n  color: red;\n  margin: 0;\n}",
  "format": "css"
}
```

---

### sql

Format and beautify SQL with proper indentation and keyword casing. Uppercases SQL keywords, adds newlines before major clauses, and indents clause bodies. Handles SELECT, FROM, WHERE, JOIN, subqueries, and all common SQL keywords.

**Keywords uppercased:**

SELECT, FROM, WHERE, JOIN, LEFT JOIN, RIGHT JOIN, INNER JOIN, OUTER JOIN, CROSS JOIN, FULL JOIN, ON, ORDER BY, GROUP BY, HAVING, INSERT, UPDATE, DELETE, CREATE, ALTER, DROP, AND, OR, IN, NOT, NULL, AS, SET, VALUES, INTO, LIMIT, OFFSET, UNION, DISTINCT, BETWEEN, LIKE, EXISTS, CASE, WHEN, THEN, ELSE, END

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` or stdin) | SQL string to format |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--indent` | `-i` | int | `2` | Number of spaces per indentation level |
| `--minify` | `-m` | bool | `false` | Put everything on one line, stripping extra whitespace |

#### Examples

```bash
# Beautify a SELECT query
openGyver format sql 'select id, name from users where active = 1 order by name'

# Beautify with 4-space indent and a JOIN
openGyver format sql --indent 4 'select * from orders o join users u on o.user_id = u.id'

# Minify a SQL file
openGyver format sql --minify --file query.sql

# Format and save to file with JSON envelope
openGyver format sql --file complex.sql --output pretty.sql --json

# Pipe SQL through the formatter
echo 'insert into users (name, email) values ("Alice", "a@b.com")' | openGyver format sql

# Format an INSERT statement
openGyver format sql 'insert into products (name, price) values ("Widget", 9.99)'

# Format a complex query with subquery
openGyver format sql 'select * from users where id in (select user_id from orders where total > 100)'

# JSON envelope output
openGyver format sql 'select count(*) from users group by status' --json
```

#### JSON Output Format

```json
{
  "input": "select id, name from users where active = 1 order by name",
  "output": "SELECT\n  id, name\nFROM\n  users\nWHERE\n  active = 1\nORDER BY\n  name",
  "format": "sql"
}
```
