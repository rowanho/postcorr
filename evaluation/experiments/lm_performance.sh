ths=(0.5 0.25 0.125 0.06125 0.0306125)
dir="real_datasets/denmark"
./postCorr -input="${dir}/den_noisy_pages" -groundtruth="${dir}/den_pages" -fp=modp -proportion=0.1 -write=False -fastAlign=false
for th in "${ths[@]}"
do
    ./postCorr -input="${dir}/den_noisy_pages" -groundtruth="${dir}/den_pages" -fp=modp -proportion=0.1 -write=False -fastAlign=false -useLM=true -lmThreshold=${th}
done
