# Backlog CLI

[Backlog](https://backlog.com/) CLI: like [gh](https://cli.github.com/)

This project aims to provide a Backlog CLI that replicates the user experience of the official GitHub CLI.

## Installation

Download the binary either by using the installation script or by directly downloading it from the [Releases](https://github.com/shufo/backlog-cli/releases).

```bash
# if you want to install to /usr/local/bin
$ curl -sSfL https://raw.githubusercontent.com/shufo/backlog-cli/main/install.sh  | sudo sh -s - -b /usr/local/bin

# if you want to install to /usr/bin
$ curl -sSfL https://raw.githubusercontent.com/shufo/backlog-cli/main/install.sh  | sudo sh -s - -b /usr/bin
```

## Authentication

You must first authenticate with the following command (Required to login to your organization)

```bash
$ backlog auth login
```

![2023-03-19-14-46-29-resize](https://user-images.githubusercontent.com/1641039/226156355-46404529-a869-45b6-9d90-fef45c8ab699.gif)

### Troubleshoot

If you had an error like below. Please logout then login to your organization.

<img width="808" alt="image" src="https://user-images.githubusercontent.com/1641039/226149627-fa45605a-3698-40e3-a0c7-f9844221398a.png">

## Usage

### Issue operation

- List Issues

```bash
$ backlog issue list
# filter isseu assigned to me
$ backlog issue list --me
# view issue list on web
$ backlog issue list --web
```

- View Issue

```bash
$ backlog issue view 123
# view issut on web
$ backlog issue view 123 -w
```

![2023-03-19-15-19-01-resize](https://user-images.githubusercontent.com/1641039/226157765-ffdb7490-7674-4031-a92b-5376236d3e4f.gif)

- Create a new Issue

```bash
$ backlog issue create
```

![2023-03-19-14-56-35-resize](https://user-images.githubusercontent.com/1641039/226156978-4658223d-d172-4522-a7b4-9ea04adf8f05.gif)

- Edit Issue

```bash
$ backlog issue edit 123
```

![2023-03-19-15-50-06-resize](https://user-images.githubusercontent.com/1641039/226159273-810e430a-2d0a-40ce-b578-57bb5dc34a8f.gif)

- View Relevant issues

```bash
$ backlog issue status
```

- Comment to issue

```bash
$ backlog issue comment 32
```

### Alias

- Create a shortcut for a `backlog` command

```bash
$ backlog alias set iv 'issue view'
```

- List aliases

```bash
$ backlog alias list
```

- Delete an alias

```bash
$ backlog alias delete iv
```

![2023-03-19-16-03-07-resize](https://user-images.githubusercontent.com/1641039/226159757-6441d5b8-b70f-4371-9ae8-2eabe8db0993.gif)

## Configuring backlog-cli

To configure project wide settings, put `backlog.json` in your repository root, backlog-cli will treat it as settings files.

```json
{
  "backlog_domain": "backlog.com",
  "organization": "your_organization", // <your_organization>.backlog.com
  "project": "MYPROJECT" // your project key
}
```

## TODO

- [ ] Add command like `gh issue comment`

## Contributing

1.  Fork it
2.  Create your feature branch (`git checkout -b my-new-feature`)
3.  Commit your changes (`git commit -am 'Add some feature'`)
4.  Push to the branch (`git push origin my-new-feature`)
5.  Create new Pull Request

## Testing

```bash
$ go test -v ./...
```

## Development

```bash
$ GO111MODULE=off go get github.com/oxequa/realize
$ realize start
$ ./app <option> <args>
```

## LICENSE

MIT
