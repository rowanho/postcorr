import os
import xml.etree.ElementTree as ET
def parse_body_to_text(dirname, outdirname):
    for filename in os.listdir(dirname):
        if filename.endswith(".xml"): 
            doc = ET.parse(os.path.join(dirname,filename))
            os.makedirs(outdirname)
            with open(os.path.join(outdirname, filename),'w') as outfile:
                outfile.write(doc.find('body').text.encode('utf-8'))
            
parse_body_to_text('COUNTER', 'counter_plain')