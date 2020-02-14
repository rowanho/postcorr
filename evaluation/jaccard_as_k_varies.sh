shingleSizes=(2 3 4 5 6 7 8 9 10)

for shingleSize in "${shingleSizes[@]}"
do
    printf "%d\n" ${shingleSize}
    ./postCorr -input=$1 -fp=modp -p=10 -jaccard=0.00 -shingleSize=${shingleSize} -write=false -align=False -writeData=true
done
