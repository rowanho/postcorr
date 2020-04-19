from flask import Flask, request, jsonify
import sys
import os
import json#
import html
from uuid import uuid4

original_dir = sys.argv[1]
corrected_dir = sys.argv[2]
app = Flask(__name__)

def to_html(text):
	return html.escape(text).replace("\n", "<br>")

@app.route('/serve_reuse', methods=['post'])
def serve_reuse():
	filename = request.form['filename']
	with open(os.path.join(corrected_dir, filename)) as b:
		file_text = b.read()
	with open(os.path.join(original_dir, filename)) as b:
		original_file_text = b.read()
	with open(os.path.join('logs', 'vote_graph.json')) as b:
		buf = b.read()
		text_json = json.loads(buf)
	with open(os.path.join('logs', 'vote_start_ends.json')) as b:
		buf = b.read()
		se_json = json.loads(buf)

	start_ends_list = se_json[filename]
	text_segment_list = text_json[filename]

	with open(os.path.join('logs', 'edit_graph.json')) as b:
		buf = b.read()
		edit_json = json.loads(buf)
	edits_list = edit_json[filename]

	reuse_map = {}
	alternating_segments = ""
	if len(text_segment_list) == 0:
		alternating_segments = f'<div>{to_html(file_text)}</div>'

	prev = 0

	for i, text_segment in enumerate(text_segment_list):
		start = start_ends_list[i]['start']
		end = start_ends_list[i]['end'] + 1
		text_segment[filename + '_unedited'] = original_file_text[start:end]
		alternating_segments += f'<span>{to_html(file_text[prev:start])}</span>'
		encoded = get_encoded(file_text[start:end], edits_list[i])
		uid = str(uuid4())
		alternating_segments += f'<span class="reused" uid="{uid}">{encoded}</span>'
		prev = end
		if i == len(text_segment_list) - 1:
			alternating_segments += f'<span>{to_html(file_text[end:])}</span>'
		escaped_text_segment = {}
		reuse_map[uid] = get_modal_html(text_segment, filename, edits_list[i])
	return jsonify({'segments':alternating_segments, 'reuse_map': reuse_map})


def get_encoded(text_segment, edits):
	if len(edits) == 0:
		return text_segment
	text_segment_chars = list(text_segment)
	pos = 0
	for i in range(len(text_segment_chars)):
		if str(i) in edits:
			typ = edits[str(i)]
			text_segment_chars[i] = f'<mark class="{typ}">' + to_html(text_segment_chars[i]) + '</mark>'
			print(text_segment_chars[i])
		else:
			text_segment_chars[i] = to_html(text_segment_chars[i])
	return ''.join(text_segment_chars)

def get_modal_html(text_segment, leader_key, edits):
	html_text = []
	html_text.append('<table style="table-layout:fixed; width:100%; word-break: break-all">')
	text_segment[leader_key] = get_encoded(text_segment[leader_key], edits)
	build_row(f'Edited text ({leader_key})', text_segment[leader_key], html_text, esc=False)
	build_row(f'Unedited text ({leader_key})', text_segment[leader_key + '_unedited'], html_text, esc=False)
	for key, val in text_segment.items():
		if key != leader_key and key != (leader_key + '_unedited'):
			build_row(f'Witness ({key})', val, html_text)
	html_text.append('</table>')
	return '\n'.join(html_text)

def build_row(key, val, html_text, esc=True):
	html_text.append('<tr>')
	html_text.append(f'<td>{to_html(key)}</td>')
	if esc:
		escaped = to_html(val).replace('<br>', ' ')
	else:
		escaped = val.replace('<br>', ' ')
	html_text.append(f'<td>{escaped}</td>')
	html_text.append('</tr>')



@app.route('/', methods=['get'])
def serve():
	return app.send_static_file('index.html')

if __name__ == "__main__":
    app.run(port=3000)
