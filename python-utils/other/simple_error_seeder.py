import os
import random
import sys
import shutil

common_mistakes = {'m':['n', 'w'], 'n':['m', 'o'], 'o':['a', 'u'], 'a':['o', 'u'],
'o':['a', 'u'], 'e':['a'], 'l':['I'], 'I':['l'], 'F':['E'], 'E':['F'], 'j':['i'], 'i':['j', 'l'], 'd':['b'], 'b':['d']}
# Randomnly seeds documents with errors
def seed_errors(dirname, new_dir, error_prob):
    print(dirname)
    for filename in os.listdir(new_dir):
        with open(os.path.join(new_dir, filename), 'r') as file:
            charlist = list(file.read())
        for i, c in enumerate(charlist):
            if random.random() < error_prob:
                if charlist[i] in common_mistakes:
                    charlist[i] = random.choice(common_mistakes[charlist[i]])
        content = ''.join(charlist)
        with open(os.path.join(new_dir, filename), 'w') as file:
            file.write(content)

def seed_errors_rec(dirname, new_dir, error_prob):
    for dir_name, subdir_list, file_list in os.walk(dirname):
        if len(subdir_list) > 0:
            for sub in subdir_list:
                print('for')
                seed_errors_rec(os.path.join(dirname, sub), os.path.join(new_dir, sub), error_prob)
        else:
            seed_errors(dirname, new_dir, error_prob)
        break
def create_new_dir(src_dir, dst):
    shutil.copytree(src_dir, dst, symlinks=False, ignore=None)

def main():
    src_dir = sys.argv[1]
    dst = sys.argv[2]
    error_rate = float(sys.argv[3])
    create_new_dir(src_dir, dst)
    seed_errors_rec(src_dir, dst, error_rate)
if __name__ == "__main__":
    main()
