import collections
import urllib.request
import os


def scrape_doc(url):
    ret = []
    try:
        data = urllib.request.urlopen(url)
        for line in data:
            ret.append(line.decode('utf-8'))
        return '\n'.join(ret)
    except urllib.error.HTTPError as e:
        return ""

def get_reuse_metric(k, text):
    freqs = collections.defaultdict(int)
    for i in range(0, len(text) + 1- k):
        freqs[text[i: i + k]] += 1
    s = 0
    i = 0
    for val in freqs.values():
        s += val
        i += 1
    if i == 0:
        return 0
    else:
        return s / i
        
def main():
    base = "http://www.gutenberg.org/files"
    metrics = []
    for no in range(1, 10):
        print(no)
        url = f"{base}/{no}/{no}-0.txt"
        alt_url = f"{base}/{no}/{no}.txt"
        text = scrape_doc(url)
        text2 = scrape_doc(alt_url)
        text = text2 if len(text2) > len(text) else text
        print(len(text))
        m = get_reuse_metric(5, text)
        metrics.append((m, no))
    
    metrics.sort(key=lambda x: x[0], reverse=True)
    print(metrics[:10])

if __name__ == "__main__":
    main()
    
    