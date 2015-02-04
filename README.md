# Tic-tac-toe in Terminal

[![Build Status](https://drone.io/github.com/wujiang/tic-tac-toe/status.png)](https://drone.io/github.com/wujiang/tic-tac-toe/latest)
[![Build Status](https://travis-ci.org/wujiang/tic-tac-toe.svg)](https://travis-ci.org/wujiang/tic-tac-toe)
[![Coverage Status](https://coveralls.io/repos/wujiang/tic-tac-toe/badge.svg?branch=master)](https://coveralls.io/r/wujiang/tic-tac-toe?branch=master)

[Tic-tac-toe](http://en.wikipedia.org/wiki/Tic-tac-toe) is a 2-player
zero-sum game on a 3 x 3 grid. This implementation provides an
experience of playing the game within a terminal. The API server can
be running anywhere.

## Usage

- Run server: `ttt-server-openbsd-amd64`
- Run client: `ttt-client-openbsd-amd64`

![Demo](./demo.gif)

[Download](https://github.com/wujiang/binaries/tree/master/tic-tac-toe/)
and try it.

## Install

0. [Go setup](http://golang.org/doc/install)
1. Please use [godep](https://github.com/tools/godep) `godep restore ./...`
   to install all dependencies.
2. Run `go install` in `ttt-client` and `ttt-server`
