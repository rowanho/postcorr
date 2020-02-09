import os
import random
import sys
import shutil

# Randomnly seeds documents with errors
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
    
def create_new_dir(src_dir, dst):
    shutil.copytree(src_dir, dst, symlinks=False, ignore=None)

def main():
    src_dir = sys.argv[1]
    dst = sys.argv[2]
    error_rate = float(sys.argv[3])
    seed_errors(src_dir, dst, error_rate)
if __name__ == "__main__":
    main()
        