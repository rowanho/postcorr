proportion=(0.1 0.2 0.3 0.4 0.5 0.6 0.7 0.8 0.9 1.0)

for p in ${proportion[@]}
do
  ./postCorr -input=real_datasets/meetings/11 -groundtruth=real_datasets/meetings/ocr -candidate_proportion=${p} -fp=minhash -use_lm=true -lm_threshold=0.2 -affine=true -fast_align=true -band_width=100
done
