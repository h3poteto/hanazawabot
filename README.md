# hanazawabot

[![CircleCI](https://circleci.com/gh/h3poteto/hanazawabot.svg?style=svg)](https://circleci.com/gh/h3poteto/hanazawabot)

Twitter bot of Kana Hanazawa.
https://twitter.com/hanazawa_sick


## Install
Download from [Releases](https://github.com/h3poteto/hanazawabot/releases).

Example:

```
$ curl -L https://github.com/h3poteto/hanazawabot/releases/download/v0.1.0/hanazawabot_0.1.0_linux_amd64.zip > hanazawabot.zip
$ unzip hanazawabot.zip
```

## Setup
At the first, create database.

```
$ mysql -u root
mysql> CREATE DATABASE hanazawabot CHARACTER SET utf8mb4;
```

Next, prepare database table, and import all youtube data.

```
$ ./hanazawabot migrate
$ ./hanazawabot seed
$ ./hanazawabot allcrawl
```

## Run
Example:

```
$ ./hanazawabot userstream -p pids/userstream.pid
```
