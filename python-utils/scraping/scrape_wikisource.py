import collections
import os
import sys
from bs4 import BeautifulSoup
import requests

base = 'https://en.wikisource.org'    

def get_validated_links():
    ext = '/w/index.php?title=Category:Index_Validated'
    links_found = []
    has_next = True
    while has_next:
        res = requests.get(base + ext)
        parser = BeautifulSoup(res.text, 'html.parser')
        link_lists = parser.find_all('div', class_='mw-category-group')
        
        for lst in link_lists:
            links = lst.find_all('a')
            for l in links:
                links_found.append(l['href'])
                
        # Find the 'next page' link
        next = parser.find('a', title='Category:Index Validated', text='next page')
        if next == None:
            has_next = False
        else:
            ext = next['href']
            
    return links_found
    
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

def get_page_links(ext):
    links = []
    res = requests.get(base + ext)
    parser = BeautifulSoup(res.text, 'html.parser')
    potential_page_links = parser.find_all('a', class_='prp-pagequality-4 quality4')
    for p in potential_page_links:
        if p.getText().isnumeric():
            links.append(p['href'])
    return links

def get_high_scorers(links, n, k):
    link_scores = []
    for l in links:
        page_texts = []
        page_links = get_page_links(l)
        for pl in page_links:
            page_texts.append(get_page_text(pl))
        text = '\n'.join(page_texts)
        score = get_reuse_metric(k, text)
        print(len(text), score)
        link_scores.append((score, l))
        
    link_scores.sort(key=lambda x: x[0], reverse=True)
    return link_scores[:n]    
        
def get_page_text(ext):
    text = []
    res = requests.get(base + ext)
    parser = BeautifulSoup(res.text, 'html.parser')
    wrapper_div = parser.find('div', class_='mw-parser-output')
    paragraphs = wrapper_div.find_all('p')
    for p in paragraphs:
        text.append(p.getText())
    return '\n'.join(text)
    
    
def download_page_texts(dir, ext):
    n = 0
    page_links = get_page_links(ext)
    os.mkdir(dir)
    for l in page_links:
        text = get_page_text(l)    
        with open(os.path.join(dir, f'{n}.txt')) as file:
            file.write(text)
        n += 1

def download_page_images(dir, ext):
    n = 0
    page_links = get_page_links(ext)
    os.mkdir(dir)
    for l in page_links:
        text = get_page_text(l)    
        with open(os.path.join(dir, f'{n}.txt')) as file:
            file.write(text)
        n += 1
    
def hs():
    links = get_validated_links()
    hs = get_high_scorers(links, 10, 30)
    print(hs)

def main():
    ext = sys.argv[1]
    image_dir = sys.argv[2]
    text_dir = sys.argv[3]
    
    download_page_images(ext, image_dir)
    download_page_texts(ext, text_dir)
    
if __name__ == '__main__':
    main()