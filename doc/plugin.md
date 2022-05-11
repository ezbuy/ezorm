## Binary Plugin

Per [#185](https://github.com/ezbuy/ezorm/issues/185), we decide to support more drivers using the plugin way.

The plugin works like `protoc` , developers should impls theirs generator.

ezorm.v2 will parse the YAML metadata and the call the bin , passes the metadata with `JSON` schema to `STDIN`, so ,
the user-side binary should read the metadata from `STDIN` and can use our `pkg/plugin` to decode the metadata to the the generator.

Check the [e2e example](../e2e/plugins/hello-driver/) for the detail usage.
