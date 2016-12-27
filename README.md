# ezorm

ezorm is an code-generation based ORM lib for golang, supporting mongodb/sql server/mysql/redis.

data model is defined with YAML file like:

```yaml
Blog:
  db: mongo
  fields:
    - Title: string
    - Hits: int32
    - Slug: string
      flags: [unique]
    - Body: string
    - User: int32
    - CreateDate: datetime
      flags: [sort]
    - IsPublished: bool
      flags: [index]
  indexes: [[User, IsPublished]]
```

Id field will be automatically included for mongo/mysql/sql server.

# Setup

ezorm templates are defined in tpl folder and managed via [go-bindata](http://github.com/jteeuwen/go-bindata/).

After `go get github.com/ezbuy/ezorm`, you should do:

	make init
	make debugTpl

to initialize the dependencies & link tpl for debug usages.

# Usage

	go install github.com/ezbuy/ezorm
	ezorm gen -i blog.yaml -o .

To generate codes, for model like `Blog`, a blog manager will be generated, supporting ActiveRecord like:

```go
p := blog.BlogMgr.NewBlog()
p.Title = "I like ezorm"
p.Slug = "ezorm"
p.Save()

p, err := blog.BlogMgr.FindBySlug("ezorm")
if err != nil {
	t.Error("find fail")
}
fmt.Println("%v", p)
page.PageMgr.RemoveByID(p.Id())

_, err = blog.BlogMgr.FindBySlug("ezorm")
if err == nil {
	t.Error("delete fail")
}
```
use
  ezorm -h
for more help
  ezorm genmsyaml -d="server=...;user id=...;password=...;DATABASE=..." -t=...  -o=...  -p=...
to generate yaml file
  ezorm genmsorm -d="server=...;user id=...;password=...;DATABASE=..." -t=...  -o=...  -p=...
to generate orm go file directly


