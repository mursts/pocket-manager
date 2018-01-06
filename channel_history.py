#!/usr/bin/env python
# coding: utf-8

import datetime
import json
import logging
import re
import time
import urllib

import webapp2
from google.appengine.api import urlfetch
from pytz import timezone

import config

rgx = re.compile('.+? \((.*)\)', re.I)


class MainHandler(webapp2.RequestHandler):
    def get(self):
        self.response.write('Hello Slack Event.')

        tz = timezone('Asia/Tokyo')
        d = datetime.datetime.now(tz=tz) + datetime.timedelta(days=-1)
        unixtime = time.mktime(d.utctimetuple())

        logging.error(unixtime)

        query_string = {
            'token': config.token,
            'channel': config.channel_id,
            'oldest': unixtime,
        }
        url = 'https://slack.com/api/channels.history'

        logging.error(urllib.urlencode(query_string))

        r = urlfetch.fetch(url + '?' + urllib.urlencode(query_string),
                           method=urlfetch.GET)

        logging.error(r.content)
        json_obj = json.loads(r.content)
        logging.error(json_obj)
        messages = json_obj["messages"]
        for i in messages:
            match = rgx.search(i["text"])
            if match:
                logging.error(match.group(1))

        # TODO request pocket api for item archive


class SlackEventReceiveHandler(webapp2.RequestHandler):
    def post(self):

        body = self.request.body

        json_obj = json.loads(body)
        logging.error(body)

        event_type = json_obj.get('type', None)

        if event_type == 'url_verification':
            self.response.headers['Content-type'] = 'application/x-www-form-urlencoded'
            self.response.write(json_obj.get('challenge', None))

        elif event_type == 'event_callback':
            logging.error('event callback')

        else:
            self.response.set_status(400, 'Bad Request.')


app = webapp2.WSGIApplication([
    ('/slack/history', MainHandler),
    ('/slack/event', SlackEventReceiveHandler),
], debug=True)
