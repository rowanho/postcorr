import os
import sys
import re

from PIL import Image
import pytesseract

# Runs tesseract ocr 
def tesseract_plaintext(filename):
    text = pytesseract.image_to_string(Image.open(filename)) 
    return text

# Traverses the dataset and runs ocr
def perform_traversal(rootDir, outdir):
    os.mkdir(outdir)
    for dir_name, subdir_list, file_list in os.walk(rootDir):
        if len(subdir_list) > 0:
            for sub in subdir_list:
                perform_traversal_single(os.path.join(rootDir, sub), os.path.join(outdir, sub))
        else:
            perform_traversal_single(rootDir, outdir)            

# Traverses the dataset and runs ocr
def perform_traversal_single(rootDir, outdir):
    os.mkdir(outdir)
    for dir_name, subdir_list, file_list in os.walk(rootDir):
        for fname in file_list:
            image_text = tesseract_plaintext(os.path.join(dir_name, fname))
            froot = fname.split('.')[0]
            with open(os.path.join(outdir,  froot + '.txt'), 'a') as f:
                f.write(image_text)
def main():
    image_path = sys.argv[1]
    output_path = sys.argv[2]
    perform_traversal(image_path, output_path)
         
    
if __name__ == '__main__':
    main()