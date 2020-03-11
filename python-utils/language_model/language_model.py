import math
import nltk
import re
from nltk.lm import MLE, Vocabulary
from nltk.lm.preprocessing import padded_everygram_pipeline
from nltk.corpus import brown
from flask import Flask, request
app = Flask(__name__)

n = 4
train, vocab = padded_everygram_pipeline(n, brown.sents())
lang_model = MLE(n)
lang_model.fit(train, vocab)
print(lang_model.score('bob', ['is']))

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
            valid_words.append(s)
    print(word)
    print(valid_words)
    logscore = lang_model.logscore(word, valid_words)
    print(str(-logscore))
    return str(-logscore)

if __name__ == "__main__":
    app.run()
