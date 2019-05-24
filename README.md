#deploy

1. define amb env 'SQ-SECRET' with a random string (must be secret)  
2. define SHA1-PWD with passwords (sperated by ':')
3. define REDISTOGO_URL with redis url

#TODO
* use hash in token for space save?
* must register token in a db for blacklist
* inform queue size
-- config redis heroku style --
* Authorization with ':'