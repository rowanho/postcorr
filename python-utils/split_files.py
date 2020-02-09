import os
import sys


def main(in_fn, out_dir, split_length):
    
    os.mkdir(out_dir)
    input_file = open(in_fn, 'r')
    text = input_file.read()
    input_file.close()
    i = 0
    for start in range(0, len(text), split_length):
        i += 1
        with open(os.path.join(out_dir, str(i)), 'w') as outfile:
            outfile.write(text[start: start + split_length])
    
    
    
    

if __name__ == "__main__":
    in_file = sys.argv[1]
    output_dir = sys.argv[2]
    split_length = int(sys.argv[3])
    main(in_file, output_dir, split_length)
    