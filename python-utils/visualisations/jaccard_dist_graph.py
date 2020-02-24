import sys
import matplotlib.pyplot as plt

def jaccards_from_file(fn):
    jaccards = []
    with open(fn) as f:
        for j in f:
            jaccards.append(float(j))
    return jaccards
    
def main():
    root_dir = sys.argv[1]
    ks = [3,5,7]
    labels = [f'k={k}' for k in ks]
    hist_arrays = []
    for k in ks:
        fn = f'{root_dir}_jaccard_indexes{k}.txt'
        hist_arrays.append(jaccards_from_file(fn))
    plt.hist(hist_arrays, bins=100, label=labels)
    plt.xlabel('Jaccard Index Score')
    plt.ylabel('Frequency')
    plt.legend()
    plt.show()
    
if __name__ == "__main__":
    main()