import os
from error_seeder import create_artifical_reuse, seed_errors, seed_errors_word
from essential_generators import DocumentGenerator
from random import randint
from distutils.dir_util import copy_tree

# Generates random words up to c characters
def gen_words_up_to_c(c):
    gen = DocumentGenerator()
    wc = 0
    char_count = 0
    sents = []
    while char_count < c:
        sent = gen.sentence()
        if sent[-1] != '.':
            sent = sent + '.'
        print(len(sent))
        sents.append(sent)
        char_count += len(sent)
    return ''.join(sents)
    
    
        
def gen_random_files(dir_name, num_files, c, sub_err_rate, deletion_rate, reuse_rate, num_reuse, reuse_sent_length):
    os.mkdir(dir_name)
    gt = os.path.join(dir_name,'gt')
    err = os.path.join(dir_name,'err')
    os.mkdir(gt)
    for f in range(num_files):
        text = gen_words_up_to_c(c)
        with open(os.path.join(gt, f'{f}.txt'), 'w') as file:
            file.write(text)
    
    create_artifical_reuse(gt, reuse_rate, num_reuse, reuse_sent_length)
    seed_errors(gt, err,  sub_err_rate, deletion_rate)
    
def gen_random_files_word(dir_base, exts, num_files, c, sub_err_rate, word_sub_rates, reuse_rate, num_reuse, reuse_sent_length):
    
    for i, ext in enumerate(exts):
        dir_name = f'{dir_base}/{ext}_word'
        os.mkdir(dir_name)
        gt = os.path.join(dir_name,'gt')
        err = os.path.join(dir_name,'err')
        if ext == 0:
            os.mkdir(gt)
            for f in range(num_files):
                text  = gen_words_up_to_c(c)
                with open(os.path.join(gt, f'{f}.txt'), 'w') as file:
                    file.write(text)
            create_artifical_reuse(gt, reuse_rate, num_reuse, reuse_sent_length)
        else:
            
            seed_errors_word(f'{dir_base}/{0}_word/gt', gt, sub_err_rate)
        seed_errors(gt, err, sub_err_rate, 0.0)

def gen_subs_delete():
    gen_random_files('synthetic_data/align_performance/2_reuse', 20, 2000, 0.03, 0.01, 1.0, 1, 2)    
    gen_random_files('synthetic_data/align_performance/4_reuse', 20, 2000, 0.03, 0.01, 1.0, 1, 4)
    gen_random_files('synthetic_data/align_performance/6_reuse', 20, 2000, 0.03, 0.01, 1.0, 1, 6)
    gen_random_files('synthetic_data/align_performance/8_reuse', 20, 2000, 0.03, 0.01, 1.0, 1, 8)
    gen_random_files('synthetic_data/align_performance/10_reuse', 20, 2000, 0.03, 0.01, 1.0, 1, 10)
    
    
def gen_subs_words():
    gen_random_files_word('synthetic_data/gap_performance', [0, 4, 8, 12, 16, 20], 20, 2000, 0.05, [0.0, 0.04, 0.08, 0.12, 0.16, 0.2], 1.0, 3, 5)
    
    
if __name__ == '__main__':
    gen_subs_words()
