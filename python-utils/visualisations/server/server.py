from flask import Flask, request, jsonify
import sys
import os
import json
path = sys.argv[1]

app = Flask(__name__)

@app.route('/serve_file', methods=['post'])
def serve_file():
	filename = request.form['filename']
	with open(os.path.join(path,filename)) as b:
		buf = b.read()
	return buf

@app.route('/serve_reuse', methods=['post'])
def serve_reuse():
	filename = request.form['filename']
	with open(os.path.join('logs', 'logs_reuse_graph.json')) as b:
		buf = b.read()
		j = json.loads()
	return jsonify(j[filename])

@app.route('/serve_improvements', methods=['post'])
def serve_improvements():
	filename = request.form['filename']
	with open(os.path.join('logs', 'logs_edit_graph.json')) as b:
		buf = b.read()
		j = json.loads()
	return jsonify(j[filename])


@app.route('/', methods=['get'])
def serve():
	return app.send_static_file('index.html')

if __name__ == "__main__":
    app.run(port=3000)
