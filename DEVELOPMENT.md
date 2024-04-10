# Local development

Setting up local development of the GovCMS CLI is fairly straightforward. You should have Go already installed on your machine.

First, clone the repository:
```
git clone https://github.com/govcms-tests/govcms-cli.git
```

Then inside the repository, either using the command line or your IDE, build it:

```
go build -o /usr/local/bin/govcms
```

The govcms executable can now be used anywhere and tested.
