import os
from error_seeder import create_artifical_reuse, seed_errors
from essential_generators import DocumentGenerator
from random import random, randint
from distutils.dir_util import copy_tree

# Generates random words up to c characters
def gen_words_up_to_c(c):
    gen = DocumentGenerator()
    wc = 0
    char_count = 0
    sents = []
    while char_count < c:
        sent = gen.sentence()
        print(sent)
        sents.append(sent)
        char_count += len(sent)
    return ''.join(sents)
    
    
        
def gen_random_files(dir_name, num_files, c, sub_err_rate, deletion_rate, reuse_rate, num_reuse, reuse_sent_length):
    os.mkdir(dir_name)
    gt = os.path.join(dir_name,'gt')
    err = os.path.join(dir_name,'err')
    os.mkdir(gt)
    os.mkdir(err)
    for f in range(num_files):
        text = gen_words_up_to_c(c)
        with open(os.path.join(gt, f'{f}.txt'), 'w') as file:
            file.write(text)
    
    create_artifical_reuse(gt, reuse_rate, num_reuse, reuse_sent_length)
    copy_tree(gt, err)
    seed_errors(err,  sub_err_rate, deletion_rate)
    
if __name__ == '__main__':
    gen_random_files('synthetic_data/benchmark_align/300_chars', 100, 300, 0.03, 0.01, 0.3, 2, 5)