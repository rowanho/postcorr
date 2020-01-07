import urllib.request
import os
from datetime import date, timedelta

# The following functions download 'chronicling america archived newspapers'
class Chronicling_america():
    
    def __init__(self):
        self.base  = "https://chroniclingamerica.loc.gov/"
        self.ocr_end = "/ed-1/seq-1/ocr.txt"
    
    # Evening herald 1914 - 22
    def evening_herald(self):
        start_date = date(1914, 1, 1)
        end_date = date(1922, 12, 31)
        delta = timedelta(days=1)
        folder = 'evening_herald_1914-1922' 
        os.mkdir(folder)
        while start_date <= end_date:
            date_string = start_date.strftime("%Y-%m-%d")
            url = self.base +"lccn/sn92070582/" + date_string + self.ocr_end
            try:
                data = urllib.request.urlopen(url)
                with open(os.path.join(folder, date_string + '.txt'), 'w') as file:
                    for line in data:
                        file.write(line.decode("utf-8"))
            except urllib.error.HTTPError as e:
                pass
            print(url)
            start_date += delta
    
    # Evening star 1854-1972
    def evening_star(self):
        start_date = date(1856, 1, 17)
        end_date = date(1972, 12, 31)
        delta = timedelta(days=1)
        folder = 'evening_star_1854-1972' 
        #os.mkdir(folder)
        while start_date <= end_date:
            date_string = start_date.strftime("%Y-%m-%d")
            url = self.base +"lccn/sn83045462/" + date_string + self.ocr_end
            try:
                data = urllib.request.urlopen(url)
                with open(os.path.join(folder, date_string + '.txt'), 'w') as file:
                    for line in data:
                        file.write(line.decode("utf-8"))
            except urllib.error.HTTPError as e:
                pass
            print(url)
            start_date += delta



def main():
    c = Chronicling_america()
    c.evening_star()
    
    
if __name__ == "__main__":
    main()
