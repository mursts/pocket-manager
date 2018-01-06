#!/usr/bin/env python
# coding: utf-8

import datetime
import json
import logging
import threading
import time

import webapp2
from google.appengine.api import urlfetch
from pytz import timezone

import config

POCKET_API_ENDPOINT = 'https://getpocket.com/v3'


class SlackPostThread(threading.Thread):

    def __init__(self, id, title, url):
        super(SlackPostThread, self).__init__()
        self.id = id
        self.title = title
        self.url = url

    def run(self):
        payload = {"text": '- {} ({})\n    {}\n\n'.format(self.title, self.id, self.url)}

        r = urlfetch.fetch(config.post_url,
                           payload=json.dumps(payload),
                           method=urlfetch.POST,
                           headers={'Content-Type': 'application/json'})
        logging.debug(r.content)


class MainHandler(webapp2.RequestHandler):
    def get(self):
        self.response.write('Hello world!')


def get_last_pocket_post():
    """
    前日にストックした記事を取得します

    :return: json
    """

    tz = timezone('Asia/Tokyo')

    # 前日のunixtime
    d = datetime.datetime.now(tz=tz) + datetime.timedelta(days=-1)
    d = d.replace(hour=0, minute=0, second=0, microsecond=0)
    unixtime = int(time.mktime(d.utctimetuple()))

    payload = {'consumer_key': config.consumer_key,
               'access_token': config.access_token,
               'state': 'unread',
               'sort': 'newest',
               'since': unixtime}

    r = urlfetch.fetch(POCKET_API_ENDPOINT + '/get',
                       payload=json.dumps(payload),
                       method=urlfetch.POST,
                       headers={'Content-Type': 'application/json'})
    return r.content


class PostSlackHandler(webapp2.RequestHandler):
    def get(self):

        content = get_last_pocket_post()

        logging.debug(json.loads(content))

        json_obj = json.loads(content).get('list', {})
        if len(json_obj) < 1:
            return

        for article in json_obj.values():
            id = article.get('item_id')
            title = article.get('resolved_title') if article.get('resolved_title', '') != '' else article.get('given_title')
            if title is not None:
                title = title.encode('utf-8')
            else:
                title = ''
            url = article.get('resolved_url') if article.get('resolved_url', '') != '' else article.get('given_url')

            t = SlackPostThread(id, title, url)
            t.start()


app = webapp2.WSGIApplication([
    ('/', MainHandler),
    ('/post', PostSlackHandler)
], debug=True)
