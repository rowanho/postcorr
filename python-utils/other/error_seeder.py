import os
import random
import shutil
import sys
import argparse

import nltk


def get_errs(text, sub_err_rate, deletion_rate):
    char_list = list(text)
    l = len(char_list)
    while i < l:
        if random() < sub_error_rate:
            char_list[i] = str(chr(random.randint(35, 120)))
        if random() < deletion_rate:
            char_list.pop(i)
            l -= 1
        i += 1
    return ''.join(char_list)
    
tokenizer = nltk.data.load('tokenizers/punkt/english.pickle')

# Randomnly seeds documents with errors
def seed_errors(dirname, new_dir, sub_err_rate, deletion_rate):
    create_new_dir(dirname, new_dir)
    for filename in os.listdir(new_dir):
        with open(os.path.join(new_dir, filename), 'r') as file:
            content = get_errs(file.read(), sub_err_rate, deletion_rate)
        with open(os.path.join(new_dir, filename), 'w') as file:
            file.write(content)
            
# Artifically inserts local text reuse into our dataset
# Avg reuse length is the reuse in number of sentences
def create_artifical_reuse(dirname, reuse_prob, num_reuses_per_doc, avg_reuse_length):
    docs = {}
    docs_no_reuse = {}
    i = 0
    for filename in os.listdir(dirname):
        with open(os.path.join(dirname, filename), 'r') as file:
            buff = file.read().decode('utf-8')
            sents = tokenizer.tokenize(buff)
            for i,_ in enumerate(sents):
                sents[i] = sents[i] + '\n'
            docs[i] = sents
            docs_no_reuse[filename] = sents.copy()
        i += 1
    
    no_docs = i
    j = 0
    for filename in os.listdir(dirname):
        if random.random() < reuse_prob: # Decide whether this doc is the source of some text to be reused
            l = len(docs_no_reuse[filename])
            size = avg_reuse_length
            start = random.randint(0, l - size - 1 )
            # Reuse the text n times in other documents
            for n in range(num_reuses_per_doc - 1, num_reuses_per_doc + 1):
                p = random.randint(0, no_docs-1)
                while p == j:
                    p = random.randint(0, no_docs-1)    
                # Insert the reused text at some random position
                line_indexes = []
                reuse_point = docs[p][random.randint(0, len(docs[p]) - 1)]  + 1
                docs[p] = docs[p][:reuse_point] +  docs_no_reuse[filename][start:start + size]  + docs[p][reuse_point:]
        j += 1        
    
    i = 0
    for filename in os.listdir(dirname):
        with open(os.path.join(dirname, filename), 'w') as file:
            file.write(docs[i])
        i += 1    
    

def create_new_dir(src_dir, dst):
    shutil.copytree(src_dir, dst, symlinks=False, ignore=None)


def main():
    src_dir = sys.argv[1]
    gold_stand = sys.argv[2]
    dst =  sys.argv[3]  
    prob_reuse = float(sys.argv[4])
    num_reuses = int(sys.argv[5])
    reuse_size  = int(sys.argv[6])
    error_rate = float(sys.argv[7])
    create_new_dir(src_dir, gold_stand)
    create_artifical_reuse(gold_stand, prob_reuse, num_reuses, reuse_size)
    seed_errors(gold_stand, dst, error_rate, 0.01)
    
if __name__ == '__main__':
    main()