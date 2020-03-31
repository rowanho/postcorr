exts=( "1973-11-10.pdf" "1974-06-25.pdf" "1974-11-16.pdf" "1975-06-17.pdf" "1975-11-15.pdf" "1976-07-13.pdf
" "1976-11-12.pdf" "1977-07-06.pdf" "1978-06-29.pdf")

for ext in "${exts[@]}"
do
	python python-utils/scraping/scrape_wikisource.py "/wiki/Index:AASHTO_USRN_${ext}" "real_datasets/meetings/imgs/${ext}" "real_datasets/meetings/gt/${ext}"
done