# Terraform provider

```
                  ...                        .....   ..         ..  ............    ........  
                 .:l,.                      .;ccc'  .,c.       .,c. '::::clc:::,. .,::::::::;.  
                .:ccl,.                    .:c,,l;. .;l.       .;l' .....:l'.... .:c'.......:c'  
              .'lc.'cc,.                  .cc' .:c. .;l'       .;l'     .;l'     'l;.       'l;  
             .,lc. .;;;;.               .'c:.  .,l, .;l'.      .;l'     .:l.     'l;.       'l;.  
            .,cc'   ,;.;;.             .,l:.    'c:. .cc;'..   .;l'     .:c.     .cc..     .:l'  
           .;::,.   .:'.,:.           .,c;..,;;;:lc.  .;::c:::::lc.     .;c.      .:c:;;;;:c:'.  
          .;,;;.    .;;..,:.          .... ....'....     ..........      ...       ...'''''..  
         ';',;.      ,:. .,c'  
       .,;',c'       .:,  .;l'           ..''''''.  ...            .',,,,,'..   ...       ... ..'''''''..  
      .;c:cl:;,.      ;c.  ;xo,        .;::;;;;;;. .:c.          .;c:;;;;;:c;. .:c.      .:c. .cl:;;;;::c,.  
     .:l::c,..,;,.    'c, .l:;l;.     .:c'.        .:l.         .:l'.     .'l:..:c.      .:c. .c:.     .;l,.  
    .:c'.;;.   .,:;.  .c:.;o' ,o:.    ,c'          .:l'         .cc.       .cc..:c.      .:c. .c:.      .c:.  
    .,:.;:.      .,c;. ,l;lc. 'oo.    'c,.         .:l'         .:c.       .cc..;l'      .:c. .c:.      .c:.  
     .:lc.         .;c:;okd;,cc;.     .;c;........ .:l,..........,l;.......:l,. .:c;......:c. .cc......':c'  
      .::,,;;;;;,,,'';lx0XOoc'.        ..;:::::::'..,c::::::::;. .';:::::::;'.   .';::::::c;. .;c:::::::,.  
         .....'',,;;;::lddc.               ......    ..........     .......          .......   ........  

```
<br />
<br />

AutoCloud Terraform Provider <!-- omit in toc -->
=========================================

Contents <!-- omit in toc -->
--------
- [Terraform provider](#terraform-provider)
  - [Overview](#overview)
  - [Getting Started](#getting-started)
  - [Setup gitlab](#setup-gitlab)
  - [Testing](#testing)
  - [Contributing](#contributing)



Overview
--------

This code is intended to represent the developer ergonimics for Terraform users deploying the AutoCloud Terraform provider in their IaC codebases. It can serve as a design guide and a reference for test cases. This is a work in progress, and should not be considered definitive until noted in this README and matched with a 1.x tag.



Getting Started
---------------

This code assumes the following dependencies are installed locally:

| Package       | Version |
|---------------|---------|
| go            | 1.19.x  |
| pre-commmit   | 2.20.x  |
| commitizen    | 2.32.x  |
| fd            | 8.3.x   |
| golangci-lint | 1.50.x  |
| gosec         | -       |
| shellcheck    | 0.8.x   |

On MacOS, these dependencies may be installed with Homebrew:

```bash
brew install \
  commitizen \
  fd \
  pre-commit \
  shellcheck \
  go \
  golangci-lint
```

`gosec` must be installed via go modules:

```bash
go install github.com/securego/gosec/v2/cmd/gosec@latest
```

To get started with this code, clone this repository to your local machine and run:

 ```bash
 pre-commit install --hook-type pre-commit --hook-type commit-msg --install-hooks
 ```



Setup gitlab
-----------

You need to set your machine to communicate with private gitlab repositories

1. setup a GOPRIVATE env 
```bash
go env -w GOPRIVATE=gitlab.com/auto-cloud
```

2. Then you need to force git to use ssh instead of gitlab (assuming you already are in autoclouds gitlab org)
```bash
git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/"   
```

3. Finally you have to create a [personal token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html) in gitlab with the following scopes `read_api` and `read_repository` and paste it in you ~/.netrc file
```
machine gitlab.com
login <your user name>
password <your token>
```

4. you should be able to run `go mod tidy` successfuly 


Testing
-------

Before testing, the go dependencies need to be installed:

```bash
go mod tidy && go mod vendor
```

Setup env variables and login to AWS:

```bash
aws sso login --profile autocloud-aws-sso-sandbox-developer
export AWS_PROFILE=autocloud-aws-sso-sandbox-developer`

export $(grep -v '^#' .env | xargs)`
```

Finally run `go run main.go` to test the SDK.

Feel free to modify main.go to test the available commands.





Contributing
------------

This codebase expects and enforces conventional commit messages. See the [documentation](https://www.conventionalcommits.org/en/v1.0.0/) for examples.
