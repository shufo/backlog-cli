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
$ bk auth login
```

![2023-03-19-14-46-29-resize](https://user-images.githubusercontent.com/1641039/226156355-46404529-a869-45b6-9d90-fef45c8ab699.gif)

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

![2023-03-19-15-19-01-resize](https://user-images.githubusercontent.com/1641039/226157765-ffdb7490-7674-4031-a92b-5376236d3e4f.gif)

- Create a new Issue

```bash
$ bk issue create
```

![2023-03-19-14-56-35-resize](https://user-images.githubusercontent.com/1641039/226156978-4658223d-d172-4522-a7b4-9ea04adf8f05.gif)

- Edit Issue

```bash
$ bk issue edit 123
```

![2023-03-19-15-50-06-resize](https://user-images.githubusercontent.com/1641039/226159273-810e430a-2d0a-40ce-b578-57bb5dc34a8f.gif)

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
