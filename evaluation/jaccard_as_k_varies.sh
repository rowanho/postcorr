shingleSizes=(3 4 5 6 7 8)

for shingleSize in "${shingleSizes[@]}"
do
    printf "%d\n" ${shingleSize}
    ./postCorr -input=$1 -fp=modp -p=1 -jaccard=0.1 -shingleSize=${shingleSize} -write=False -align=False -writeData=True
done
