import requests
import os
import sys

def write_image(url, out_dir):
    image_data = requests.get(url).content
    with open(os.path.join(out_dir, url.split('/')[-1]), 'wb') as handler:
        handler.write(image_data)
    
def main(base_path, out_dir):
    ext = "png"
    n = 1
    os.mkdir(out_dir)
    while True:
        try:
             path = "{}{:03d}.png".format(base_path, n)
             write_image(path, out_dir)
        except:
            break
        n += 1
    

if __name__ == "__main__":
    main(sys.argv[1], sys.argv[2])
    
    