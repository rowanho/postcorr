from bs4 import BeautifulSoup
import requests

def get_validated_links():
    base = 'https://en.wikisource.org'    
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
        print(next)
        if next == None:
            has_next = False
        else:
            ext = next['href']
            
    return links_found
    
    
if __name__ == '__main__':
    links = get_validated_links()
    print(links)
