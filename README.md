![Astria server header](assets/logo/logo-sbs-bg.svg)

## `This is all pre-release code, continue at your own peril.`

## What is Astria?

Astria is an open-source, community driven, media manager and transcoding server. The main interface is the [olaris-react](https://gitlab.com/olaris/olaris-react) project although in due time we hope to support multiple clients / applications.

Our core values are:

### Community driven development
We want Olaris to be a community project which means we will heavily prioritise features based on our user feedback.

### Focus on one feature at a time
We will work on features until they are perfect (or as close to it as possible). We would rather have a product where three features work really well than a product with 30 unfinished features.

This does not mean we won't work in parallel, it simply means we will not start anything new until we are happy the new feature works to a high standard.

### Our users are not our product
We don't want to collect metadata, we don't want to sell metadata your data is yours and yours alone.

### Singular Focus: Video.
Our sole focus is on video and video alone, anything that does not meet this requirement will not be considered. This means for example we will never add music support due to different approach that would be required throughout the application. 

### Open-source
Everything we build should be open-source. We feel strongly that more can be achieved with free open-source software. That's why we are aiming to be and to remain open-source instead of open-core where certain features are locked behind a paywall.

## How to run olaris

### Local install

#### Unpack to `/opt`

    sudo unzip olaris-linux-amd64-v0.3.0.zip -d /opt/olaris

Replace the name of the zipfile with the name of the file you downloaded.

### Configuration

Olaris can be configured via configuration file, environment variables, or command-line flags. An `olaris.toml.sample` configuration file is included in the `docs/` folder; rename it to `olaris.toml` and place in `$HOME/.config/olaris`. You can also override the configuration directory location with the `OLARIS_CONFIG_DIR` environment variable or the `--config_dir` command-line flag.

The configuration options available are documented in [CONFIGURATION.md](CONFIGURATION.md).

#### Run as daemon using systemd

To run Olaris as a daemon you may use the supplied systemd unit file:

    mkdir -p ~/.config/systemd/user/
    cp /opt/olaris/doc/config-examples/systemd/olaris.service ~/.config/systemd/user/
    systemctl --user daemon-reload
    systemctl --user start olaris.service

To start Olaris automatically:

    # Allow systemd to start in user mode without a login session
    sudo loginctl enable-linger $USER
    systemctl --user enable olaris.service

### Run using Docker

The following command runs Olaris in a Docker container under your own user‘s UID, ensuring that the Olaris config files end up in your home directory with the correct permissions. It exposes Olaris on port 8080 only on your local machine.

The command below mounts `~/Videos` to `/var/media` in the container --- please update this path to match the location of your media files. When you create a library in Olaris, please keep in mind that Olaris is running inside the container and will see your media at `/var/media/`.

    mkdir -p ~/.config/olaris ~/.config/rclone
    docker run \
      -p 127.0.0.1:8080:8080/tcp \
      -v $HOME/media/:/var/media \
      -v $HOME/.config/olaris:/home/olaris/.config/olaris \
      -v $HOME/.config/rclone:/home/olaris/.config/rclone \
      -e OLARIS_UID=$(id -u) -e OLARIS_GID=$(id -g) \
      olaristv/olaris-server

#### Running the latest build in Docker

To run the latest build from our CI (Continous Integration) infrastructure, use the `olaristv/olaris-server:from-ci` image instead. This will download a new build every time the container is started. Please note that this runs a bleeding-edge development version, which may be horribly unstable!

## How to build

See the [hacking](HACKING.md) document for instructions on how to build Olaris yourself
