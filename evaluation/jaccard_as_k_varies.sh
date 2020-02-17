shingleSizes=(5 6 7 8 9 10)

for shingleSize in "${shingleSizes[@]}"
do
    printf "%d\n" ${shingleSize}
    ./postCorr -input=$1 -groundtruth=$2 -fp=modp -p=10 -jaccard=0.1 -shingleSize=${shingleSize} -write=False -align=True -writeData=True -fastAlign=True -parallel=True
done
