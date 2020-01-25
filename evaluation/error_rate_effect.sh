errorRates=(0.05 0.1 0.15 0.2)
numDuplicates=(3 5 7 9)

for dups in "${numDuplicates[@]}"
do
    for err in "${errorRates[@]}"
    do
        printf "%d\n" ${dups} 
        printf "%f\n" ${err}
        python python-utils/error_seeder.py datasets/test_data gs_data lossy_data 0.5 ${dups} 50 ${err}
        ./postCorr -input=lossy_data -groundtruth=gs_data -fp=modp -p=1 -jaccard=0.01 -shingleSize=5 -write=false
        rm -r lossy_data gs_data
    done
done
