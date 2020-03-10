import math
import nltk
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
    word = request.form['word']
    prev_word = request.form['sentence']
    logscore = lang_model.unmasked_score(word, sentence.split(' '))
    return str(-score)

if __name__ == "__main__":
    app.run()
