# hanazawabot

[![Build Status](https://circleci.com/gh/h3poteto/hanazawabot.svg?style=shield&circle-token=8c81a54ae9fb7455eb8e742c4de3eb818e2c7e6c)]()

## install

## initial setup
At the first, you should prepare database data, and import all youtube data.
```
./bin/seeds
./bin/all_crawl_youtube
```

## setup crontab
```
01 0 * * * /bin/bash -l -c 'cd /home/akira/projects/hanazawabot/bin && ./daily_crawl_youtube >> ../log/daily_crawl_youtube.log 2>> ../log/daily_crawl_youtube.err.log'

11 * * * * /bin/bash -l -c 'cd /home/akira/projects/hanazawabot/bin && ./refollow >> ../log/refollow.log 2>> ../log/refollow.err.log'
```
