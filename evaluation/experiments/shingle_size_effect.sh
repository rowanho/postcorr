shingleSizes=(3 5 7 9 11)
errorRates=(0.00 0.05 0.1 0.15)

for err in "${errorRates[@]}"
    do
        python python-utils/simple_error_seeder.py datasets/1000_different_data lossy_data ${err} 
        python python-utils/simple_error_seeder.py datasets/1000_same_data1 lossy_data1 ${err} 
        python python-utils/simple_error_seeder.py datasets/1000_same_data2 lossy_data2 ${err} 
        python python-utils/simple_error_seeder.py datasets/1000_same_data3 lossy_data3 ${err} 
        printf "%f\n" ${err}
        for shingleSize in "${shingleSizes[@]}"
        do
            printf "%f\n" ${shingleSize}
            ./postCorr -input=lossy_data1  -fp=modp -p=1 -jaccard=0.00 -shingleSize=${shingleSize} -write=false -align=False
            ./postCorr -input=lossy_data2 -fp=modp -p=1 -jaccard=0.00 -shingleSize=${shingleSize} -write=false -align=False
            ./postCorr -input=lossy_data3 -fp=modp -p=1 -jaccard=0.00 -shingleSize=${shingleSize} -write=false -align=False
            ./postCorr -input=lossy_data -fp=modp -p=1 -jaccard=0.00 -shingleSize=${shingleSize} -write=false -align=False
        done
        rm -r lossy_data1 lossy_data2 lossy_data3 
        rm -r lossy_data
    done
