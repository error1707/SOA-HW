import json
from urllib.parse import urlparse

import pika
import requests
import validators
from bs4 import BeautifulSoup


class Searcher:
    def __init__(self, url):
        self.paths = {url: []}
        self.visited = {url}

    def step_and_check(self, check_link):
        new_paths = {}
        for cur, path in self.paths.items():
            try:
                new_links = get_links(cur)
                new_path = path + [cur]
                for link in new_links:
                    # print(f"*** Got link: {link}", flush=True)
                    if check_link == link:
                        return new_path + [check_link]
                    if link in self.visited:
                        continue
                    else:
                        self.visited.add(link)
                    new_paths[link] = new_path
            except RuntimeError as err:
                print(err)
        self.paths = new_paths


def get_base(url):
    parsed_uri = urlparse(url)
    return '{uri.scheme}://{uri.netloc}'.format(uri=parsed_uri)


def get_links(base_url):
    try:
        resp = requests.get(base_url)
    except Exception:
        raise RuntimeError(f"Bad url {base_url}")
    if resp.status_code / 100 != 2:
        raise RuntimeError(f"Not 2** status code for {base_url}")
    soup = BeautifulSoup(resp.text, 'html.parser')
    for link in soup.findAll('a'):
        url = link.get('href')
        if not url:
            continue
        if "wikipedia" in url and validators.url(url):
            yield url.strip()
        if url.startswith("/wiki"):
            url = get_base(base_url) + url
            if validators.url(url):
                yield url.strip()


def measure_path_length(url1, url2):
    s1 = Searcher(url1)
    for i in range(1000):
        res = s1.step_and_check(url2)
        if res:
            return res
        print(f"Iter {i} from 1000 done", flush=True)
    return []


def main():
    connection = pika.BlockingConnection(pika.ConnectionParameters('rabbitmq'))
    channel = connection.channel()
    channel.queue_declare(queue='requests')
    channel.queue_declare(queue='responses')

    def callback(ch, method, properties, body):
        req = json.loads(body)
        print(f"Received {req}", flush=True)
        req["Path"] = measure_path_length(
            req['URL1'],  # "https://en.wikipedia.org/wiki/Friends",
            req['URL2']  # "https://en.wikipedia.org/wiki/The_Pilot_(Friends)"
        )
        channel.basic_publish(exchange='',
                              routing_key='responses',
                              body=json.dumps(req).encode('UTF-8'))
    print("Ready", flush=True)
    channel.basic_consume(queue='requests',
                          auto_ack=True,
                          on_message_callback=callback)
    channel.start_consuming()


main()
