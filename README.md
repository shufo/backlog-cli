# Backlog CLI

## Installation

Download binary by installation script

```bash
# if you want to install to /usr/local/bin
$ curl -sSfL https://raw.githubusercontent.com/shufo/ecs-fargate-oneshot/master/install.sh  | sudo sh -s - -b /usr/local/bin

# if you want to install to /usr/bin
$ curl -sSfL https://raw.githubusercontent.com/shufo/ecs-fargate-oneshot/master/install.sh  | sudo sh -s - -b /usr/bin
```

## Authentication

```
$ ba auth login
```

## Usage

- List Issues

```
$ ba issue list
# filter isseu assigned to me
$ ba issue list --me
# view issue list on web
$ ba issue list --web
```

- View Issue

```
$ ba issue view 123
# view issut on web
$ ba issue view 123 -w
```

- Create Issue

```
$ ba issue create
```
