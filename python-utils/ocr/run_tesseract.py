import os
import sys

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
        for fname in file_list:
            image_text = tesseract_plaintext(os.path.join(dir_name, fname))
            froot = fname.split('.')[0]
            with open(os.path.join(outdir, froot + '.txt'), 'a') as f:
                f.write(image_text)
            

# Traverses the dataset and runs ocr
def perform_traversal_single(rootDir, outfile):
    total_text = []
    for dir_name, subdir_list, file_list in os.walk(rootDir):
        for fname in file_list:
            image_text = tesseract_plaintext(os.path.join(dir_name, fname))
            total_text.append(image_text)    
    with open(outfile, 'w') as f:
        f.write('\n'.join(total_text))

def main():
    image_path = sys.argv[1]
    output_path = sys.argv[2]
    try:
        output_type = sys.argv[3]
    except:
        output_type = "multiple"
    if output_type == "multiple":
        perform_traversal(image_path, output_path)
    else:
        perform_traversal_single(image_path, output_path)
         
    
if __name__ == '__main__':
    main()