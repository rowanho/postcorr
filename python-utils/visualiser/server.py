from flask import Flask, request, jsonify
import sys
import os
import json#
import html
from uuid import uuid4

original_dir = sys.argv[1]
corrected_dir = 'corrected'
app = Flask(__name__)
def to_html(text):
	return html.escape(text).replace("\n", "<br>")

def open_log_json(fn):
	with open(os.path.join('logs', fn)) as b:
		buf = b.read()
		data_json = json.loads(buf)
		return data_json

text_json = open_log_json('vote_graph.json')
old_se_json = open_log_json('vote_start_ends_old.json')
new_se_json = open_log_json('vote_start_ends_new.json')
sub_edit_json = open_log_json('edit_graph_sub.json')
del_edit_json = open_log_json('edit_graph_del.json')
ins_edit_json = open_log_json('edit_graph_ins.json')
vote_json = open_log_json('vote_details.json')

@app.route('/serve_reuse', methods=['post'])
def serve_reuse():
	filename = request.form['filename']
	with open(os.path.join(corrected_dir, filename)) as b:
		file_text = b.read()
	with open(os.path.join(original_dir, filename)) as b:
		original_file_text = b.read()
	old_start_ends_list = old_se_json[filename]
	new_start_ends_list = new_se_json[filename]
	text_segment_list = text_json[filename]
	sub_edits_list = sub_edit_json[filename]
	del_edits_list = del_edit_json[filename]
	ins_edits_list = ins_edit_json[filename]
	vote_details = vote_json[filename]
	reuse_map = {}
	alternating_segments = ""
	if len(text_segment_list) == 0:
		alternating_segments = f'<div>{to_html(file_text)}</div>'

	prev = 0
	z = list(zip(text_segment_list, old_start_ends_list, new_start_ends_list))
	z.sort(key=lambda x: x[2]['start'])
	uz = list(zip(*z))
	text_segment_list, old_start_ends_list, new_start_ends_list = tuple(uz)
	for i, text_segment in enumerate(text_segment_list):
		old_start = old_start_ends_list[i]['start']
		old_end = old_start_ends_list[i]['end'] + 1
		start = new_start_ends_list[i]['start']
		end = new_start_ends_list[i]['end'] + 1
		text_segment[filename + '_unedited'] = get_encoded(original_file_text[old_start:old_end], old_start, del_edits_list)
		alternating_segments += f'<span>{to_html(file_text[prev:start])}</span>'
		encoded = get_encoded_double(file_text[start:end], start, sub_edits_list, ins_edits_list, vote_details)
		text_segment[filename] = encoded
		uid = str(uuid4())
		alternating_segments += f'<span class="reused" uid="{uid}">{encoded}</span>'
		prev = end
		if i == len(text_segment_list) - 1:
			alternating_segments += f'<span>{to_html(file_text[end:])}</span>'
		escaped_text_segment = {}
		reuse_map[uid] = get_modal_html(text_segment, filename)
	return jsonify({'segments':alternating_segments, 'reuse_map': reuse_map})

def get_encoded(text_segment, start, edits):
	if len(edits) == 0:
		return text_segment
	text_segment_chars = list(text_segment)
	pos = 0
	for i in range(start, len(text_segment_chars) + start):
		if str(i) in edits:
			typ = edits[str(i)]
			text_segment_chars[i - start] = f'<mark class="{typ}">' + to_html(text_segment_chars[i - start]) + '</mark>'
		else:
			text_segment_chars[i - start] = to_html(text_segment_chars[i - start])
	return ''.join(text_segment_chars)


def get_encoded_double(text_segment, start, sub_edits, ins_edits, vote_details):
	if len(sub_edits) + len(ins_edits) == 0:
		return text_segment
	text_segment_chars = list(text_segment)
	pos = 0
	ins_streak = 0
	uid = ''
	for i in range(start, len(text_segment_chars) + start):
		if str(i) in sub_edits:
			ins_streak = 0
			typ = sub_edits[str(i)]
			hover = get_hover_text(vote_details[str(i)], True)
			text_segment_chars[i - start] = f'<mark class="{typ} vote">' + to_html(text_segment_chars[i - start]) + hover + '</mark>'
		elif str(i) in ins_edits:
			ins_streak += 1
			hover = get_hover_text(vote_details[str(i - ins_streak)], False)
			typ = ins_edits[str(i)]
			text_segment_chars[i - start] = f'<mark class="{typ} vote">' + to_html(text_segment_chars[i - start]) + hover + '</mark>'
		else:
			ins_streak = 0
			text_segment_chars[i - start] = to_html(text_segment_chars[i - start])
	return ''.join(text_segment_chars)

def get_modal_html(text_segment, leader_key):
	html_text = []
	html_text.append('<table style="table-layout:fixed; width:100%; word-break: break-all">')
	build_row(f'Edited text ({leader_key})', text_segment[leader_key], html_text, esc=False)
	build_row(f'Unedited text ({leader_key})', text_segment[leader_key + '_unedited'], html_text, esc=False)
	for key, val in text_segment.items():
		if key != leader_key and key != (leader_key + '_unedited'):
			build_row(f'Witness ({key})', val, html_text)
	html_text.append('</table>')
	return '\n'.join(html_text)

def get_hover_text(jsn, is_sub):
	if is_sub:
		dct = jsn['EditDict']
	else:
		dct = jsn['InsertDict']
	mx = 0
	mx_key = ''
	for key, val in dct.items():
		if val > mx:
			mx = val
			mx_key = key
	text = []
	text.append(f'{process(mx_key)}: {mx}')
	for key, val in dct.items():
		if key != mx_key:
				text.append(f'{process(key)}: {val}')
	t = '\n'.join(text)
	r =  f'<span class = "hovertext"> {t}</span>'
	return r

def process(txt):
	return to_html(txt).replace('<br>', ' ')

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
