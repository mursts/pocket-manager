==============================
Pocket Manager
==============================

`Pocket <getpocket.com>`_ の積読解消のために前日ストックしたものをSlackに通知します

Access token を取得
==============================
 
.. code-block:: bash

   $ curl -X POST -H "Content-Type: application/json; charset=UTF-8" -d '{"consumer_key":"your consumer key","redirect_uri":"pocketapp1234:authorizationFinished"}' https://getpocket.com/v3/oauth/request

   > code=xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx

   $ open https://getpocket.com/auth/authorize?request_token={your code}

   $ curl -X POST -H "Content-Type: application/json; charset=UTF-8" -d '{"consumer_key":"your consuker key","code":"your code"}' https://getpocket.com/v3/oauth/authorize
   
   > access_token=xxxxxxx-xxxx-xxxx-xxxx-xxxxxx&username={your user name}

Todo
==============================

- Slackの `Add Reaction` でアーカイブ

