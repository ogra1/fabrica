# Fabrica - snap build factory

Fabrica is a web service to be run on an lxd appliance. It spawns a
web application that allows you to point to cloneable git trees, initializes
lxd and builds snap packages of the provided source trees.

Once the snap is installed, it needs its interfaces to be connected:
```bash
snap install lxd
sudo lxd init # hit enter for all questions

sudo snap install fabrica

sudo snap connect fabrica:lxd lxd:lxd
sudo snap connect fabrica:mount-observe :mount-observe
sudo snap connect fabrica:system-observe :system-observe
```

## Options
Fabrica has a configuration option that may be useful for a development environment:
- `debug=true`: (default false) if a build fails, the LXD container is retained for debugging

Options can be set using a `snap set` command e.g.
```
sudo snap set fabrica debug=true
```

If the `debug` option is used, the container will need to be deleted manually to
recover disk space.

## Development Environment
The build needs Go 13.* and npm installed.

### Building the web pages
The web pages use [create-react-app](https://github.com/facebook/create-react-app)
which needs an up-to-date version of Node.
 ```
cd webapp
npm install
./build.sh
```

### Building the application
The application is packaged as a [snap](https://snapcraft.io/docs) and can be
built using the `snapcraft` command. The [snapcraft.yaml](snap/snapcraft.yaml)
is the source for building the application and the name of the snap needs to be
modified in that file.

For testing purposes, it can also be run via:
```
go run fabrica.go
```

## Credits
Icons provided by [Font Awesome](https://fontawesome.com/)
