import os

from PIL import Image
import pytesseract

PATH_TO_DATASET = 'image_dataset/'
PATH_TO_OUTPUT = 'corrected_dataset/'

# Runs tesseract ocr 
def tesseract_plaintext(filename):
    text = pytesseract.image_to_string(Image.open(filename)) 
    return text

# Traverses the dataset and runs ocr
def perform_traversal(rootDir, outdir):
    for dir_name, subdir_list, file_list in os.walk(rootDir):
        for fname in file_list:
            #image_text = tesseract_plaintext(os.path.join(dir_name, fname))
            image_text = 'j'
            directory = ''.join(dir_name.split('/')[1:]) 
            if not os.path.exists(outdir + '/' + directory):
                print(outdir + directory)
                os.makedirs(outdir + directory)
            with open(outdir  + directory + '/' + fname.split('.')[0] + '.txt', 'w') as f:
                f.write(image_text)
            
perform_traversal(PATH_TO_DATASET, PATH_TO_OUTPUT)