import math
import nltk
import re
from nltk.lm import MLE, Vocabulary
from nltk.lm.preprocessing import padded_everygram_pipeline
from nltk.corpus import brown
from flask import Flask, request
app = Flask(__name__)

n = 2
print(brown.sents())
def lower_list(l):
    keep = []
    for i in range(len(l)):
        if re.fullmatch('[a-zA-Z]+', l[i]):
            keep.append(l[i].lower())
    return keep
        
train, vocab = padded_everygram_pipeline(n, [lower_list(sent) for sent in brown.sents()])

#train, vocab = padded_everygram_pipeline(n, [['then', 'she', 'said'], ['then', 'he', 'said']])
lang_model = MLE(n)
lang_model.fit(train, vocab)
print(lang_model.logscore('atlanta', ['city', 'of']))
print(lang_model.logscore('the', ['he', 'saw']))

@app.route('/', methods=['post'])
def get_lang_probability():
    j = request.get_json()
    word = j['word']
    sentence = j['sentence']
    word = re.sub('[^a-zA-Z]+', '', word)
    sen_words = sentence.split(' ')
    for i, _ in enumerate(sen_words):
        sen_words[i] = re.sub('[^a-zA-Z]+', '', sen_words[i])
    valid_words = []
    for s in sen_words:
        if len(s) > 0:
            valid_words.append(s.lower())
    print(word)
    print(valid_words)
    logscore = lang_model.logscore(word.lower(), valid_words)
    print(logscore)
    if logscore == 0.0:
        return "0.0"
    else:
        return str(-logscore)

if __name__ == "__main__":
    app.run()
