# Clears out the test indexes 
# For dev purposes 

curl -X DELETE "localhost:9200/test_documents"
curl -X DELETE "localhost:9200/test_fingerprints"
curl -X DELETE "localhost:9200/test_minhash_fingerprints"
curl -X DELETE "localhost:9200/test_alignments"