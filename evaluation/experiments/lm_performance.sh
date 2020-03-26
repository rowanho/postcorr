ths=(0.5 0.25 0.125 0.06125 0.0306125)
dir="real_datasets/copyright"
./postCorr -input="${dir}/ocr" -groundtruth="${dir}/gt" -fp=modp -proportion=0.01 -write=False -fastAlign=false
for th in "${ths[@]}"
do
    ./postCorr -input="${dir}/ocr" -groundtruth="${dir}/gt" -fp=modp -proportion=0.01 -write=False -fastAlign=false -useLM=true -lmThreshold=${th}
done
