import editdistance
import os

PATH_TO_UNCORRECTED = 'eval/ocr_dataset/'
PATH_TO_CORRECTED = 'eval/corrected_dataset/'
PATH_TO_GROUNDTRUTH = 'eval/groundtruth_dataset/'

def traverse_and_compare(uncorrected_path, corrected_path, groundtruth_path):
    for dir_name, subdir_list, file_list in os.walk(groundtruth_path):
        corrected_dir = corrected_path + ''.join(dir_name.split('/')[2:]) 
        uncorrected_dir = uncorrected_path + ''.join(dir_name.split('/')[2:])
        for fn in file_list:
            with open(dir_name + '/' + fn, 'r') as groundtruth_file:
                groundtruth_text = groundtruth_file.read()
                print(dir_name + '/' + fn)
                with open(uncorrected_dir + '/' + fn, 'r') as uncorrected_file:
                    uncorrected_text = uncorrected_file.read()
                    e = editdistance.eval(groundtruth_text, uncorrected_text)
                    print(e)
                with open(corrected_dir + '/' + fn, 'r') as corrected_file:
                    corrected_text = corrected_file.read()
                    e = editdistance.eval(groundtruth_text, corrected_text)
                    print(e)
                    
                    
traverse_and_compare(PATH_TO_UNCORRECTED, PATH_TO_CORRECTED, PATH_TO_GROUNDTRUTH)