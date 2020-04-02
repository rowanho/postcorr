errorRates=(0.05 0.1 0.15 0.20 0.25 0.3 0.35 0.4)

for err in "${errorRates[@]}"
do
    printf "%f\n" ${err}
    for i in {1..4}
    do
        python python-utils/simple_error_seeder.py datasets/100_data lossy_data  ${err}
        ./postCorr -input=lossy_data -groundtruth=datasets/100_data -fp=modp -p=1 -jaccard=0.00 -shingleSize=5 -write=false -fastAlign=true
        rm -r lossy_data 
    done
done
