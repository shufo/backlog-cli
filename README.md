# Backlog CLI

[Backlog](https://backlog.com/) CLI: like [gh](https://cli.github.com/)

## Installation

Download binary by installation script

```bash
# if you want to install to /usr/local/bin
$ curl -sSfL https://raw.githubusercontent.com/shufo/backlog-cli/main/install.sh  | sudo sh -s - -b /usr/local/bin

# if you want to install to /usr/bin
$ curl -sSfL https://raw.githubusercontent.com/shufo/backlog-cli/main/install.sh  | sudo sh -s - -b /usr/bin
```

## Authentication

You must first authenticate with the following command (Required to login to your organization)

```bash
$ bk auth login
```

![2023-03-19 12-25-16](https://user-images.githubusercontent.com/1641039/226151964-9c5d4e9b-df6e-443e-981b-19a21160d920.gif)

### Troubleshoot

If you had an error like below. Please logout then login to your organization.

<img width="808" alt="image" src="https://user-images.githubusercontent.com/1641039/226149627-fa45605a-3698-40e3-a0c7-f9844221398a.png">

## Usage

### Issue operation

- List Issues

```bash
$ bk issue list
# filter isseu assigned to me
$ bk issue list --me
# view issue list on web
$ bk issue list --web
```

- View Issue

```bash
$ bk issue view 123
# view issut on web
$ bk issue view 123 -w
```

- Create a new Issue

```bash
$ bk issue create
```

- Edit Issue

```bash
$ bk issue create
$ bk issue edit 123
```

- View Relevant issues

```bash
$ bk issue status
```

### Alias

- Create a shortcut for a `bk` command

```bash
$ bk alias set iv 'issue view'
```

- List aliases

```bash
$ bk alias list
```

- Delete an alias

```bash
$ bk alias delete iv
```
