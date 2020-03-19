gaps=(0 4 8 12 16 20)
dir="synthetic_data/gap_performance"
for g in "${gaps[@]}"
do
    ./postCorr -input="${dir}/${g}_word/err" -groundtruth="${dir}/${g}_word/gt" -fp=modp -proportion=1.0 -write=False -fastAlign=false
    ./postCorr -input="${dir}/${g}_word/err" -groundtruth="${dir}/${g}_word/gt" -fp=modp -proportion=1.0 -write=False -fastAlign=false -affine=true
done
