import os
import random
import shutil
import sys

def seed_errors(dirname, new_dir, error_prob):
    create_new_dir(dirname, new_dir)
    for filename in os.listdir(new_dir):
        with open(os.path.join(new_dir, filename), 'r') as file:
            charlist = list(file.read())
        for i, c in enumerate(charlist):
            if random.random() < error_prob:
                charlist[i] = str(chr(random.randint(35, 120)))
        content = ''.join(charlist)
        with open(os.path.join(new_dir, filename), 'w') as file:
            file.write(content)
            

def create_artifical_reuse(dirname, reuse_prob, num_reuses_per_doc, avg_reuse_length):
    docs = {}
    docs_no_reuse = {}
    i = 0
    for filename in os.listdir(dirname):
        with open(os.path.join(dirname, filename), 'r') as file:
            buff = file.read()
            docs[filename] = buff
            docs_no_reuse[i] = buff
        i += 1
    for filename in os.listdir(dirname):
        if random.random() < reuse_prob:
            lines = docs[filename].split('\n')
            for n in range(num_reuses_per_doc - 1, num_reuses_per_doc + 1):
                p = random.randint(0, i-1)
                l = len(docs_no_reuse[p])
                size = random.randint(avg_reuse_length // 2, (avg_reuse_length * 3) // 2)
                start = random.randint(0, l - size - 1 )
                reuse_point = random.randint(0, len(lines) - 1)
                lines.insert(reuse_point, docs_no_reuse[p][start:start+size])
            with open(os.path.join(dirname, filename), 'w') as file:
                file.write('\n'.join(lines))
            
    

def create_new_dir(src_dir, dst):
    shutil.copytree(src_dir, dst, symlinks=False, ignore=None)


def main():
    src_dir = sys.argv[1]
    gold_stand = sys.argv[2]
    dst =  sys.argv[3]   
    
    create_new_dir(src_dir, gold_stand)
    create_artifical_reuse(gold_stand, 0.5, 3, 60)
    seed_errors(gold_stand, dst, 0.05)
    
if __name__ == '__main__':
    main()