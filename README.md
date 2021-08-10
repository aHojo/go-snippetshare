# Setup the DB
`docker run -v "$PWD/data":/var/lib/mysql -p 3306:3306 --name snippetsql -e MYSQL_ROOT_PASSWORD=password -d mysql:latest`

```sql

CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

use snippetbox

CREATE TABLE snippets (
  id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(100) NOT NULL,
  content TEXT NOT NULL,
  created DATETIME NOT NULL,
  expires DATETIME NOT NULL
);

CREATE INDEX idx_snippets_create ON snippets(created);

INSERT INTO snippets (title, content, created, expires) VALUES ('HELLO WORLD', 'package main
import "fmt"
func main() {
    fmt.Println("hello world")
}', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY));
INSERT INTO snippets (title, content, created, expires) VALUES ('HELLO WORLD', 'package main
import "fmt"
func main() {
    fmt.Println("hello world")
}', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY));
INSERT INTO snippets (title, content, created, expires) VALUES ('Functions', 'package main

import "fmt"

func plus(a int, b int) int {

    return a + b
}

func plusPlus(a, b, c int) int {
    return a + b + c
}

func main() {

    res := plus(1, 2)
    fmt.Println("1+2 =", res)

    res = plusPlus(1, 2, 3)
    fmt.Println("1+2+3 =", res)
}', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY));
INSERT INTO snippets (title, content, created, expires) VALUES ('Maps', 'package main

import "fmt"

func main() {

    m := make(map[string]int)

    m["k1"] = 7
    m["k2"] = 13

    fmt.Println("map:", m)

    v1 := m["k1"]
    fmt.Println("v1: ", v1)

    fmt.Println("len:", len(m))

    delete(m, "k2")
    fmt.Println("map:", m)

    _, prs := m["k2"]
    fmt.Println("prs:", prs)

    n := map[string]int{"foo": 1, "bar": 2}
    fmt.Println("map:", n)
}', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY));

-- Create a new user
CREATE USER 'appadmin'@'localhost';
GRANT SELECT, INSERT ON snippetbox.* TO 'appadmin'@'localhost';
ALTER USER 'appadmin'@'localhost' IDENTIFIED BY '<insert password here';

-- test the user
mysql -D snippetbox -u appadmin -p

```

## Get the mysql golang driver
`go get github.com/go-sql-driver/mysql@v1`


# Test Inserting a snippet
Only faked data is in createSnippet() right now. 
`curl -iL -X POST http://localhost:4000/snippet/create`