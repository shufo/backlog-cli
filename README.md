# Backlog CLI

[Backlog](https://backlog.com/) CLI like [gh](https://cli.github.com/)

## Installation

Download binary by installation script

```bash
# if you want to install to /usr/local/bin
$ curl -sSfL https://raw.githubusercontent.com/shufo/ecs-fargate-oneshot/master/install.sh  | sudo sh -s - -b /usr/local/bin

# if you want to install to /usr/bin
$ curl -sSfL https://raw.githubusercontent.com/shufo/ecs-fargate-oneshot/master/install.sh  | sudo sh -s - -b /usr/bin
```

## Authentication

```bash
$ ba auth login
```

### Troubleshoot

If you had an error like `you are not logged in`

## Usage

### Issue operation

- List Issues

```bash
$ ba issue list
# filter isseu assigned to me
$ ba issue list --me
# view issue list on web
$ ba issue list --web
```

- View Issue

```bash
$ ba issue view 123
# view issut on web
$ ba issue view 123 -w
```

- Create a new Issue

```bash
$ ba issue create
```

- Edit Issue

```bash
$ ba issue create
$ ba issue edit 123
```

- View Relevant issues

```bash
$ ba issue status
```

### Alias

- Create a shortcut for a `ba` command

```bash
$ ba alias set iv 'issue view'
```

- List aliases

```bash
$ ba alias list
```

- Delete an alias

```bash
$ ba alias delete iv
```
