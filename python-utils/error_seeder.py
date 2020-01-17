import os
import random
def seed_errors(dirname, error_prob):
    for filename in os.listdir(dirname):
        with open(os.path.join(dirname, filename), 'r') as file:
            charlist = list(file.read())
        for i, c in enumerate(charlist):
            if random.random() < error_prob:
                charlist[i] = str(chr(random.randint(35, 120)))
        content = ''.join(charlist)
        with open(os.path.join(dirname, filename), 'w') as file:
            file.write(content)
            
            
if __name__ == '__main__':
    seed_errors('lossy_test_data', 0.05)