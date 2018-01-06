#!/usr/bin/env python
# coding: utf-8
import json
import logging

import webapp2


CHANNEL_ID = 'C21RF6B0W'


class MainHandler(webapp2.RequestHandler):
    def get(self):
        self.response.write('Hello Slack Event.')


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
    ('/slack/event1', MainHandler),
    ('/slack/event', SlackEventReceiveHandler),
], debug=True)
