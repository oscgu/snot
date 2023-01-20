# Snot

A simple notes taking app ;)

# WIP

Really needs a cleanup and some bugfixing. Use at your own risk!

## Config

A config file is created on first run which you can find at `~/.snot/config.yml`

The database resides in the same dir

```
editor: [e.g. vim]             # optional field if you are not happy with the default editor
user:
    name: Mustermann
    group: DevOps
server:
    address: ""
    port: ""
    active: false
```

The server part doesn't do anything (for now).

## Installation

`$ sudo make`

## Usage

Writing a note:

`$ snot note [topic] (title)`

![editor](https://user-images.githubusercontent.com/94227101/211210472-5f4b188f-8139-4389-a2b2-f28ca4b89ce3.png)

Browsing through notes:

`$ snot view `

![view](https://user-images.githubusercontent.com/94227101/211210512-507ae398-ca4a-4b56-b988-a301459a89d6.png)

## Uninstall

`$ make clean`

Deleteting the config and db:

`$ rm -rf ~/.snot`