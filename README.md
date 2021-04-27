# ezorm v2

**WORK IN PROGRESS**

> NOTE: v2 is not fully compatible with the previous [ezorm(v1)](https://github.com/ezbuy/ezorm/tree/v1)

Supported Drivers:
> ⭐️ : NEW features in v2 and **RECOMMEND TO USE**
> 📁 : DEPRECATED and only do some bugfixes

* ⭐️ mysql.v2: based on the [redis-orm](https://github.com/ezbuy/redis-orm)'s mysql driver and some optimizations
* ⭐️ mongo.v2: based on the official [mongo driver](https://github.com/mongodb/mongo-go-driver)
* 📁 mysql: the previous ezorm mysql driver
* 📁 mongo: the previous ezorm mongo driver
* 📁 elasticSearch: the previous elasticSearch driver
* 📁 redis: the previous redis driver

Goals:
- [ ] Deprecated the old go-bindata dependency , uses the new Go embedded template since Go1.16.
- [ ] Fully Go Modules support , and follow Go Module package semantic management.
- [ ] Better Raw SQL Query support.
- [ ] Better Code Coverage.
- [ ] Reforged mongo connection pool.
- [ ] Deprecated no-context functions , better query tracing integration.
- [ ] More useful template functions.

