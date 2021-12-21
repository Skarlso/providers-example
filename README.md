# providers-example

Providers example project for my blog post about Providers Pattern.

The blog is [here](https://skarlso.github.io/2021/12/21/providers-pattern/).

Some sample commands:

```
providers add --name bob --image skarlso/providers:echo-v1 --type container
```

```
providers run --name bob --args 'echo this'
5:52PM INF Getting plugin... name=bob
{"status":"Pulling from skarlso/providers","id":"echo-v1"}
{"status":"Digest: sha256:dc09554d11862dd2d3800b6f65352f89b2639f9ec877ef35697d8b959f17c9dd"}
{"status":"Status: Image is up to date for skarlso/providers:echo-v1"}
5:52PM INF Creating container...
5:52PM INF Starting running command... name=bob
5:52PM INF Starting container...
5:52PM INF Successfully finished command. Output:
echo this

5:52PM INF All done.
```

Listing plugins:

```
providers list
+------+-----------+---------------------------+
| NAME |   TYPE    |      IMAGE/LOCATION       |
+------+-----------+---------------------------+
| bob  | container | skarlso/providers:echo-v1 |
+------+-----------+---------------------------+
```

# Restriction

For simplicity, we will use `~/.config/providers` as a plugin folder. Name of the file will correspond with the name in
db.