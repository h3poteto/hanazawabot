# hanazawabot

[![Build Status](https://circleci.com/gh/h3poteto/hanazawabot.svg?style=shield&circle-token=8c81a54ae9fb7455eb8e742c4de3eb818e2c7e6c)]()

## install
You have to build all golang files for your environments.

Example:
```
gom build -o bin/seeds_386 seeds/seeds.go
```

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

49 * * * * /bin/bash -l -c 'cd /home/akira/projects/hanazawabot/bin && ./tweet >> ../log/tweet.log 2>> ../log/tweet.err.log'

*/5 * * * * /bin/bash -l -c 'cd /home/akira/projects/hanazawabot/bin && if [ ! -e ../tmp/pids/userstream.pid ] || ! ps $(cat ../tmp/pids/userstream.pid) ; then ./userstream >> ../log/userstream.log 2>> ../log/userstream.err.log ; fi'
```
