import requests
from html.parser import HTMLParser

class Data:
    def __init__(self, x, y, character):
        self.x = x
        self.y = y
        self.character = character

class ParseHTML(HTMLParser):
    def __init__(self):
        super().__init__()
        self.in_table = False
        self.data = []

    def handle_starttag(self, tag, attrs):
        if tag == "table":
            self.in_table = True

    def handle_endtag(self, tag):
        if tag == "table":
            self.in_table = False

    def handle_data(self, data):
        if self.in_table:
            cleaned = data.strip()
            if cleaned:
                self.data.append(cleaned)

def decode_from_url(url):
    response = requests.get(url)
    parser = ParseHTML()
    parser.feed(response.text)
    return parser.data

def insert_data(raw):
    dt = []
    for i in range(0, len(raw), 3):
        x = int(raw[i])
        char = raw[i + 1]
        y = int(raw[i + 2])
        dt.append(Data(x, y, char))
    return dt

def get_array_size(dt):
    max_x = max(data.x for data in dt)
    max_y = max(data.y for data in dt)
    return max_x + 1, max_y + 1

def map_data(dt):
    x, y = get_array_size(dt)
    arr = [[" " for _ in range(x)] for _ in range(y)]

    for data in dt:
        arr[data.y][data.x] = data.character

    for row in arr:
        print("".join(row))

def main():
    url = "https://docs.google.com/document/d/e/2PACX-1vSZ1vDD85PCR1d5QC2XwbXClC1Kuh3a4u0y3VbTvTFQI53erafhUkGot24ulET8ZRqFSzYoi3pLTGwM/pub"
    raw_data = decode_from_url(url)
    dt = insert_data(raw_data[3:])
    map_data(dt)

if __name__ == "__main__":
    main()