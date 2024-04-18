# Pocket Manager

[Pocket](getpocket.com) の積読解消のために前日ストックしたものをSlackに通知します

## Deplog

```sh
export POCKET_CONSUMER_KEY=xxxxxx
export POCKET_ACCESS_TOKEN=xxxxxx
export SLACK_POST_URL=xxxxxx

gcloud builds submit \
  --config cloudbuild.yaml \
  --substitutions=_POCKET_CONSUMER_KEY="${POCKET_CONSUMER_KEY}",_POCKET_ACCESS_TOKEN="${POCKET_ACCESS_TOKEN}",_SLACK_POST_URL="${SLACK_POST_URL}"
  .
```

## Access token を取得

```shell script

   $ curl -X POST -H "Content-Type: application/json; charset=UTF-8" \
       -d '{"consumer_key":"your consumer key","redirect_uri":"pocketapp1234:authorizationFinished"}' \
       https://getpocket.com/v3/oauth/request

   > code=xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx

   $ open https://getpocket.com/auth/authorize?request_token={your code}

   $ curl -X POST -H "Content-Type: application/json; charset=UTF-8" \
       -d '{"consumer_key":"your consuker key","code":"your code"}' \
       https://getpocket.com/v3/oauth/authorize

   > access_token=xxxxxxx-xxxx-xxxx-xxxx-xxxxxx&username={your user name}
```

## Todo

- Slackの `Add Reaction` でアーカイブ
